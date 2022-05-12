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
	var HostPoolID string
	flag.StringVar(&HostPoolID, "hostpoolid", "", "HostPoolID")
	flag.Parse()
	fmt.Println("[FLAG]HostPoolID: ", HostPoolID)
	c := new(bluemix.Config)

	trace.Logger = trace.NewLogger("true")

	var removeDedicatedHostPool = v2.RemoveDedicatedHostPoolRequest{
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
	dedicatedHostPoolAPI := v2Client.DedicatedHostPool()

	err = dedicatedHostPoolAPI.RemoveDedicatedHostPool(removeDedicatedHostPool, target)
	if err != nil {
		fmt.Printf("Remove was not successful: %v \n", err)
		return
	}
	fmt.Printf("Remove dedicated hostpool was successful \n")
}
