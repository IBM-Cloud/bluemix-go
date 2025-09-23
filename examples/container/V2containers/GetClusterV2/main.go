package main

import (
	"fmt"
	"log"

	bluemix "github.com/Mavrickk3/bluemix-go"

	v2 "github.com/Mavrickk3/bluemix-go/api/container/containerv2"
	"github.com/Mavrickk3/bluemix-go/session"

	"github.com/Mavrickk3/bluemix-go/trace"
)

func main() {

	c := new(bluemix.Config)

	trace.Logger = trace.NewLogger("false")

	target := v2.ClusterTargetHeader{}

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	clusterClient, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	clustersAPI := clusterClient.Clusters()

	out, err := clustersAPI.GetCluster("bm64u3ed02o93vv36hb0", target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", out)
}
