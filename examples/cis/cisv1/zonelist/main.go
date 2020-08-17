package main

import (
	"flag"
	"log"
	"os"

	"github.com/IBM-Cloud/bluemix-go/api/cis/cisv1"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	var cisID string
	flag.StringVar(&cisID, "cisID", "", "CRN of the CIS service instance")

	// var domain string
	// flag.StringVar(&domain, "domain", "", "DNS domain name for zone")

	// flag.Parse()

	if cisID == "" {
		flag.Usage()
		os.Exit(1)
	}

	trace.Logger = trace.NewLogger("true")
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	cisClient, err := cisv1.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	zonesAPI := cisClient.Zones()

	myZonePtr, err := zonesAPI.ListZones(cisID)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("myZonePtr", myZonePtr)

}
