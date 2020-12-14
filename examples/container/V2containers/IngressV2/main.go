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

	// CREATE INGRESS SECRET
	req := v2.SecretCreateConfig{
		Cluster:     clusterID,
		Name:        "testabc123",
		CRN:         certcrn,
		Persistence: true,
	}
	resp, err := ingressAPI.CreateIngressSecret(req)
	fmt.Println("err=", err)

	// Get INGRESS SECRET
	_, err = ingressAPI.GetIngressSecret(clusterID, "testabc123", resp.Namespace)
	fmt.Println("err=", err)

	// Delete INGRESS SECRET
	req1 := v2.SecretDeleteConfig{
		Cluster:   clusterID,
		Name:      "testabc123",
		Namespace: resp.Namespace,
	}
	err = ingressAPI.DeleteIngressSecret(req1)
	fmt.Println("err=", err)
}
