package main

import (
	"flag"
	"fmt"
	"log"

	bluemix "github.com/Mavrickk3/bluemix-go"
	"github.com/Mavrickk3/bluemix-go/session"
	"github.com/Mavrickk3/bluemix-go/trace"

	v2 "github.com/Mavrickk3/bluemix-go/api/container/containerv2"
)

func main() {

	var Zone string
	flag.StringVar(&Zone, "zone", "", "Zone")
	flag.Parse()
	fmt.Println("[FLAG]Zone: ", Zone)
	c := new(bluemix.Config)

	trace.Logger = trace.NewLogger("true")

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	target := v2.ClusterTargetHeader{}

	v2Client, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	dedicatedHostFlavorAPI := v2Client.DedicatedHostFlavor()

	dhf, err := dedicatedHostFlavorAPI.ListDedicatedHostFlavors(Zone, target)
	if err != nil {
		fmt.Printf("ListDedicatedHostFlavors was not successful: %v \n", err)
		return
	}
	fmt.Printf("ListDedicatedHostFlavors was successful: %v \n", dhf)

}
