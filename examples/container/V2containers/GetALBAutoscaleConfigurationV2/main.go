package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	bluemix "github.com/Mavrickk3/bluemix-go"
	"github.com/Mavrickk3/bluemix-go/session"
	"github.com/Mavrickk3/bluemix-go/trace"

	v2 "github.com/Mavrickk3/bluemix-go/api/container/containerv2"
)

func main() {

	var albID, clusterID string
	flag.StringVar(&albID, "albID", "", "ALB ID")
	flag.StringVar(&clusterID, "clusterNameOrID", "", "cluster name or ID")
	flag.Parse()

	trace.Logger = trace.NewLogger("true")

	if clusterID == "" || albID == "" {
		flag.Usage()
		os.Exit(1)
	}

	c := new(bluemix.Config)

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	target := v2.ClusterTargetHeader{}

	clusterClient, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}

	albAPI := clusterClient.Albs()

	getAutoscaleConf, err := albAPI.GetALBAutoscaleConfiguration(clusterID, albID, target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("getAutoscaleConf.Config: %+v\n", getAutoscaleConf.Config)
}
