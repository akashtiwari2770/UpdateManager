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

	// Connect to MongoDB
	cfg := database.DefaultConfig()
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
	port := os.Getenv("PORT")
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
		log.Printf("API available at http://localhost:%s/api/v1", port)
		log.Printf("Health check at http://localhost:%s/health", port)
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
