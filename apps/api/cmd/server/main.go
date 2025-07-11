package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jatifjr/app-unw-toefl/apps/api/internal/config"
	"github.com/jatifjr/app-unw-toefl/apps/api/internal/handler"
	"github.com/jatifjr/app-unw-toefl/apps/api/internal/repository"
	"github.com/jatifjr/app-unw-toefl/apps/api/internal/router"
	"github.com/jatifjr/app-unw-toefl/apps/api/internal/service"
)

func main() {
	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to listen for OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Start a goroutine to handle graceful shutdown
	go func() {
		<-signalChan
		log.Println("Shutdown signal received, exiting...")
		cancel()
	}()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	if cfg.Environment == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database connection
	db, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	scheduleRepo := repository.NewScheduleRepository(db)

	// Initialize services
	scheduleService := service.NewScheduleService(scheduleRepo)

	// Initialize handlers
	handlers := handler.NewHandler(scheduleService)

	// Initialize router
	r := router.NewRouter(handlers)
	engine := r.Setup()

	// Configure CORS
	engine.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigin := "*"

		// If we have specific allowed origins, check if the request origin is allowed
		if len(cfg.AllowedOrigins) > 0 && cfg.AllowedOrigins[0] != "*" {
			for _, allowed := range cfg.AllowedOrigins {
				if allowed == origin {
					allowedOrigin = origin
					break
				}
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Get port from environment variable or use default
	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}

	// Create a server with timeout configurations
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for cancel context
	<-ctx.Done()

	// Create a deadline for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}
