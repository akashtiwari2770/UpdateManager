package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"updatemanager/internal/api/router"
	"updatemanager/internal/service"
	"updatemanager/pkg/database"
)

func main() {
	ctx := context.Background()

	// 1. Load the default config (localhost)
	cfg := database.DefaultConfig()

	// 2. CHECK FOR CLOUD CONFIG (This is the fix!)
	// If Render provides a MONGODB_URI, we use that instead of localhost
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		cfg.URI = uri
	}
	if dbName := os.Getenv("MONGODB_DATABASE"); dbName != "" {
		cfg.Database = dbName
	}

	// 3. Now Connect
	db, err := database.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Disconnect(ctx)

	log.Println("Connected to MongoDB successfully")

	// Initialize services
	services := service.NewServiceFactory(db.Database)
	log.Println("Services initialized")

	// Setup router
	r := router.NewRouter(services)
	handler := r.Handler()

	// Start server
	port := os.Getenv("PORT") // Render sets this automatically to 10000 or similar
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}