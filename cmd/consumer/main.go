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

	// Contexto base para los consumidores
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Canal para errores de los consumidores
	consumerErrors := make(chan error, 1)

	// Iniciar consumidores en goroutine
	go func() {
		if err := app.StartEventConsumers(ctx); err != nil {
			consumerErrors <- err
		}
	}()

	// Canal para señales del sistema operativo
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Bloquear hasta recibir señal de apagado o error
	select {
	case err := <-consumerErrors:
		logger.Fatal("consumer error", zap.Error(err))

	case sig := <-shutdown:
		logger.Info("received shutdown signal",
			zap.String("signal", sig.String()),
		)

		// Cancelar contexto de consumidores
		cancel()

		// Contexto con timeout para graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		// Apagar aplicación de forma ordenada
		if err := app.Shutdown(shutdownCtx); err != nil {
			logger.Error("error during shutdown", zap.Error(err))
			os.Exit(1)
		}

		logger.Info("consumer stopped gracefully")
	}
}
