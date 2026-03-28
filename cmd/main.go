package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"git.woda.ink/Woda_OA/config"
	"git.woda.ink/Woda_OA/internal/model"
	"git.woda.ink/Woda_OA/pkg/logger"
	"git.woda.ink/Woda_OA/router"
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

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	go func() {
		logger.L.Infof("server starting on %s", addr)
		if err := r.Run(addr); err != nil {
			logger.L.Fatalf("server failed to start: %v", err)
		}
	}()

	<-quit
	logger.L.Info("server shutting down...")
}
