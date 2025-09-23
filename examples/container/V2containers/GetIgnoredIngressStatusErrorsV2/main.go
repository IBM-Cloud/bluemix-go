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

	var clusterID string
	flag.StringVar(&clusterID, "clusterNameOrID", "", "cluster name or ID")
	flag.Parse()

	trace.Logger = trace.NewLogger("true")

	if clusterID == "" {
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

	ignoredErrors, getErr := albAPI.GetIgnoredIngressStatusErrors(clusterID, target)
	fmt.Println("err: ", getErr)
	fmt.Printf("ignoredErrors: %+v\n", ignoredErrors.IgnoredErrors)
}
