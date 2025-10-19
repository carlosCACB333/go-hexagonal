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
			// Ignorar errores de sync en algunos sistemas
			log.Printf("Failed to sync logger: %v", err)
		}
	}()

	// Canal para errores del servidor
	serverErrors := make(chan error, 1)

	// Iniciar servidor HTTP en goroutine
	go func() {
		logger.Info("starting HTTP server")
		if err := app.StartHTTPServer(); err != nil {
			serverErrors <- err
		}
	}()

	// Canal para señales del sistema operativo
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Bloquear hasta recibir señal de apagado o error
	select {
	case err := <-serverErrors:
		logger.Fatal("server error", zap.Error(err))

	case sig := <-shutdown:
		logger.Info("received shutdown signal",
			zap.String("signal", sig.String()),
		)

		// Contexto con timeout para graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Apagar aplicación de forma ordenada
		if err := app.Shutdown(ctx); err != nil {
			logger.Error("error during shutdown", zap.Error(err))
			os.Exit(1)
		}

		logger.Info("server stopped gracefully")
	}
}
