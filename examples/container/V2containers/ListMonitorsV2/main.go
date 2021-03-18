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

	c := new(bluemix.Config)

	var cluster string
	flag.StringVar(&cluster, "cluster", "", "Clusetr Name")
	flag.Parse()

	trace.Logger = trace.NewLogger("true")
	if cluster == "" {
		flag.Usage()
		os.Exit(1)
	}

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	target := v2.MonitoringTargetHeader{}

	monitoringClient, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	monitoringAPI := monitoringClient.Monitoring()

	out, err := monitoringAPI.ListAllMonitors(cluster, target)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Monitors ", out)
}
