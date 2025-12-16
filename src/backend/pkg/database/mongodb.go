package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB represents a MongoDB database connection
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// Config holds MongoDB connection configuration
type Config struct {
	URI      string
	Database string
	Timeout  time.Duration
}

// DefaultConfig returns default MongoDB configuration
func DefaultConfig() *Config {
	return &Config{
		URI:      "mongodb://admin:admin123@localhost:27017/updatemanager?authSource=admin",
		Database: "updatemanager",
		Timeout:  10 * time.Second,
	}
}

// Connect establishes a connection to MongoDB
func Connect(ctx context.Context, cfg *Config) (*MongoDB, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Set connection timeout
	connectCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	// Create client options
	clientOptions := options.Client().ApplyURI(cfg.URI)

	// Connect to MongoDB
	client, err := mongo.Connect(connectCtx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify connection
	if err := client.Ping(connectCtx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(cfg.Database)

	return &MongoDB{
		Client:   client,
		Database: db,
	}, nil
}

// Disconnect closes the MongoDB connection
func (m *MongoDB) Disconnect(ctx context.Context) error {
	if m.Client != nil {
		return m.Client.Disconnect(ctx)
	}
	return nil
}

// Collection returns a MongoDB collection
func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}
