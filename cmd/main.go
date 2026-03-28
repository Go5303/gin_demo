package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Go5303/gin_demo/config"
	"github.com/Go5303/gin_demo/internal/model"
	"github.com/Go5303/gin_demo/pkg/logger"
	"github.com/Go5303/gin_demo/router"
	"github.com/gin-gonic/gin"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config/config.yaml", "path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("[Config] Failed to load config: %v", err)
	}

	// Initialize logger
	logger.Init(cfg.Log.Level, cfg.Log.ErrorPath+"/app.log")
	defer logger.L.Sync()
	logger.L.Infof("config loaded from %s", *configPath)

	// Initialize MySQL
	if err := model.InitDB(&cfg.Database); err != nil {
		logger.L.Fatalf("[DB] %v", err)
	}
	defer model.CloseDB()
	logger.L.Info("MySQL connected")

	// Initialize Redis
	if err := model.InitRedis(&cfg.Redis); err != nil {
		logger.L.Warnf("Redis init failed: %v (continuing without redis)", err)
	} else {
		defer model.CloseRedis()
		logger.L.Info("Redis connected")
	}

	// Set Gin mode
	switch cfg.Server.Mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin engine
	r := gin.New()

	// Setup routes
	router.Setup(r, cfg)

	// Create HTTP server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Start server in goroutine
	go func() {
		logger.L.Infof("server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L.Fatalf("server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.L.Info("server shutting down...")

	// Give outstanding requests 10 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.L.Errorf("server forced to shutdown: %v", err)
	}

	logger.L.Info("server exited")
}
