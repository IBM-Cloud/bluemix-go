package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	var clusterID string
	flag.StringVar(&clusterID, "clusterNameOrID", "", "cluster name or ID")

	var region string
	flag.StringVar(&region, "region", "us-south", "region of cluster")

	flag.Parse()

	trace.Logger = trace.NewLogger("true")
	if clusterID == "" {
		flag.Usage()
		os.Exit(1)
	}

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	addOnClient, err := v1.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	addOnAPI := addOnClient.AddOns()
	target := v1.ClusterTargetHeader{
		Region: region,
	}

	//Enable the AddOns

	addOnConfig := v1.ConfigureAddOns{
		AddonsList: []v1.AddOn{},
		Enable:     false,
		Update:     false,
	}
	var addOn = v1.AddOn{
		Name: "istio",
	}
	addOnConfig.AddonsList = append(addOnConfig.AddonsList, addOn)
	_, err = addOnAPI.ConfigureAddons(clusterID, &addOnConfig, target)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	addons, err := addOnAPI.GetAddons(clusterID, target)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("The avalable addons in a cluster ", addons)

}
