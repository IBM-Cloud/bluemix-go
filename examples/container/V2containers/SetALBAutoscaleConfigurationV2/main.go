package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
)

func main() {

	var clusterID, albID string
	var cpuAverageUtilization, minReplicas, maxReplicas int
	flag.StringVar(&clusterID, "clusterNameOrID", "", "cluster name or ID")
	flag.StringVar(&albID, "albID", "", "ALB ID")
	flag.IntVar(&minReplicas, "minReplicas", 0, "minimum ALB replicas")
	flag.IntVar(&maxReplicas, "maxReplicas", 0, "maximum ALB replicas")
	flag.IntVar(&cpuAverageUtilization, "cpuAverageUtilization", 0, "the CPU Average Utilization")
	flag.Parse()

	trace.Logger = trace.NewLogger("true")

	if clusterID == "" || albID == "" {
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

	autoscaleConf := v2.AutoscaleDetails{
		Config: &v2.AutoscaleConfig{
			MinReplicas:           minReplicas,
			MaxReplicas:           maxReplicas,
			CPUAverageUtilization: cpuAverageUtilization,
		},
	}

	albAPI := clusterClient.Albs()

	err = albAPI.SetALBAutoscaleConfiguration(clusterID, albID, autoscaleConf, target)
	fmt.Println("setErr: ", err)

	getAutoscaleConf, err := albAPI.GetALBAutoscaleConfiguration(clusterID, albID, target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("getAutoscaleConf.Config=%+v\n", getAutoscaleConf.Config)
}
