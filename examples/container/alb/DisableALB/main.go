package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM-Cloud/bluemix-go/session"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	var albID string
	flag.StringVar(&albID, "albID", "", "ALB ID")

	var region string
	flag.StringVar(&region, "region", "us-south", "region of cluster")

	flag.Parse()

	trace.Logger = trace.NewLogger("true")
	if albID == "" {
		flag.Usage()
		os.Exit(1)
	}

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	albClient, err := v1.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	albAPI := albClient.Albs()
	target := v1.ClusterTargetHeader{
		Region: region,
	}

	err = albAPI.DisableALB(albID, target)
	if err != nil {
		log.Fatal(err)
	}

	// wait...
	time.Sleep(30 * time.Second)

	alb, err := albAPI.GetALB(albID, target)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Get ALB with ID ", albID)
	fmt.Printf("alb: %+v", alb)
}
