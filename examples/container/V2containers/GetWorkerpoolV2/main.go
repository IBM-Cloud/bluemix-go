package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
)

func main() {

	var Cluster string
	flag.StringVar(&Cluster, "cluster", "", "Cluster")

	var WorkerPool string
	flag.StringVar(&WorkerPool, "workerpool", "", "WorkerPool")

	flag.Parse()

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

	out, err := workerpoolAPI.GetWorkerPool(Cluster, WorkerPool, target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Get workerpool request was successful")
	json, _ := json.Marshal(out)
	fmt.Println("Response:", string(json))
}
