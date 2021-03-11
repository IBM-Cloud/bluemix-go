package main

import (
	"fmt"
	"log"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/session"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"

	"github.com/IBM-Cloud/bluemix-go/trace"
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

	target := v2.MonitoringTargetHeader{}

	monitoringClient, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	monitoringAPI := monitoringClient.Monitoring()

	var monitoringInfo = v2.MonitoringDeleteRequest{
		Cluster:  "vpc-cluster2",
		Instance: "b22d551c-238b-4849-9b16-63edf1e45e7d",
	}

	err1 := monitoringAPI.DeleteMonitoringConfig(monitoringInfo, target)
	if err1 != nil {
		log.Fatal(err)
	}
	fmt.Println("out=", err1)

}
