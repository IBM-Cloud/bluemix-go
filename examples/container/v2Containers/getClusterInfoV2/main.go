package main

import (
	"flag"
	"fmt"
	"log"

	bluemix "github.com/IBM-Cloud/bluemix-go"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/session"

	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	c := new(bluemix.Config)

	var zone string
	flag.StringVar(&zone, "zone", "us-south-1", "Zone")
	var location string
	flag.StringVar(&location, "location", "dallas", "location")

	var region string
	flag.StringVar(&location, "region", "us-south", "region")

	flag.Parse()

	trace.Logger = trace.NewLogger("true")

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
