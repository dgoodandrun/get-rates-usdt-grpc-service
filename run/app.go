package run

import (
	"context"
	"get-rates-usdt-grpc-service/config"
	"get-rates-usdt-grpc-service/internal/infrastracture/errors"
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


	a.grpcSrv = grpc.NewServer()
	pb.RegisterGetRatesServer(a.grpcSrv, ....)

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
