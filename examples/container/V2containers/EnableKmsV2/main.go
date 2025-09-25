package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
)

func main() {

	var clusterID string
	flag.StringVar(&clusterID, "clusterID", "", "cluster name or ID")

	var kmsID string
	flag.StringVar(&kmsID, "kmsID", "", "kms Id")

	var rootKey string
	flag.StringVar(&rootKey, "rootKey", "", "root Key")

	var privateEndpoint bool
	flag.BoolVar(&privateEndpoint, "privateEndpoint", false, " private Endpoint(true/false)")

	flag.Parse()

	trace.Logger = trace.NewLogger("true")
	if kmsID == "" || clusterID == "" || rootKey == "" {
		flag.Usage()
		os.Exit(1)
	}

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	kmsClient, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	kmsAPI := kmsClient.Kms()
	target := v2.ClusterHeader{}

	//Enable the Kms
	var kmsConfig = v2.KmsEnableReq{
		Cluster:         clusterID,
		Kms:             kmsID,
		Crk:             rootKey,
		PrivateEndpoint: privateEndpoint,
	}

	err = kmsAPI.EnableKms(kmsConfig, target)
	fmt.Println("err=", err)

}
