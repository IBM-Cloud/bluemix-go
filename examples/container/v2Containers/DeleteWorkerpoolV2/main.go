package main

import (
	"flag"
	"fmt"
	"log"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/session"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"

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

	err1 := workerpoolAPI.DeleteWorkerPool("bm64u3ed02o93vv36hb0", "bm64u3ed02o93vv36hb0-502aed1", target)

	//out,err=

	fmt.Println("err=", err1)

}
