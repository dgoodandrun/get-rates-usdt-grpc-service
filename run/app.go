package run

import (
	"context"
	"get-rates-usdt-grpc-service/config"
	"get-rates-usdt-grpc-service/internal/infrastracture/errors"
	"get-rates-usdt-grpc-service/internal/modules/controller"
	"get-rates-usdt-grpc-service/internal/modules/service"
	"get-rates-usdt-grpc-service/internal/modules/storage"
	pb "get-rates-usdt-grpc-service/protogen/golang/get-rates"
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

	chStorage, err := storage.NewClickHouseStorage(a.conf.ClickHouse)
	if err != nil {
		a.logger.Fatal("Failed to init ClickHouse: ", err)
	}

	ratesService := service.NewRatesService(chStorage, a.conf.GarantexURL)
	ratesController := controller.NewRatesController(ratesService)

	a.grpcSrv = grpc.NewServer()
	pb.RegisterRatesServiceServer(a.grpcSrv, ratesController)

	// возвращаем приложение
	return a
}

// Run - запуск приложения
func (a *App) Run() int {
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
			a.logger.Info("Geo service stopped gracefully")
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	// Горутина для запуска gRPC сервера
	eg.Go(func() error {
		a.logger.Infof("Geo service started on port: %s", a.conf.Port)
		if err := a.grpcSrv.Serve(a.lis); err != nil && err != grpc.ErrServerStopped {
			a.logger.Errorf("Failed to serve Geo service: %v", err)
			return err
		}
		return nil
	})

	// Ожидаем завершения всех горутин
	if err := eg.Wait(); err != nil {
		a.logger.Errorf("Geo service shutdown with error: %v", err)
		return errors.GeneralError
	}

	return errors.NoError
}
