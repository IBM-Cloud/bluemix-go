package main

import (
	"flag"
	"log"
	"os"

	"github.com/IBM-Bluemix/bluemix-go/api/mccp/mccpv2"
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

	region := sess.Config.Region
	client, err := mccpv2.New(sess)

	if err != nil {
		log.Fatal(err)
	}

	orgAPI := client.Organizations()

	err = orgAPI.Create(org)
	if err != nil {
		log.Fatal(err)
	}

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

	err = orgAPI.Delete(updatedOrg.GUID, true)
	if err != nil {
		log.Fatal(err)
	}

	_, err = orgAPI.List(region)
	if err != nil {
		log.Fatal(err)
	}
}
