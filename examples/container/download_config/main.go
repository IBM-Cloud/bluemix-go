package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	v2 "github.com/Mavrickk3/bluemix-go/api/container/containerv2"
	"github.com/Mavrickk3/bluemix-go/session"
	"github.com/Mavrickk3/bluemix-go/trace"
)

func main() {
	var clusterName string
	flag.StringVar(&clusterName, "clustername", "", "The cluster whose config will be downloaded")

	var path string
	flag.StringVar(&path, "path", "", "The Path where the config will be downloaded")

	var resourceGroup string
	flag.StringVar(&resourceGroup, "resourcegroup", "", "ResourceGroup where the cluster is deployed")

	var endpointType string
	flag.StringVar(&endpointType, "endpoint", "", "Endpoint defines how the kubeconfig will connect to the cluster. Can be public, private and vpe in case of VPC")

	var admin bool
	flag.BoolVar(&admin, "admin", false, "If true download the admin config")

	var network bool
	flag.BoolVar(&network, "network", false, "If true download the calico network config")

	flag.Parse()
	trace.Logger = trace.NewLogger("true")
	if clusterName == "" || path == "" {
		flag.Usage()
		os.Exit(1)
	}
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	target := v2.ClusterTargetHeader{
		ResourceGroup: resourceGroup,
	}
	clusterClient, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	clustersAPI := clusterClient.Clusters()

	if network {
		kubeConfig, configPath, err := clustersAPI.StoreConfigDetail(clusterName, path, admin, network, target, endpointType)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(kubeConfig, configPath.FilePath)
	} else {
		configPath, err := clustersAPI.GetClusterConfigDetail(clusterName, path, admin, target, endpointType)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(configPath.FilePath)
	}
}
