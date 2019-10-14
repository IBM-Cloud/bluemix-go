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

	var zone string
	flag.StringVar(&zone, "zone", "us-south-1", "Zone")
	flag.Parse()

	trace.Logger = trace.NewLogger("true")

	var albinfo = v2.AlbCreateReq{
		Cluster:         "bm64u3ed02o93vv36hb0",
		EnableByDefault: true,
		Type:            "private",
		ZoneAlb:         "us-south-1",
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

	albAPI := clusterClient.Albs()

	err2 := albAPI.CreateAlb(albinfo, target)

	fmt.Println("err=", err2)

}
