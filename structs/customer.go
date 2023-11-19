package structs

import (
	"context"
	"fmt"
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

type CreateCustomerPayload struct {
	FirstName string `json:"FirstName" bson:"FirstName" validate:"required"`
	LastName  string `json:"LastName" bson:"LastName" validate:"required"`
	Gender    string `json:"Gender" bson:"Gender" validate:"required,oneof=Masculino Feminino"`
	Dob       string `json:"Dob" bson:"Dob" validate:"required,isDate"`
}

type UpdateCustomerPayload struct {
	FirstName string `json:"FirstName" bson:"FirstName" validate:"omitempty"`
	LastName  string `json:"LastName" bson:"LastName" validate:"omitempty"`
	Gender    string `json:"Gender" bson:"Gender" validate:"omitempty,oneof=Masculino Feminino"`
	Dob       string `json:"Dob" bson:"Dob" validate:"omitempty,isDate"`
}

type Customer struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id" validate:"required"`
	FirstName string             `json:"FirstName" bson:"FirstName" validate:"required"`
	LastName  string             `json:"LastName" bson:"LastName" validate:"required"`
	Gender    string             `json:"Gender" bson:"Gender" validate:"required,oneof=Masculino Feminino"`
	Dob       string             `json:"Dob" bson:"Dob" validate:"required,isDate"`
}

const CustomerCollection string = "customer"

func InsertCustomer(client *mongo.Client, database string, newCustomer CreateCustomerPayload) (mongo.InsertOneResult, error) {
	if newCustomer.Dob == "" {
		newCustomer.Dob = randomDate().Format(time.DateOnly)
	}
	insertResult, err := client.Database(database).Collection(CustomerCollection).InsertOne(context.TODO(), newCustomer)
	if err != nil {
		return mongo.InsertOneResult{}, err
	}
	return *insertResult, nil
}

func UpdateCustomer(client *mongo.Client, database string, foundCustomer Customer, customer UpdateCustomerPayload) (mongo.UpdateResult, error) {
	if len(customer.Dob) == 0 {
		customer.Dob = foundCustomer.Dob
	}
	if len(customer.FirstName) == 0 {
		customer.FirstName = foundCustomer.FirstName
	}
	if len(customer.LastName) == 0 {
		customer.LastName = foundCustomer.LastName
	}
	if len(customer.Gender) == 0 {
		customer.Gender = foundCustomer.Gender
	}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "Dob", Value: customer.Dob},
		{Key: "FirstName", Value: customer.FirstName},
		{Key: "LastName", Value: customer.LastName},
		{Key: "Gender", Value: customer.Gender},
	}}}

	fmt.Println(update)

	updateResult, err := client.Database(database).Collection(CustomerCollection).UpdateByID(context.TODO(), foundCustomer.Id, update)
	if err != nil {
		return mongo.UpdateResult{}, err
	}
	return *updateResult, nil
}

func DeleteCustomer(client *mongo.Client, database string, Id string) (mongo.DeleteResult, error) {
	_id, _ := primitive.ObjectIDFromHex(Id)
	filter := bson.D{{Key: "_id", Value: _id}}
	deleteResult, err := client.Database(database).Collection(CustomerCollection).DeleteOne(context.TODO(), filter)
	if err != nil {
		return mongo.DeleteResult{}, err
	}
	return *deleteResult, nil
}

func ListCustomers(client *mongo.Client, database string) ([]Customer, error) {
	var customers []Customer
	cursor, err := client.Database(database).Collection(CustomerCollection).Find(context.TODO(), bson.M{})
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
