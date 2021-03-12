package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	bluemix "github.com/IBM-Cloud/bluemix-go"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/session"

	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	c := new(bluemix.Config)

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

	trace.Logger = trace.NewLogger("true")

	target := v2.MonitoringTargetHeader{}

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	monitoringClient, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	monitoringAPI := monitoringClient.Monitoring()

	out, err := monitoringAPI.GetMonitoringConfig(cluster, InstanceID, target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("out=", out)
}
