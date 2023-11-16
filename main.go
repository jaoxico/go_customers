package main

import (
	"booking-app/structs"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		_, err := response.Write([]byte("Testando!"))
		if err != nil {
			log.Fatal(err)
		}
		//		response.WriteHeader(http.StatusOK)
	}).Methods("GET")
	router.HandleFunc("/costumer", func(response http.ResponseWriter, request *http.Request) {
		response.Header().Add("content-type", "application/json")
		client := getMongoConnection()
		var customer structs.Customer
		err := json.NewDecoder(request.Body).Decode(&customer)
		if err != nil {
			http.Error(response, err.Error(), http.StatusBadRequest)
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
	router.HandleFunc("/costumers", func(response http.ResponseWriter, request *http.Request) {
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
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	fmt.Println(`Listening on http://localhost:8000`)
	log.Fatal(srv.ListenAndServe())
}
