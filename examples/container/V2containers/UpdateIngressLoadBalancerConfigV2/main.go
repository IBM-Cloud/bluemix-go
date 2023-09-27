package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
)

func main() {

	var clusterID, lbType string
	flag.StringVar(&clusterID, "clusterNameOrID", "", "cluster name or ID")
	flag.StringVar(&lbType, "lbType", "", "loadbalancer type")
	flag.Parse()

	trace.Logger = trace.NewLogger("true")

	if clusterID == "" || lbType == "" {
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

	target := v2.ClusterTargetHeader{}

	clusterClient, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}

	lbConf := v2.ALBLBConfig{
		Cluster: clusterID,
		ProxyProtocol: &v2.ALBLBProxyProtocolConfig{
			Enable: true,
		},
	}

	albAPI := clusterClient.Albs()

	err = albAPI.UpdateIngressLoadBalancerConfig(lbConf, target)
	fmt.Println("updateErr: ", err)

	getLbConf, err := albAPI.GetIngressLoadBalancerConfig(clusterID, lbType, target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("getLbConf.ProxyProtocol: %+v\n", getLbConf.ProxyProtocol)
}
