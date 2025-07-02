package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/TicketsBot-cloud/status-updates/internal/config"
	"github.com/TicketsBot-cloud/status-updates/internal/daemon"
	"github.com/TicketsBot-cloud/status-updates/internal/db"
	"github.com/TicketsBot-cloud/status-updates/internal/httpserver"
	"github.com/TicketsBot-cloud/status-updates/internal/statuspage"
	"go.uber.org/zap"

	_ "github.com/joho/godotenv/autoload" // Load environment variables from .env file
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	config.Conf = conf
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if err := db.InitDB(); err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	statusPageClient := statuspage.NewClient(logger)
	daemon := daemon.NewDaemon(logger, statusPageClient)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := daemon.Start(ctx); err != nil {
			logger.Error("Failed to start daemon", zap.Error(err))
		}
	}()

	server := httpserver.NewServer(logger.With(zap.String("component", "server")), conf)

	go func() {
		logger.Info("Starting HTTP server on :8080")
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down HTTP server...")
}
