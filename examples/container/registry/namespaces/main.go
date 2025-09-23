package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	registryv1 "github.com/Mavrickk3/bluemix-go/api/container/registryv1"
	"github.com/Mavrickk3/bluemix-go/session"
	"github.com/Mavrickk3/bluemix-go/trace"
)

func main() {

	var namespace string
	flag.StringVar(&namespace, "namespace", "", "Namespace")
	var accountID string
	flag.StringVar(&accountID, "accountID", "", "Account ID")
	var resourceGroup string
	flag.StringVar(&resourceGroup, "resourceGroup", "", "Resource Group ID")

	flag.Parse()

	trace.Logger = trace.NewLogger("true")
	if namespace == "" {
		flag.Usage()
		os.Exit(1)
	}

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	crClient, err := registryv1.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	crAPI := crClient.Namespaces()

	// CREATE cr SECRET
	req := registryv1.NamespaceTargetHeader{
		AccountID: accountID,
	}
	resp, err := crAPI.AddNamespace(namespace, req)
	fmt.Println("err=", resp)

	req.ResourceGroup = resourceGroup
	resp3, err := crAPI.AssignNamespace(namespace, req)
	fmt.Println("err=", resp3)

	resp5, err := crAPI.GetNamespaces(req)
	fmt.Println("err=", resp5)

	resp4, err := crAPI.GetDetailedNamespaces(req)
	fmt.Println("err=", resp4)

	err = crAPI.DeleteNamespace(namespace, req)
	fmt.Println("err=", err)
}
