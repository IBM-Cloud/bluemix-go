package main

import (
	"log"

	"github.com/IBM-Cloud/bluemix-go/api/cse/csev2"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {
	trace.Logger = trace.NewLogger("true")
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	cseClient, err := csev2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	seAPI := cseClient.ServiceEndpoints()

	// create a serviceendpoint
	payload := make(map[string]interface{})
	payload["service"] = "test-terrafor-11"
	payload["customer"] = "test-customer-11"
	payload["serviceAddresses"] = []string{"10.102.33.131", "10.102.33.133"}
	payload["region"] = "us-south"
	payload["dataCenters"] = []string{"dal10"}
	payload["tcpports"] = []int{8080, 80}

	log.Println("create a serviceendpoint with ", payload)

	newSrvId, err := seAPI.CreateServiceEndpoint(payload)
	if err != nil {
		log.Fatal(err)
	}
	// query the serviceendpoint
	log.Println("query the serviceendpoint ", newSrvId)

	srvObj, err := seAPI.GetServiceEndpoint(newSrvId)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(srvObj.Endpoints)

	log.Println("delete the service endpoint ", newSrvId)

	err = seAPI.DeleteServiceEndpoint(newSrvId)
	if err != nil {
		log.Fatal(err)
	}
}
