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
	log.Printf("[Config] Loaded from %s", *configPath)

	// Initialize MySQL
	if err := model.InitDB(&cfg.Database); err != nil {
		log.Fatalf("[DB] %v", err)
	}
	defer model.CloseDB()

	// Initialize Redis
	if err := model.InitRedis(&cfg.Redis); err != nil {
		log.Printf("[Redis] Warning: %v (continuing without redis)", err)
	} else {
		defer model.CloseRedis()
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
		log.Printf("[Server] Woda OA starting on %s", addr)
		if err := r.Run(addr); err != nil {
			log.Fatalf("[Server] Failed to start: %v", err)
		}
	}()

	<-quit
	log.Println("[Server] Shutting down...")
}
