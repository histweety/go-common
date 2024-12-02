package connections

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URI string
	DB  string
}

func MongoConnect(cfg MongoConfig) *mongo.Database {
	var client *mongo.Client
	var err interface{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.URI)
	client, err = mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("[MongoDB]: failed to connect:", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal("[MongoDB]: failed to disconnect client: ", err)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("[MongoDB]: failed to ping:", err)
	}

	log.Info("[MongoDB]: connected!")

	return client.Database(cfg.DB)
}
