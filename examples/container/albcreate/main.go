package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM-Cloud/bluemix-go/session"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	var enable bool
	flag.BoolVar(&enable, "enable", false, "enable alb")

	var clusterID string //mandatory
	flag.StringVar(&clusterID, "clusterNameOrID", "", "cluster name or ID - mandatory")

	var albtype string //mandatory
	flag.StringVar(&albtype, "type", "", "type of alb - mandatory")

	var vlanID string //mandatory
	flag.StringVar(&vlanID, "vlanID", "", "vlanID of alb - mandatory")

	var zone string //mandatory
	flag.StringVar(&vlanID, "zone", "", "zone of alb - mandatory")

	var region string //mandatory
	flag.StringVar(&region, "region", "us-south", "region of cluster - mandatory")

	flag.Parse()

	trace.Logger = trace.NewLogger("true")
	if region == "" || albtype == "" || clusterID == "" || vlanID == "" || zone == "" {
		flag.Usage()
		os.Exit(1)
	}

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	albClient, err := v1.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	albAPI := albClient.Albs()
	target := v1.ClusterTargetHeader{
		Region: region,
	}

	params := v1.CreateALB{
		Zone: "testZone", VlanID: "testVlan", Type: "testType", EnableByDefault: true, IP: "1.2.3.4", NLBVersion: "testnlbVersion", IngressImage: "testingressImage",
	}

	albResp, err := albAPI.CreateALB(params, clusterID, target)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	alb, err := albAPI.GetALB(albResp.Alb, target)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created ALB: ", alb)

}
