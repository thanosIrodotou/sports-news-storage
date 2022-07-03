package mongodb

import (
	"context"
	"fmt"
	"time"

	"com.thanos/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoClient returns a mongo client
func NewMongoClient(cfg config.Mongo) (*mongo.Client, error) {
	mcOpts, err := clientOptions(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid mongodb client options: %w", err)
	}

	mClient, err := mongo.Connect(context.Background(), mcOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to init mongo client: %w", err)
	}

	err = mClient.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongo DB, client options: %v, error %w", mcOpts, err)
	}

	return mClient, nil
}

func clientOptions(cfg config.Mongo) (*options.ClientOptions, error) {
	cOpts := options.Client()

	// escapePassword := url.QueryEscape(cfg.Password)
	mongoUri := fmt.Sprintf(`mongodb://%s:%d/?tls=false`, cfg.Host, cfg.Port)

	cOpts.ApplyURI(mongoUri)
	cOpts.SetServerSelectionTimeout(10 * time.Second)
	cOpts.SetSocketTimeout(15 * time.Second)
	cOpts.SetMaxConnIdleTime(60 * time.Second)
	cOpts.SetRetryWrites(true)

	if err := cOpts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid mongodb client options: %w", err)
	}

	return cOpts, nil
}
