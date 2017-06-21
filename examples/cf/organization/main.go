package main

import (
	"flag"
	"log"
	"os"

	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/IBM-Bluemix/bluemix-go/session"
	"github.com/IBM-Bluemix/bluemix-go/trace"
)

func main() {
	var org string
	flag.StringVar(&org, "org", "", "Bluemix Organization")

	var neworg string
	flag.StringVar(&neworg, "neworg", "", "Bluemix Organization")

	flag.Parse()

	if org == "" || neworg == "" {
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

	//  Authorization required to create organization.
	//	err = orgAPI.Create(org)
	//	if err != nil {
	//		log.Fatal(err)
	//	}

	region := sess.Config.Region
	myorg, err := orgAPI.FindByName(org, region)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(myorg.GUID, myorg.Name)

	err = orgAPI.Update(myorg.GUID, neworg)
	if err != nil {
		log.Fatal(err)
	}

	updatedOrg, err := orgAPI.FindByName(neworg, region)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(updatedOrg.GUID, updatedOrg.Name)

	//  Authorization required to delete organization.
	//	err = orgAPI.Delete(updatedOrg.GUID, true)
	//	if err != nil {
	//		log.Fatal(err)
	//	}

	_, err = orgAPI.List()
	if err != nil {
		log.Fatal(err)
	}
}
