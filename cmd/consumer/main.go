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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumerErrors := make(chan error, 1)

	go func() {
		if err := app.StartEventConsumers(ctx); err != nil {
			consumerErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-consumerErrors:
		logger.Fatal("consumer error", zap.Error(err))

	case sig := <-shutdown:
		logger.Info("received shutdown signal",
			zap.String("signal", sig.String()),
		)

		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := app.Shutdown(shutdownCtx); err != nil {
			logger.Error("error during shutdown", zap.Error(err))
			os.Exit(1)
		}

		logger.Info("consumer stopped gracefully")
	}
}
