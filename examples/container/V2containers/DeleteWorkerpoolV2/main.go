package main

import (
	"fmt"
	"log"

	bluemix "github.com/Mavrickk3/bluemix-go"
	"github.com/Mavrickk3/bluemix-go/session"

	v2 "github.com/Mavrickk3/bluemix-go/api/container/containerv2"

	"github.com/Mavrickk3/bluemix-go/trace"
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

	target := v2.ClusterTargetHeader{}

	clusterClient, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	workerpoolAPI := clusterClient.WorkerPools()

	err1 := workerpoolAPI.DeleteWorkerPool("bm64u3ed02o93vv36hb0", "bm64u3ed02o93vv36hb0-502aed1", target)

	//out,err=

	fmt.Println("err=", err1)

}
