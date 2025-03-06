package connections

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoConfig struct {
	URI        string
	DB         string
	AuthSource string
	Username   string
	Password   string
}

func MongoConnect(cfg MongoConfig) *mongo.Database {
	var client *mongo.Client

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		AuthSource:    cfg.AuthSource,
		Username:      cfg.Username,
		Password:      cfg.Password,
	}
	clientOpts := options.Client().ApplyURI(cfg.URI).SetAuth(credential)
	client, _ = mongo.Connect(clientOpts)
	err := client.Ping(ctx, nil)

	if err != nil {
		log.Fatal("[MongoDB]: failed to ping: ", err)
	}

	log.Info("[MongoDB]: connected!")

	return client.Database(cfg.DB)
}
