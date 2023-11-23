package main

import (
	"context"
	"encoding/json"
	"fmt"
	customvalidators "go_customers/customValidators"
	"go_customers/structs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoHost = os.Getenv(`BOOKING_MONGODB_URL`) // "mongodb://localhost:27017/?connect=direct"
var database = os.Getenv(`BOOKING_DB`)           // "booking"

func getContext() *context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return &ctx
}

func getMongoConnection() *mongo.Client {
	clientOpts := options.Client().ApplyURI(mongoHost)
	client, err := mongo.Connect(*getContext(), clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

var validate = validator.New(validator.WithRequiredStructEnabled())
var _ = validate.RegisterValidation("isDate", customvalidators.IsAValidDate)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Println("incoming / GET")
		_, err := response.Write([]byte("Testando!"))
		if err != nil {
			log.Fatal(err)
		}
		//		response.WriteHeader(http.StatusOK)
	}).Methods("GET")
	router.HandleFunc("/costumer", func(response http.ResponseWriter, request *http.Request) {
		fmt.Println("incoming /customer POST")
		response.Header().Add("content-type", "application/json")
		client := getMongoConnection()
		var customer structs.CreateCustomerPayload
		err := json.NewDecoder(request.Body).Decode(&customer)
		if err != nil {
			http.Error(response, err.Error(), http.StatusBadRequest)
			return
		}
		err = validate.Struct(customer)
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
		insertResult, err := structs.InsertCustomer(client, database, customer)
		if err != nil {
			fmt.Println(err.Error())
			response.WriteHeader(http.StatusInternalServerError)
			_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}
		_ = json.NewEncoder(response).Encode(insertResult)
	}).Methods("POST")
	router.HandleFunc("/costumer", func(response http.ResponseWriter, request *http.Request) {
		fmt.Println("incoming /customer GET")
		response.Header().Add("content-type", "application/json")
		client := getMongoConnection()
		var customers []structs.Customer
		customers, err := structs.ListCustomers(client, database)
		if err != nil {
			fmt.Println(err.Error())
			response.WriteHeader(http.StatusInternalServerError)
			_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}
		_ = json.NewEncoder(response).Encode(customers)
	}).Methods("GET")
	router.HandleFunc("/costumer/{id}", func(response http.ResponseWriter, request *http.Request) {
		fmt.Println("incoming /customer GET")
		response.Header().Add("content-type", "application/json")
		var Id = mux.Vars(request)["id"]
		_id, _ := primitive.ObjectIDFromHex(Id)
		client := getMongoConnection()
		var foundCustomer structs.Customer
		err := client.Database(database).Collection(structs.CustomerCollection).FindOne(context.TODO(), bson.D{{Key: "_id", Value: _id}}).Decode(&foundCustomer)
		if err != nil {
			fmt.Println(err)
			response.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(response).Encode(Id)
			return
		}
		response.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(response).Encode(foundCustomer)
	}).Methods("GET")
	router.HandleFunc("/costumer/{id}", func(response http.ResponseWriter, request *http.Request) {
		fmt.Println("incoming /customer DELETE")
		response.Header().Add("content-type", "application/json")
		var Id = mux.Vars(request)["id"]
		client := getMongoConnection()
		deleteResult, err := structs.DeleteCustomer(client, database, Id)
		if err != nil {
			fmt.Println(err.Error())
			response.WriteHeader(http.StatusInternalServerError)
			_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}
		_ = json.NewEncoder(response).Encode(deleteResult)
	}).Methods("DELETE")
	router.HandleFunc("/costumer/{id}", func(response http.ResponseWriter, request *http.Request) {
		fmt.Println("incoming /customer PATCH")
		response.Header().Add("content-type", "application/json")
		var customer structs.UpdateCustomerPayload
		err := json.NewDecoder(request.Body).Decode(&customer)
		if err != nil {
			http.Error(response, err.Error(), http.StatusBadRequest)
			return
		}
		err = validate.Struct(customer)
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
		var Id = mux.Vars(request)["id"]
		_id, _ := primitive.ObjectIDFromHex(Id)
		client := getMongoConnection()

		var foundCustomer structs.Customer
		err = client.Database(database).Collection(structs.CustomerCollection).FindOne(context.TODO(), bson.D{{Key: "_id", Value: _id}}).Decode(&foundCustomer)
		if err != nil {
			fmt.Println(err)
			response.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(response).Encode(Id)
			return
		}
		updateResult, err := structs.UpdateCustomer(client, database, foundCustomer, customer)
		if err != nil {
			fmt.Println(err.Error())
			response.WriteHeader(http.StatusInternalServerError)
			_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}
		_ = json.NewEncoder(response).Encode(updateResult)
	}).Methods("PATCH")
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	fmt.Println(`Listening on http://localhost:8000`)
	log.Fatal(srv.ListenAndServe())
}
