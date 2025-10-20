package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/carloscacb333/go-hexagonal/app/bootstrap"
	"go.uber.org/zap"
)

func main() {
	// Inicializar aplicación
	app, err := bootstrap.NewApplication()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	logger := app.GetLogger()
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("Failed to sync logger: %v", err)
		}
	}()

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("starting HTTP server")
		if err := app.StartHTTPServer(); err != nil {
			serverErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-serverErrors:
		logger.Fatal("server error", zap.Error(err))

	case sig := <-shutdown:
		logger.Info("received shutdown signal", zap.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := app.Shutdown(ctx); err != nil {
			logger.Error("error during shutdown", zap.Error(err))
			os.Exit(1)
		}

		logger.Info("server stopped gracefully")
	}
}
