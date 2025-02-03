package run

import (
	"context"
	"get-rates-usdt-grpc-service/config"
	"get-rates-usdt-grpc-service/internal/infrastracture/errors"
	"get-rates-usdt-grpc-service/internal/infrastracture/metrics"
	"get-rates-usdt-grpc-service/internal/infrastracture/trace"
	"get-rates-usdt-grpc-service/internal/modules/controller"
	"get-rates-usdt-grpc-service/internal/modules/service"
	"get-rates-usdt-grpc-service/internal/modules/storage"
	pb "get-rates-usdt-grpc-service/protogen/golang/get-rates"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	conf         config.AppConf
	logger       *zap.SugaredLogger
	grpcSrv      *grpc.Server
	Sig          chan os.Signal
	lis          net.Listener
	tracerCloser func()
}

// NewApp - конструктор приложения
func NewApp(conf config.AppConf, logger *zap.SugaredLogger) *App {
	return &App{conf: conf, logger: logger, Sig: make(chan os.Signal, 1)}
}

// Bootstrap - инициализация приложения
func (a *App) Bootstrap(options ...interface{}) Runner {
	metrics.InitMetrics(a.conf.MetricsPort)

	closer, err := trace.InitTracer(a.conf.AppName, a.logger)
	if err != nil {
		a.logger.Fatalf("failed to init tracer: %v", err)
	}
	a.tracerCloser = closer

	lis, err := net.Listen("tcp", ":"+a.conf.Port)
	if err != nil {
		a.logger.Fatalf("failed to listen tcp %s: %v", a.conf.AppName, err)
	}
	a.lis = lis

	pgStorage, err := storage.NewPostgresStorage(a.conf.Postgres)
	//chStorage, err := storage.NewClickHouseStorage(a.conf.ClickHouse)
	if err != nil {
		a.logger.Fatalf("Failed to init database: %v", err)
	}

	ratesService := service.NewRatesService(pgStorage, a.conf.GarantexURL, a.conf.Market)
	ratesController := controller.NewRatesController(ratesService)

	a.grpcSrv = grpc.NewServer()
	pb.RegisterRatesServiceServer(a.grpcSrv, ratesController)
	reflection.Register(a.grpcSrv)

	// возвращаем приложение
	return a
}

// Run - запуск приложения
func (a *App) Run() int {
	defer func() {
		if a.tracerCloser != nil {
			a.tracerCloser()
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Устанавливаем обработчик сигналов
	signal.Notify(a.Sig, syscall.SIGINT, syscall.SIGTERM)

	eg, ctx := errgroup.WithContext(ctx)

	// Горутина для обработки сигналов завершения работы
	eg.Go(func() error {
		select {
		case sig := <-a.Sig:
			a.logger.Infof("Received signal: %v, shutting down...", sig)
			cancel()
			a.grpcSrv.GracefulStop()
			a.logger.Infof("%s stopped gracefully", a.conf.AppName)
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	// Горутина для запуска gRPC сервера
	eg.Go(func() error {
		a.logger.Infof("%s started on port: %s", a.conf.AppName, a.conf.Port)
		if err := a.grpcSrv.Serve(a.lis); err != nil && err != grpc.ErrServerStopped {
			a.logger.Errorf("Failed to serve %s: %v", a.conf.AppName, err)
			return err
		}
		return nil
	})

	// Ожидаем завершения всех горутин
	if err := eg.Wait(); err != nil {
		a.logger.Errorf("%s shutdown with error: %v", a.conf.AppName, err)
		return errors.GeneralError
	}

	return errors.NoError
}
