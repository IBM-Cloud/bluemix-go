package main

import (
	"flag"
	"fmt"
	"log"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
)

func main() {
	var Flavor string
	flag.StringVar(&Flavor, "flavor", "", "Flavor")
	var Zone string
	flag.StringVar(&Zone, "zone", "us-south-1", "Zone")
	var HostPoolID string
	flag.StringVar(&HostPoolID, "hostpoolid", "", "HostPoolID")
	flag.Parse()
	fmt.Println("[FLAG]Flavor: ", Flavor)
	fmt.Println("[FLAG]Zone: ", Zone)
	fmt.Println("[FLAG]HostPoolID: ", HostPoolID)

	c := new(bluemix.Config)

	trace.Logger = trace.NewLogger("true")

	var createDedicatedHost = v2.CreateDedicatedHostRequest{
		Flavor:     Flavor,
		HostPoolID: HostPoolID,
		Zone:       Zone,
	}
	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	target := v2.ClusterTargetHeader{}

	v2Client, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	dedicatedHostAPI := v2Client.DedicatedHost()

	dh, err := dedicatedHostAPI.CreateDedicatedHost(createDedicatedHost, target)
	if err != nil {
		fmt.Printf("Create was not successful: %v \n", err)
		return
	}
	fmt.Printf("Create dedicated host response: %v \n", dh)
}
