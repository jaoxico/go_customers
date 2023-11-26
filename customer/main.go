package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	customvalidators "go_customers/customValidators"
	"go_customers/database"
	"math/rand"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Validate = validator.New(validator.WithRequiredStructEnabled())
var _ = Validate.RegisterValidation("isDate", customvalidators.IsAValidDate)

func RandomDate() time.Time {
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

const Collection string = "customer"

func InsertCustomer(client *mongo.Client, database string, newCustomer CreateCustomerPayload) (mongo.InsertOneResult, error) {
	if newCustomer.Dob == "" {
		newCustomer.Dob = RandomDate().Format(time.DateOnly)
	}
	insertResult, err := client.Database(database).Collection(Collection).InsertOne(context.TODO(), newCustomer)
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

	updateResult, err := client.Database(database).Collection(Collection).UpdateByID(context.TODO(), foundCustomer.Id, update)
	if err != nil {
		return mongo.UpdateResult{}, err
	}
	return *updateResult, nil
}

func DeleteCustomer(client *mongo.Client, database string, Id string) (mongo.DeleteResult, error) {
	_id, _ := primitive.ObjectIDFromHex(Id)
	filter := bson.D{{Key: "_id", Value: _id}}
	deleteResult, err := client.Database(database).Collection(Collection).DeleteOne(context.TODO(), filter)
	if err != nil {
		return mongo.DeleteResult{}, err
	}
	return *deleteResult, nil
}

func ListCustomers(client *mongo.Client, database string) ([]Customer, error) {
	var customers []Customer
	cursor, err := client.Database(database).Collection(Collection).Find(context.TODO(), bson.M{})
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

func HandlePost(response http.ResponseWriter, request *http.Request) {
	fmt.Println("incoming /customer POST")
	response.Header().Add("content-type", "application/json")
	client := database.GetMongoConnection()
	var customer CreateCustomerPayload
	err := json.NewDecoder(request.Body).Decode(&customer)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	err = Validate.Struct(customer)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		// if _, ok := err.(*validator.InvalidValidationError); ok {
		// 	fmt.Println(err)
		// 	return
		// }

		type failedFieldsStruct struct {
			Field   string
			Value   interface{}
			Message string
		}

		var failedFields []failedFieldsStruct

		for _, err := range err.(validator.ValidationErrors) {
			var failedField failedFieldsStruct
			failedField.Field = err.StructField()
			failedField.Value = err.Value()
			failedField.Message = err.Tag() + " - " + err.Param()
			failedFields = append(failedFields, failedField)
			fmt.Println(`Invalid field`, failedField)
		}
		response.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(response).Encode(failedFields)
		return
	}
	insertResult, err := InsertCustomer(client, database.Database, customer)
	if err != nil {
		fmt.Println(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	_ = json.NewEncoder(response).Encode(insertResult)
}

func HandleGetForList(response http.ResponseWriter, _ *http.Request) {
	fmt.Println("incoming /customer GET")
	response.Header().Add("content-type", "application/json")
	client := database.GetMongoConnection()
	var customers []Customer
	customers, err := ListCustomers(client, database.Database)
	if err != nil {
		fmt.Println(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	_ = json.NewEncoder(response).Encode(customers)
}

func HandleGetById(response http.ResponseWriter, request *http.Request) {
	fmt.Println("incoming /customer GET")
	response.Header().Add("content-type", "application/json")
	var Id = mux.Vars(request)["id"]
	_id, _ := primitive.ObjectIDFromHex(Id)
	client := database.GetMongoConnection()
	var foundCustomer Customer
	err := client.Database(database.Database).Collection(Collection).FindOne(context.TODO(), bson.D{{Key: "_id", Value: _id}}).Decode(&foundCustomer)
	if err != nil {
		fmt.Println(err)
		response.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(response).Encode(Id)
		return
	}
	response.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(response).Encode(foundCustomer)
}

func HandleDelete(response http.ResponseWriter, request *http.Request) {
	fmt.Println("incoming /customer DELETE")
	response.Header().Add("content-type", "application/json")
	var Id = mux.Vars(request)["id"]
	client := database.GetMongoConnection()
	deleteResult, err := DeleteCustomer(client, database.Database, Id)
	if err != nil {
		fmt.Println(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	_ = json.NewEncoder(response).Encode(deleteResult)
}

func HandlePatch(response http.ResponseWriter, request *http.Request) {
	fmt.Println("incoming /updateCustomer PATCH")
	response.Header().Add("content-type", "application/json")
	var updateCustomer UpdateCustomerPayload
	err := json.NewDecoder(request.Body).Decode(&updateCustomer)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	err = Validate.Struct(updateCustomer)
	if err != nil {
		type failedFieldsStruct struct {
			Field   string
			Value   interface{}
			Message string
		}

		var failedFields []failedFieldsStruct

		for _, err := range err.(validator.ValidationErrors) {
			var failedField failedFieldsStruct
			failedField.Field = err.StructField()
			failedField.Value = err.Value()
			failedField.Message = err.Tag() + " - " + err.Param()
			failedFields = append(failedFields, failedField)
			fmt.Println(`Invalid field`, failedField)
		}
		response.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(response).Encode(failedFields)
		return
	}
	var Id = mux.Vars(request)["id"]
	_id, _ := primitive.ObjectIDFromHex(Id)
	client := database.GetMongoConnection()

	var foundCustomer Customer
	err = client.Database(database.Database).Collection(Collection).FindOne(context.TODO(), bson.D{{Key: "_id", Value: _id}}).Decode(&foundCustomer)
	if err != nil {
		fmt.Println(err)
		response.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(response).Encode(Id)
		return
	}
	updateResult, err := UpdateCustomer(client, database.Database, foundCustomer, updateCustomer)
	if err != nil {
		fmt.Println(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	_ = json.NewEncoder(response).Encode(updateResult)
}
