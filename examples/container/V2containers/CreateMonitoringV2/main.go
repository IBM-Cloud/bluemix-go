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

	var monitoringInfo = v2.MonitoringCreateRequest{
		Cluster: "DragonBoat-cluster",
		//IngestionKey:    "",
		SysidigInstance: "ec4f0886-edc4-409e-8720-574035538f91",
		//PrivateEndpoint: true,
	}

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

	out, err := monitoringAPI.CreateMonitoringConfig(monitoringInfo, target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("out=", out)
}
