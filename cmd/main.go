package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/imansohibul/otp-service/config"
	"github.com/imansohibul/otp-service/internal/handler"
)

// Author: MOCHAMAD SOHIBUL IMAN - iman@imansohibul.my.id

func main() {
	ctx := context.Background()

	restAPIServer, err := config.NewRestAPI()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize REST API server")
	}

	// Get server address from config or environment
	address := ":8080" // You should get this from your config

	// Graceful shutdown handler
	idleConnsClosed := make(chan struct{})
	go handleGracefulShutdown(ctx, restAPIServer, idleConnsClosed)

	log.Info().Msgf("Starting REST API server on %s...", address)
	if err := restAPIServer.Start(address); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("REST API server stopped with error")
	}

	<-idleConnsClosed
	log.Info().Msg("Server shut down gracefully")
}

func handleGracefulShutdown(ctx context.Context, restAPIServer *handler.RestAPIServer, done chan struct{}) {
	// Listen for interrupt signal (e.g., Ctrl+C, SIGTERM)
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	log.Warn().Msg("Shutdown signal received")

	// Create a context with timeout for shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := restAPIServer.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("Error during server shutdown")
	}

	close(done)
}
