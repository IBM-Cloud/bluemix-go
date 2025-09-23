package main

import (
	"flag"
	"fmt"
	"log"

	bluemix "github.com/Mavrickk3/bluemix-go"
	"github.com/Mavrickk3/bluemix-go/session"

	v2 "github.com/Mavrickk3/bluemix-go/api/container/containerv2"
)

func main() {

	var clusterIDOrName string
	flag.StringVar(&clusterIDOrName, "cluster-id-or-name", "", "ID or Name of the targeted cluster")

	var requestType string
	flag.StringVar(&requestType, "request-type", "", "'enable' or 'disable' image security enforcement")

	flag.Parse()

	if clusterIDOrName == "" {
		fmt.Println("Please provide a cluster ID or Name to target with -cluster-id-or-name!")
		return
	}

	c := new(bluemix.Config)

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	v2ClusterClient, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}

	v2ClusterAPI := v2ClusterClient.Clusters()

	switch requestType {
	case "enable":
		err = v2ClusterAPI.EnableImageSecurityEnforcement(
			clusterIDOrName,
			v2.ClusterTargetHeader{},
		)
		if err != nil {
			log.Fatal(err)
		}
	case "disable":
		if requestType == "disable" {
			err = v2ClusterAPI.DisableImageSecurityEnforcement(
				clusterIDOrName,
				v2.ClusterTargetHeader{},
			)
			if err != nil {
				log.Fatal(err)
			}
		}
	default:
		fmt.Println("Please provide the request type using request-type!")
	}
}
