package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"
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

	targetV1 := v1.ClusterTargetHeader{}

	clusterClientV1, err := v1.New(sess)
	if err != nil {
		log.Fatal(err)
	}

	albAPIV1 := clusterClientV1.Albs()

	err = albAPIV1.UpdateALBs(clusterID, targetV1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("forced one-time update for ALBs")

}
