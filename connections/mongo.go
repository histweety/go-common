package connections

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoConfig struct {
	URI string
	DB  string
}

func MongoConnect(cfg MongoConfig) *mongo.Database {
	var client *mongo.Client

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.URI)
	clientOpts.Auth.AuthSource = cfg.DB
	client, _ = mongo.Connect(clientOpts)
	err := client.Ping(ctx, nil)

	if err != nil {
		log.Fatal("[MongoDB]: failed to ping: ", err)
	}

	log.Info("[MongoDB]: connected!")

	return client.Database(cfg.DB)
}
