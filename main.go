package main

import (
	"fmt"
	"github.com/joho/godotenv"
	customvalidators "go_customers/customValidators"
	"go_customers/customer"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var Validate = validator.New(validator.WithRequiredStructEnabled())
var _ = Validate.RegisterValidation("isDate", customvalidators.IsAValidDate)

func main() {

	var environment = os.Getenv("GO_CUSTOMERS_ENVIRONMENT")
	if environment == "" {
		environment = "local"
	}

	err := godotenv.Load("environment/.env." + environment)

	if err != nil {
		fmt.Println(err)
		log.Fatal()
	}

	router := mux.NewRouter()
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Println("incoming / GET")
		_, err := response.Write([]byte("Testando!"))
		if err != nil {
			log.Fatal(err)
		}
		//		response.WriteHeader(http.StatusOK)
	}).Methods("GET")
	router.HandleFunc("/costumer", customer.HandlePost).Methods("POST")
	router.HandleFunc("/costumer", customer.HandleGetForList).Methods("GET")
	router.HandleFunc("/costumer/{id}", customer.HandleGetById).Methods("GET")
	router.HandleFunc("/costumer/{id}", customer.HandleDelete).Methods("DELETE")
	router.HandleFunc("/costumer/{id}", customer.HandlePatch).Methods("PATCH")
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	fmt.Println(`Listening on http://localhost:8000`)
	log.Fatal(srv.ListenAndServe())
}
