package structs

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Ticket struct {
	Customer primitive.ObjectID
	When     time.Time
	Order    uint8
}

const ticketCollection string = "ticket"

func OrderTicket(client mongo.Client, database string, newOrder Ticket) (mongo.InsertOneResult, error) {
	insertResult, err := client.Database(database).Collection(ticketCollection).InsertOne(context.TODO(), newOrder)
	if err != nil {
		return mongo.InsertOneResult{}, err
	}
	return *insertResult, nil
}
