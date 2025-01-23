package main

import (
	"github.com/histweety/go-common/connections"
)

func main() {
	connections.MongoConnect(connections.MongoConfig{
		URI: "mongodb://pos:password@localhost:27017",
		DB:  "pos",
	})
}
