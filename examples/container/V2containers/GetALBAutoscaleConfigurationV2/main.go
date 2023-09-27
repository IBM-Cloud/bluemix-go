package main

import (
	"fmt"
	"log"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
)

func main() {

	c := new(bluemix.Config)

	trace.Logger = trace.NewLogger("true")

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	clusterID := "ck6k27hd0s1542093c6g"
	albID := "public-crck6k27hd0s1542093c6g-alb1"

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
