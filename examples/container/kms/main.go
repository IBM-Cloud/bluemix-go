package main

import (
	"flag"
	"log"
	"os"

	"github.com/IBM-Cloud/bluemix-go/session"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/trace"
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

	kmsClient, err := v1.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	kmsAPI := kmsClient.Kms()
	target := v1.ClusterHeader{}

	//Enable the Kms
	var kmsConfig = v1.KmsEnableReq{
		Cluster:         clusterID,
		Kms:             kmsID,
		Crk:             rootKey,
		PrivateEndpoint: privateEndpoint,
	}
	err = kmsAPI.EnableKms(kmsConfig, target)
	if err != nil {
		log.Fatal(err)
	}

}
