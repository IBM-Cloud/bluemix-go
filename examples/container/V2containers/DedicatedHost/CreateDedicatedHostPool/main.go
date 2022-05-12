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
	var FlavorClass string
	flag.StringVar(&FlavorClass, "flavorclass", "", "FlavorClass")
	var Metro string
	flag.StringVar(&Metro, "metro", "dal", "Metro")
	var Name string
	flag.StringVar(&Name, "name", "bluemixV2Test", "Name")
	flag.Parse()
	fmt.Println("[FLAG]FlavorClass: ", FlavorClass)
	fmt.Println("[FLAG]Metro: ", Metro)
	fmt.Println("[FLAG]Name: ", Name)
	c := new(bluemix.Config)

	trace.Logger = trace.NewLogger("true")

	var createDedicatedHostPool = v2.CreateDedicatedHostPoolRequest{
		FlavorClass: FlavorClass,
		Metro:       Metro,
		Name:        Name,
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

	dh, err := dedicatedHostPoolAPI.CreateDedicatedHostPool(createDedicatedHostPool, target)
	if err != nil {
		fmt.Printf("Create was not successful: %v \n", err)
		return
	}
	fmt.Printf("Create dedicated hostpool response: %v \n", dh)
}
