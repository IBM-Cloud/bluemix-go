package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/IBM-Cloud/bluemix-go/session"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	var clusterID string
	flag.StringVar(&clusterID, "clusterNameOrID", "", "cluster name or ID")

	var albID string
	flag.StringVar(&albID, "albID", "", "ALB ID")

	var region string
	flag.StringVar(&region, "region", "us-south", "region of cluster")

	flag.Parse()

	trace.Logger = trace.NewLogger("true")
	if clusterID == "" || albID == "" {
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

	alb, err := albAPI.GetALB(albID, target)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Get ALB with ID ", albID)
	fmt.Printf("alb: %+v", alb)
}
