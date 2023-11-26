package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var mongoHost string

var Database string

func getContext() *context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return &ctx
}

func GetMongoConnection() *mongo.Client {
	mongoHost = os.Getenv(`GO_CUSTOMERS_MONGODB_URL`) // "mongodb://localhost:27017/?connect=direct"
	Database = os.Getenv(`GO_CUSTOMERS_DB`)           // "go_customers"
	clientOpts := options.Client().ApplyURI(mongoHost)
	client, err := mongo.Connect(*getContext(), clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
