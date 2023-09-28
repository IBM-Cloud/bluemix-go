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

	var clusterID string
	flag.StringVar(&clusterID, "clusterNameOrID", "", "cluster name or ID")

	var albID string
	flag.StringVar(&albID, "albID", "", "ALB ID")

	var enable bool
	flag.BoolVar(&enable, "enable", false, "enable alb")

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

	albConf := v1.ALBConfig{
		ALBID:     albID,
		ClusterID: clusterID,
		Enable:    enable,
	}

	err = albAPI.EnableALB(albID, albConf, target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ALB enabled")

	alb, err := albAPI.GetALB(albID, target)
	if err != nil {
		log.Fatal(err)
	}

	// wait...
	time.Sleep(30 * time.Second)

	fmt.Println("Get ALB with ID ", albID)
	fmt.Printf("alb: %+v", alb)
}
