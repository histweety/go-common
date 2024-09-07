package connections

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Host   string
	DBName string
}

func MongoConnect(cfg MongoConfig) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.Host)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	log.Infof("Connected to MongoDB: %s", cfg.Host)
	return client, nil
}
