package main

import (
	"flag"
	"log"
	"os"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/api/resource/resourcev2/controllerv2"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	var name string
	flag.StringVar(&name, "name", "", "Name of the service-instance")

	var serviceInstanceID string
	flag.StringVar(&serviceInstanceID, "serviceInstanceID", "", "guid of instance")

	flag.Parse()

	if name == "" || serviceInstanceID == "" {
		flag.Usage()
		os.Exit(1)
	}

	trace.Logger = trace.NewLogger("true")
	sess, err := session.New(&bluemix.Config{Debug: true})
	if err != nil {
		log.Fatal(err)
	}

	controllerClient, err := controllerv2.New(sess)

	if err != nil {
		log.Fatal(err)
	}

	resServiceInstanceAPI := controllerClient.ResourceServiceInstanceV2()

	query := controllerv2.ServiceInstanceQuery{
		Name: name,
	}
	listInsatnce, err := resServiceInstanceAPI.ListInstances(query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Resoure  Instance List :", listInsatnce)

	serviceInstance, err := resServiceInstanceAPI.GetInstance(serviceInstanceID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("\nResoure service Instance Details :", serviceInstance)

	// serviceInstanceUpdatePayload := controller.UpdateServiceInstanceRequest{
	// 	Name:          "update-service",
	// 	ServicePlanID: servicePlanID,
	// }

	// serviceInstance, err = resServiceInstanceAPI.UpdateInstance(serviceInstance.ID, serviceInstanceUpdatePayload)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Resoure service Instance Details after update :", serviceInstance)

	// err = resServiceInstanceAPI.DeleteInstance(serviceInstance.ID, true)

	// if err != nil {
	// 	log.Fatal(err)
	// }

}
