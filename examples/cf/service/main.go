package main

import (
	"flag"
	"log"
	"os"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/trace"
	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/IBM-Bluemix/bluemix-go/session"
)

func main() {
	var org string
	flag.StringVar(&org, "org", "", "Bluemix Organization")

	var space string
	flag.StringVar(&space, "space", "", "Bluemix Space")

	flag.Parse()

	if org == "" || space == "" {
		flag.Usage()
		os.Exit(1)
	}

	trace.Logger = trace.NewLogger("true")
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	client, err := cfv2.NewClient(sess)

	if err != nil {
		log.Fatal(err)
	}

	orgAPI := client.Organizations()
	myorg, err := orgAPI.FindByName(org)

	if err != nil {
		log.Fatal(err)
	}

	spaceAPI := client.Spaces()
	myspace, err := spaceAPI.FindByNameInOrg(myorg.GUID, space)

	if err != nil {
		log.Fatal(err)
	}

	serviceOfferingAPI := client.ServiceOfferings()
	myserviceOff, err := serviceOfferingAPI.FindByLabel("cloudantNoSQLDB")
	if err != nil {
		log.Fatal(err)
	}
	servicePlanAPI := client.ServicePlans()
	plan, err := servicePlanAPI.GetServicePlan(myserviceOff.GUID, "Lite")
	if err != nil {
		log.Fatal(err)
	}

	serviceInstanceAPI := client.ServiceInstances()
	myService, err := serviceInstanceAPI.Create("myservice", plan.GUID, myspace.GUID, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	updatedInstance, err := serviceInstanceAPI.Update("New instance", myService.Metadata.GUID, plan.GUID, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	serviceKeys := client.ServiceKeys()
	mykeys, err := serviceKeys.Create(updatedInstance.Metadata.GUID, "mykey", nil)
	if err != nil {
		log.Fatal(err)
	}

	myRetrievedKeys, err := serviceKeys.FindByName(myService.Metadata.GUID, "mykey")
	if err != nil {
		log.Fatal(err)
	}

	err = serviceKeys.Delete(myRetrievedKeys.GUID)
	if err != nil {
		log.Fatal(err)
	}

	err = serviceInstanceAPI.Delete(myService.Metadata.GUID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(myorg.GUID, myspace.GUID, plan.GUID, myService.Metadata.GUID, mykeys.Metadata.GUID, myRetrievedKeys)
}
