package main

import (
	"flag"
	"fmt"
	"log"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"

	//v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
)

func main() {

	c := new(bluemix.Config)

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

	var zoneinfo = v2.WorkerPoolZone{
		Cluster:      "bmfgkjed0qgub4kab82g",
		Id:           "us-south-1",
		SubnetID:     "015ffb8b-efb1-4c03-8757-29335a07493b",
		WorkerPoolID: "bmfgkjed0qgub4kab82g-330d830",
	}
	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	target := v2.ClusterTargetHeader{}

	target.Region = region

	clusterClient, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	workerpoolAPI := clusterClient.WorkerPools()

	err2 := workerpoolAPI.CreateWorkerPoolZone(zoneinfo, target)

	fmt.Println("out=", err2)
}
