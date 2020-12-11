package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	var certcrn string
	flag.StringVar(&certcrn, "certcrn", "", "Alb Id")

	var clusterID string
	flag.StringVar(&clusterID, "clusterNameOrID", "", "cluster name or ID")

	flag.Parse()

	trace.Logger = trace.NewLogger("true")
	if certcrn == "" || clusterID == "" {
		flag.Usage()
		os.Exit(1)
	}

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	ingressClient, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	ingressAPI := ingressClient.Ingresses()
	target := v2.ClusterTargetHeader{}

	// CREATE INGRESS SECRET
	req := v2.SecretCreateConfig{
		Cluster:     "bugi52rf0rtfgadjfso0",
		Name:        "testabc2",
		CRN:         "crn:v1:bluemix:public:cloudcerts:us-south:a/883079c85357a1f3f85d968780e56518:b65b5b7f-e904-4d2b-bd87-f0ccd57e76ba:certificate:333d8673f4d03c148ff81192b9edaafc",
		Persistence: true,
	}
	resp, err := ingressAPI.CreateIngressSecret(req, target)
	fmt.Println("err=", err)

	// Get INGRESS SECRET
	_, err = ingressAPI.GetIngressSecret("bugi52rf0rtfgadjfso0", "testabc2", resp.Namespace, target)
	fmt.Println("err=", err)

	// Delete INGRESS SECRET
	req1 := v2.SecretDeleteConfig{
		Cluster:   "bugi52rf0rtfgadjfso0",
		Name:      "testabc2",
		Namespace: resp.Namespace,
	}
	err = ingressAPI.DeleteIngressSecret(req1, target)
	fmt.Println("err=", err)
}
