package main

import (
	"fmt"
	"log"

	bluemix "github.com/Mavrickk3/bluemix-go"
	"github.com/Mavrickk3/bluemix-go/session"
	"github.com/Mavrickk3/bluemix-go/trace"

	v2 "github.com/Mavrickk3/bluemix-go/api/container/containerv2"
)

func main() {

	c := new(bluemix.Config)

	trace.Logger = trace.NewLogger("true")

	var albinfo = v2.AlbCreateReq{
		Cluster:         "bm64u3ed02o93vv36hb0",
		EnableByDefault: true,
		Type:            "private",
		ZoneAlb:         "us-south-1",
	}
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

	alb, err2 := albAPI.CreateAlb(albinfo, target)

	fmt.Println("err=", err2)
	fmt.Println("created albID=", alb.Alb)

}
