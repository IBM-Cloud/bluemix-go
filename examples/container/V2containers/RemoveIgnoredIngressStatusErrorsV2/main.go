package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
)

func main() {

	var errCode, clusterID string
	flag.StringVar(&errCode, "errCode", "", "ignored error code")
	flag.StringVar(&clusterID, "clusterNameOrID", "", "cluster name or ID")
	flag.Parse()

	trace.Logger = trace.NewLogger("true")

	if clusterID == "" || errCode == "" {
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

	albAPI := clusterClient.Albs()
	errorCodes := v2.IgnoredIngressStatusErrors{
		Cluster: clusterID,
		IgnoredErrors: []string{
			errCode,
		},
	}

	err = albAPI.RemoveIgnoredIngressStatusErrors(errorCodes, target)
	fmt.Println("err: ", err)
}
