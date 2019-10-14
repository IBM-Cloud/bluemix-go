package main

import (
	"flag"
	"fmt"
	"log"

	bluemix "github.com/IBM-Cloud/bluemix-go"

	//v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/session"

	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	c := new(bluemix.Config)

	var zone string
	flag.StringVar(&zone, "zone", "us-south-1", "Zone")

	// var privateVlan string
	// flag.StringVar(&privateVlan, "privateVlan", "", "Zone Private Vlan")

	// var publicVlan string
	// flag.StringVar(&publicVlan, "publicVlan", "", "Zone Public vlan")

	// var updatePrivateVlan string
	// flag.StringVar(&updatePrivateVlan, "updatePrivateVlan", "", "Zone Private vlan to be updated")

	// var updatePublicVlan string
	// flag.StringVar(&updatePublicVlan, "updatePublicVlan", "", "Zone Public vlan to be updated")

	var location string
	flag.StringVar(&location, "location", "dallas", "location")

	var region string
	flag.StringVar(&location, "region", "us-south", "region")

	// var skipDeletion bool
	// flag.BoolVar(&skipDeletion, "no-delete", false, "If provided will delete the resources created")

	flag.Parse()

	trace.Logger = trace.NewLogger("true")
	// if privateVlan == "" || publicVlan == "" || updatePrivateVlan == "" || updatePublicVlan == "" || zone == "" || location == "" {
	// 	flag.Usage()
	// 	os.Exit(1)
	// }

	target := v2.ClusterTargetHeader{}

	target.Region = region

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	clusterClient, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	clustersAPI := clusterClient.Clusters()

	out, err := clustersAPI.GetCluster("bmfghjbd0hi363tqv7c0", target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ouyt=", out)
}
