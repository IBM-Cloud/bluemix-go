package main

import (
	"flag"
	"log"
	"os"

	"github.com/Mavrickk3/bluemix-go"
	"github.com/Mavrickk3/bluemix-go/api/functions"
	"github.com/Mavrickk3/bluemix-go/session"
	"github.com/Mavrickk3/bluemix-go/trace"
)

func main() {

	c := new(bluemix.Config)
	var namespaceName string
	flag.StringVar(&namespaceName, "namespace", "", "Namespace ID")
	var resourceGroupID string
	flag.StringVar(&resourceGroupID, "resourceGroupID", "", "resourceGroupID for namespace")

	flag.Parse()

	if namespaceName == "" || resourceGroupID == "" {
		flag.Usage()
		os.Exit(1)
	}

	trace.Logger = trace.NewLogger("true")
	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	nsCFClient, err := functions.NewCF(sess)
	if err != nil {
		log.Fatal(err)
	}
	nsCFAPI := nsCFClient.Namespaces()

	log.Println(">>>>>>>>>  List CF namespaces")
	cfList, err := nsCFAPI.GetCloudFoundaryNamespaces()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfList)

	nsClient, err := functions.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	nsAPI := nsClient.Namespaces()

	log.Println(">>>>>>>>>  Namespace create")
	resourcePlanID := "functions-base-plan"
	namespaceOpts := functions.CreateNamespaceOptions{
		Name:            &namespaceName,
		ResourceGroupID: &resourceGroupID,
		ResourcePlanID:  &resourcePlanID,
	}
	namespaceResponse, err := nsAPI.CreateNamespace(namespaceOpts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(">>>>>>>>>  List namespaces")
	_, err = nsAPI.GetNamespaces()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(">>>>>>>>>  Get namespace by ID")
	getNamespaceOpts := functions.GetNamespaceOptions{
		ID: namespaceResponse.ID,
	}
	_, err = nsAPI.GetNamespace(getNamespaceOpts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(">>>>>>>>>  Update namespace by ID")
	updateName := "namespace-update-01"
	updateNamespaceOpts := functions.UpdateNamespaceOptions{
		ID:   namespaceResponse.ID,
		Name: &updateName,
	}
	_, err = nsAPI.UpdateNamespace(updateNamespaceOpts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(">>>>>>>>>  Delete Namespace")
	_, err = nsAPI.DeleteNamespace(*namespaceResponse.ID)
	if err != nil {
		log.Fatal(err)
	}

}
