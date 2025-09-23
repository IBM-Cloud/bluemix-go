package main

import (
	"flag"
	"fmt"
	"log"

	bluemix "github.com/Mavrickk3/bluemix-go"
	"github.com/Mavrickk3/bluemix-go/session"

	v1 "github.com/Mavrickk3/bluemix-go/api/container/containerv1"
	"github.com/Mavrickk3/bluemix-go/trace"
)

func main() {

	c := new(bluemix.Config)

	var ResourceGroup string
	flag.StringVar(&ResourceGroup, "resourcegroup", "", "ResourceGroup")

	var Region string
	flag.StringVar(&Region, "region", "", "Region")

	var Cluster string
	flag.StringVar(&Cluster, "cluster", "", "Cluster")

	var WorkerPool string
	flag.StringVar(&WorkerPool, "workerpool", "", "WorkerPool")

	flag.Parse()

	trace.Logger = trace.NewLogger("true")

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	target := v1.ClusterTargetHeader{}
	target.Region = Region
	target.ResourceGroup = ResourceGroup

	clusterClient, err := v1.New(sess)
	if err != nil {
		log.Fatal(err)
	}

	workerPoolAPI := clusterClient.WorkerPools()

	pool, err := workerPoolAPI.GetWorkerPool(Cluster, WorkerPool, target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WorkerPool get was successful\n %v", pool)
}
