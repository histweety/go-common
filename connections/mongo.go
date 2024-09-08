package connections

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Host string
}

func MongoConnect(cfg MongoConfig) *mongo.Client {
	var client *mongo.Client
	var err interface{}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.Host)

	for i := 0; i < 3; i++ {
		client, err = mongo.Connect(ctx, clientOpts)
		if err == nil {
			break
		}

		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		os.Exit(0)
	}

	log.Infof("Connected to MongoDB: %s", cfg.Host)

	return client
}
