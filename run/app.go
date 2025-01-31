package run

import (
	"GeoService/geo-service/geo/config"
	"GeoService/geo-service/geo/internal/infrastructure/cache"
	"GeoService/geo-service/geo/internal/infrastructure/errors"
	"GeoService/geo-service/geo/internal/infrastructure/limiter"
	"GeoService/geo-service/geo/internal/modules/controller"
	"GeoService/geo-service/geo/internal/modules/service"
	"context"
	pb "github.com/dgoodandrun/GeoService-protos/geoPb"
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"gitlab.com/ptflp/gopubsub/kafkamq"
	"gitlab.com/ptflp/gopubsub/queue"
	"gitlab.com/ptflp/gopubsub/rabbitmq"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// Application - интерфейс приложения
type Application interface {
	Runner
	Bootstraper
}

// Runner - интерфейс запуска приложения
type Runner interface {
	Run() int
}

// Bootstraper - интерфейс инициализации приложения
type Bootstraper interface {
	Bootstrap(options ...interface{}) Runner
}

// App - структура приложения
type App struct {
	conf    config.AppConf
	logger  *zap.SugaredLogger
	grpcSrv *grpc.Server
	Sig     chan os.Signal
	lis     net.Listener
	conn    *amqp.Connection
	mq      queue.MessageQueuer
}

// NewApp - конструктор приложения
func NewApp(conf config.AppConf, logger *zap.SugaredLogger) *App {
	return &App{conf: conf, logger: logger, Sig: make(chan os.Signal, 1)}
}

// Bootstrap - инициализация приложения
func (a *App) Bootstrap(options ...interface{}) Runner {
	lis, err := net.Listen("tcp", ":"+a.conf.Port)
	if err != nil {
		a.logger.Fatal("failed to listen tcp Geo service: ", err)
	}
	a.lis = lis

	if a.conf.APIKey == "" || a.conf.SecretKey == "" {
		a.logger.Fatal("API Key and Secret Key must be set")
	}

	cacheClient := redis.NewClient(&redis.Options{
		Addr:     a.conf.RedisAddr,
		Password: "",
		DB:       0,
	})

	switch a.conf.Queue {
	case "rabbitmq":
		a.conn, err = amqp.Dial(a.conf.RabbitMQAddr)
		if err != nil {
			a.logger.Fatal(err)
		}
		a.mq, err = rabbitmq.NewRabbitMQ(a.conn)
		if err != nil {
			a.logger.Fatal(err)
		}
		if err = rabbitmq.CreateExchange(a.conn, a.conf.Topic, "direct"); err != nil {
			a.logger.Fatal(err)
		}
	case "kafka":
		a.mq, err = kafkamq.NewKafkaMQ(a.conf.KafkaAddr, "myGroup")
		if err != nil {
			a.logger.Fatal(err)
		}
	default:
		a.logger.Fatal("Queue type error: ", a.conf.Queue)
	}

	rateLimiter := limiter.NewRateLimiter(a.conf.Rate)
	geoService := service.NewGeoService(a.conf.APIKey, a.conf.SecretKey)
	geoServiceProxy := cache.NewGeoServiceProxy(geoService, cacheClient)
	geoController := controller.NewGeoProvider(geoServiceProxy, rateLimiter, a.mq)

	a.grpcSrv = grpc.NewServer()
	pb.RegisterGeoProviderServer(a.grpcSrv, geoController)

	// возвращаем приложение
	return a
}

// Run - запуск приложения
func (a *App) Run() int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errGroup, ctx := errgroup.WithContext(ctx)

	// Устанавливаем обработчик сигналов
	signal.Notify(a.Sig, syscall.SIGINT, syscall.SIGTERM)

	// Горутина для обработки сигналов завершения работы
	errGroup.Go(func() error {
		defer func() {
			if a.conn != nil {
				a.conn.Close()
			}
			if a.mq != nil {
				a.mq.Close()
			}
		}()

		select {
		case sig := <-a.Sig:
			a.logger.Infof("Received signal: %v, shutting down...", sig)
			cancel()
			a.grpcSrv.GracefulStop()
			a.logger.Info("Geo service stopped gracefully")
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	// Горутина для запуска gRPC сервера
	errGroup.Go(func() error {
		a.logger.Infof("Geo service started on port: %s", a.conf.Port)
		if err := a.grpcSrv.Serve(a.lis); err != nil && err != grpc.ErrServerStopped {
			a.logger.Errorf("Failed to serve Geo service: %v", err)
			return err
		}
		return nil
	})

	// Ожидаем завершения всех горутин
	if err := errGroup.Wait(); err != nil {
		a.logger.Errorf("Geo service shutdown with error: %v", err)
		return errors.GeneralError
	}

	return errors.NoError
}
