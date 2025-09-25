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

	var certcrn, clusterID, instancecrn, fieldcrn string
	flag.StringVar(&certcrn, "certcrn", "", "CRN for certificate")
	flag.StringVar(&clusterID, "clusterNameOrID", "", "cluster name or ID")
	flag.StringVar(&instancecrn, "instance-crn", "", "crn for secrets manager instance")
	flag.StringVar(&fieldcrn, "field-crn", "", "crn for opaque secret field")
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
	ingressClient, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	ingressAPI := ingressClient.Ingresses()

	if instancecrn != "" {
		if err := instance(ingressAPI, clusterID, instancecrn); err != nil {
			fmt.Println("err=", err)
			os.Exit(1)
		}
	}

	if certcrn != "" {
		if err := secret(ingressAPI, clusterID, certcrn); err != nil {
			fmt.Println("err=", err)
			os.Exit(1)
		}
	}

	if fieldcrn != "" {
		if err := opaqueSecret(ingressAPI, clusterID, fieldcrn); err != nil {
			fmt.Println("err=", err)
			os.Exit(1)
		}
	}
}

func secret(ingressAPI v2.Ingress, clusterID, certCRN string) error {
	trace.Logger = trace.NewLogger("true")
	if certCRN == "" || clusterID == "" {
		flag.Usage()
		os.Exit(1)
	}

	// CREATE INGRESS SECRET
	req := v2.SecretCreateConfig{
		Cluster:     clusterID,
		Name:        "testabc123",
		CRN:         certCRN,
		Persistence: true,
	}
	resp, err := ingressAPI.CreateIngressSecret(req)
	if err != nil {
		return err
	}

	// Get INGRESS SECRET
	_, err = ingressAPI.GetIngressSecret(clusterID, "testabc123", resp.Namespace)
	if err != nil {
		return err
	}

	// Delete INGRESS SECRET
	req1 := v2.SecretDeleteConfig{
		Cluster:   clusterID,
		Name:      "testabc123",
		Namespace: resp.Namespace,
	}
	return ingressAPI.DeleteIngressSecret(req1)
}

func opaqueSecret(ingressAPI v2.Ingress, clusterID, fieldCRN string) error {
	trace.Logger = trace.NewLogger("true")
	if fieldCRN == "" || clusterID == "" {
		flag.Usage()
		os.Exit(1)
	}

	// CREATE INGRESS SECRET
	req := v2.SecretCreateConfig{
		Cluster:     clusterID,
		Name:        "testabc123",
		Persistence: true,
		Type:        "Opaque",
		FieldsToAdd: []v2.FieldAdd{{
			CRN: fieldCRN,
		}},
	}
	resp, err := ingressAPI.CreateIngressSecret(req)
	if err != nil {
		return err
	}

	// Get INGRESS SECRET
	secret, err := ingressAPI.GetIngressSecret(clusterID, "testabc123", resp.Namespace)
	if err != nil {
		return err
	}

	// Add INGRESS SECRET field
	req1 := v2.SecretUpdateConfig{
		Cluster:   clusterID,
		Name:      "testabc123",
		Namespace: resp.Namespace,
		FieldsToRemove: []v2.FieldRemove{{
			Name: secret.Fields[0].Name,
		}},
	}

	_, err = ingressAPI.RemoveIngressSecretField(req1)
	if err != nil {
		return err
	}

	// Delete INGRESS SECRET
	req2 := v2.SecretDeleteConfig{
		Cluster:   clusterID,
		Name:      "testabc123",
		Namespace: resp.Namespace,
	}
	return ingressAPI.DeleteIngressSecret(req2)
}

func instance(ingressAPI v2.Ingress, clusterID, instanceCRN string) error {
	trace.Logger = trace.NewLogger("true")
	if instanceCRN == "" || clusterID == "" {
		flag.Usage()
		os.Exit(1)
	}

	// CREATE INGRESS SECRET
	req := v2.InstanceRegisterConfig{
		Cluster: clusterID,
		CRN:     instanceCRN,
	}
	resp, err := ingressAPI.RegisterIngressInstance(req)
	if err != nil {
		return err
	}

	// Get INGRESS SECRET
	_, err = ingressAPI.GetIngressInstance(clusterID, resp.Name)
	if err != nil {
		return err
	}

	// Delete INGRESS SECRET
	req1 := v2.InstanceDeleteConfig{
		Cluster: clusterID,
		Name:    resp.Name,
	}
	return ingressAPI.DeleteIngressInstance(req1)
}
