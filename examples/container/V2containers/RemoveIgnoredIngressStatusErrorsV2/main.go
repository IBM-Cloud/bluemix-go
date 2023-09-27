package main

import (
	"fmt"
	"log"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
)

func main() {

	c := new(bluemix.Config)

	trace.Logger = trace.NewLogger("true")

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	clusterID := "ck6k27hd0s1542093c6g"

	target := v2.ClusterTargetHeader{}

	clusterClient, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}

	albAPI := clusterClient.Albs()
	errorCodes := v2.IgnoredIngressStatusErrors{
		Cluster: clusterID,
		IgnoredErrors: []string{
			"ERRADRUH",
		},
	}

	err = albAPI.RemoveIgnoredIngressStatusErrors(errorCodes, target)
	fmt.Println("err: ", err)
}
