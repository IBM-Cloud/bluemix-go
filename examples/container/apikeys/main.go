package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	v1 "github.com/Mavrickk3/bluemix-go/api/container/containerv1"
	"github.com/Mavrickk3/bluemix-go/session"
	"github.com/Mavrickk3/bluemix-go/trace"
)

func main() {

	var clusterID string
	flag.StringVar(&clusterID, "clusterNameOrID", "", "cluster name or ID")

	var region string
	flag.StringVar(&region, "region", "us-south", "region of cluster")

	flag.Parse()

	trace.Logger = trace.NewLogger("true")
	if clusterID == "" {
		flag.Usage()
		os.Exit(1)
	}

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	apiKeyClient, err := v1.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	apiKeyAPI := apiKeyClient.Apikeys()
	target := v1.ClusterTargetHeader{
		Region: region,
	}

	err = apiKeyAPI.ResetApiKey(target)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Reset APIKey to Cluster ", clusterID)

	resp, err := apiKeyAPI.GetApiKeyInfo(clusterID, target)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("APIKey Details %+v", resp)

}
