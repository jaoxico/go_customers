package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func getContext() *context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return &ctx
}

func GetMongoConnection(mongoHost string) *mongo.Client {
	clientOpts := options.Client().ApplyURI(mongoHost)
	client, err := mongo.Connect(*getContext(), clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
