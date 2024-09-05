package mongo

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Config *MongoConfig
}

func New(config *MongoConfig) *MongoInstance {
	return &MongoInstance{
		Config: config,
	}
}

func (m *MongoInstance) Connect() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(m.Config.Host)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	log.Infof("Connected to MongoDB: %s", m.Config.Host)
	return client, nil
}
