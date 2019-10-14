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

	var poolinfo = v2.WorkerPoolRequest{
		Cluster: "bm64u3ed02o93vv36hb0",
		WorkerPoolConfig: v2.WorkerPoolConfig{
			Flavor:      "c2.2x4",
			Name:        "mywork2",
			VpcID:       "6015365a-9d93-4bb4-8248-79ae0db2dc26",
			WorkerCount: 1,
			Zones:       []v2.Zone{},
		},
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

	out, err := workerpoolAPI.CreateWorkerPool(poolinfo, target)

	fmt.Println("out=", out)
}
