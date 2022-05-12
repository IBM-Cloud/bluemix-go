package main

import (
	"flag"
	"fmt"
	"log"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {
	var VpcID string
	flag.StringVar(&VpcID, "vpcid", "", "VpcID")
	var SubnetID string
	flag.StringVar(&SubnetID, "subnetid", "", "SubnetID")
	var Name string
	flag.StringVar(&Name, "name", "bluemixV2Test", "Name")
	var Zone string
	flag.StringVar(&Zone, "zone", "us-south-1", "Zone")
	var HostPoolID string
	flag.StringVar(&HostPoolID, "hostpoolid", "", "HostPoolID")
	flag.Parse()
	fmt.Println("[FLAG]VpcID: ", VpcID)
	fmt.Println("[FLAG]SubnetID: ", SubnetID)
	fmt.Println("[FLAG]Name: ", Name)
	fmt.Println("[FLAG]Zone: ", Zone)
	fmt.Println("[FLAG]HostPoolID: ", HostPoolID)
	c := new(bluemix.Config)
	trace.Logger = trace.NewLogger("true")
	var clusterInfo = v2.ClusterCreateRequest{
		DisablePublicServiceEndpoint: true,
		Name:                         Name,
		Provider:                     "vpc-gen2",
		WorkerPools: v2.WorkerPoolConfig{
			DiskEncryption: false,
			Flavor:         "bx2d.16x64",
			VpcID:          VpcID,
			WorkerCount:    1,
			HostPoolID:     HostPoolID,
			Zones: []v2.Zone{
				{
					ID:       Zone,
					SubnetID: SubnetID,
				},
			},
		},
	}
	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}
	target := v2.ClusterTargetHeader{}
	v2Client, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	clustersAPI := v2Client.Clusters()
	out, err := clustersAPI.Create(clusterInfo, target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("out=", out)
}
