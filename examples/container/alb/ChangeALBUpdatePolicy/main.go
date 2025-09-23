package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	bluemix "github.com/Mavrickk3/bluemix-go"
	v1 "github.com/Mavrickk3/bluemix-go/api/container/containerv1"
	"github.com/Mavrickk3/bluemix-go/session"
	"github.com/Mavrickk3/bluemix-go/trace"
)

func main() {

	var clusterID string
	flag.StringVar(&clusterID, "clusterNameOrID", "", "cluster name or ID")

	var enable bool
	flag.BoolVar(&enable, "enable", false, "turn on or off the ALB auto update")

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

	albUpdatePolicy := v1.ALBUpdatePolicy{
		AutoUpdate: enable,
	}

	fmt.Println("turn off auto update")
	updatePolicyErr := albAPIV1.ChangeALBUpdatePolicy(clusterID, albUpdatePolicy, targetV1)
	if updatePolicyErr != nil {
		log.Fatal(updatePolicyErr)
	}

	// wait...
	time.Sleep(30 * time.Second)

	autoUpdateConf, err := albAPIV1.GetALBUpdatePolicy(clusterID, targetV1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("autoUpdateConf: %+v", autoUpdateConf)

}
