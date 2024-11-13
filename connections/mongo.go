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
	URI string
	DB  string
}

func MongoConnect(cfg MongoConfig) *mongo.Database {
	var client *mongo.Client
	var err interface{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.URI)

	for i := 0; i < 3; i++ {
		log.Infof("[MongoDB]: connecting...")
		client, err = mongo.Connect(ctx, clientOpts)
		if err == nil {
			break
		}

		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("[MongoDB]: failed to connect: %v", err)
		os.Exit(0)
	}

	log.Infof("[MongoDB]: connected!")

	return client.Database(cfg.DB)
}
