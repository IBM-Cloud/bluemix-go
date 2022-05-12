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
	var HostID string
	flag.StringVar(&HostID, "hostid", "", "HostID")
	var HostPoolID string
	flag.StringVar(&HostPoolID, "hostpoolid", "", "HostPoolID")
	flag.Parse()
	fmt.Println("[FLAG]HostID: ", HostID)
	fmt.Println("[FLAG]HostPoolID: ", HostPoolID)
	c := new(bluemix.Config)

	trace.Logger = trace.NewLogger("true")

	var removeDedicatedHost = v2.RemoveDedicatedHostRequest{
		HostID:     HostID,
		HostPoolID: HostPoolID,
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

	err = dedicatedHostAPI.RemoveDedicatedHost(removeDedicatedHost, target)
	if err != nil {
		fmt.Printf("Remove was not successful: %v \n", err)
		return
	}
	fmt.Printf("Remove dedicated host was successful \n")
}
