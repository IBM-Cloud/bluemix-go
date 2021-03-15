package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/session"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"

	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	var cluster string
	flag.StringVar(&cluster, "cluster", "", "Clusetr Name")

	var InstanceID string
	flag.StringVar(&InstanceID, "InstanceID", "", " monitoring InstanceID")
	flag.Parse()

	trace.Logger = trace.NewLogger("true")
	if cluster == "" || InstanceID == "" {
		flag.Usage()
		os.Exit(1)
	}

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
		Cluster:  cluster,
		Instance: InstanceID,
	}

	resp, err1 := monitoringAPI.DeleteMonitoringConfig(monitoringInfo, target)
	if err1 != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted the monitor instance successfully", resp)

}
