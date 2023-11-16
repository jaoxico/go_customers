package structs

import (
	"context"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func randomDate() time.Time {
	minDate := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	maxDate := time.Now().Unix()
	delta := maxDate - minDate

	sec := rand.Int63n(delta) + minDate
	return time.Unix(sec, 0)
}

type Customer struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"FirstName,omitempty" bson:"FirstName,omitempty"`
	LastName  string             `json:"LastName,omitempty" bson:"LastName,omitempty"`
	Gender    string             `json:"Gender,omitempty" bson:"Gender,omitempty"`
	Dob       string             `json:"Dob,omitempty" bson:"Dob,omitempty"`
}

const customerCollection string = "customer"

func InsertCustomer(client *mongo.Client, database string, newCustomer Customer) (mongo.InsertOneResult, error) {
	if newCustomer.Dob == "" {
		newCustomer.Dob = randomDate().Format(time.DateOnly)
	}
	insertResult, err := client.Database(database).Collection(customerCollection).InsertOne(context.TODO(), newCustomer)
	if err != nil {
		return mongo.InsertOneResult{}, err
	}
	return *insertResult, nil
}

func ListCustomers(client *mongo.Client, database string) ([]Customer, error) {
	var customers []Customer
	cursor, err := client.Database(database).Collection(customerCollection).Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(context.TODO())
	}()
	for cursor.Next(context.TODO()) {
		var customer Customer
		_ = cursor.Decode(&customer)
		customers = append(customers, customer)
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	return customers, nil
}
