package customer

import (
	"github.com/benweissmann/memongo"
	"log"
	"testing"
)

func TestInsertCustomer(t *testing.T) {
	mongoServer, err := memongo.Start("4.0.5")
	if err != nil {
		t.Fatal(err)
	}
	//	databaseName := "go_customers"
	defer mongoServer.Stop()
	logger := log.Default()
	logger.Println(mongoServer.URI())
	//	client := database.GetMongoConnection(mongoServer.URI())
	//	result, err := InsertCustomer(client, databaseName, CreateCustomerPayload{
	//		FirstName: "João",
	//		LastName:  "Souza",
	//		Gender:    "Masculino",
	//		Dob:       "1973-06-10",
	//	})
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//	fmt.Println(result)
}
