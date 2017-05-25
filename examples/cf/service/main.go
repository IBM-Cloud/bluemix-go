package main

import (
	"flag"
	"log"
	"os"

	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/IBM-Bluemix/bluemix-go/helpers"
	"github.com/IBM-Bluemix/bluemix-go/session"
	"github.com/IBM-Bluemix/bluemix-go/trace"
)

func main() {
	var org string
	flag.StringVar(&org, "org", "", "Bluemix Organization")

	var space string
	flag.StringVar(&space, "space", "", "Bluemix Space")

	var name string
	flag.StringVar(&name, "name", "", "Service Instance Name")

	var newname string
	flag.StringVar(&newname, "newname", "", "Service Instance Name")

	var skipDeletion bool
	flag.BoolVar(&skipDeletion, "no-delete", false, "If provided will delete the resources created")

	flag.Parse()

	if org == "" || space == "" || name == "" || newname == "" {
		flag.Usage()
		os.Exit(1)
	}

	trace.Logger = trace.NewLogger("true")
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	client, err := cfv2.New(sess)

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
	createRequest := cfv2.ServiceInstanceCreateRequest{
		Name:      name,
		PlanGUID:  plan.GUID,
		SpaceGUID: myspace.GUID,
	}
	myService, err := serviceInstanceAPI.Create(createRequest)
	if err != nil {
		log.Fatal(err)
	}

	updateRequest := cfv2.ServiceInstanceUpdateRequest{
		Name:     helpers.String(newname),
		PlanGUID: helpers.String(plan.GUID),
	}
	updatedInstance, err := serviceInstanceAPI.Update(updateRequest, myService.Metadata.GUID)
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

	log.Println(myorg.GUID, myspace.GUID, plan.GUID, myService.Metadata.GUID, mykeys.Metadata.GUID, myRetrievedKeys)

	if !skipDeletion {
		err = serviceKeys.Delete(myRetrievedKeys.GUID)
		if err != nil {
			log.Fatal(err)
		}

		err = serviceInstanceAPI.Delete(myService.Metadata.GUID)
		if err != nil {
			log.Fatal(err)
		}

	}

}
