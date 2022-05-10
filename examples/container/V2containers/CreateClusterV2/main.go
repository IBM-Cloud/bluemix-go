package main

import (
	"flag"
	"fmt"
	"log"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/session"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"

	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	var KmsInstanceID string
	flag.StringVar(&KmsInstanceID, "kmsid", "", "KmsInstanceID")

	var WorkerVolumeCRKID string
	flag.StringVar(&WorkerVolumeCRKID, "crkid", "", "WorkerVolumeCRKID")

	var VpcID string
	flag.StringVar(&VpcID, "vpcid", "", "VpcID")

	var SubnetID string
	flag.StringVar(&SubnetID, "subnetid", "", "SubnetID")

	var Name string
	flag.StringVar(&Name, "Name", "bluemixV2Test", "Name")

	var Zone string
	flag.StringVar(&Zone, "Zone", "us-south-1", "Zone")

	flag.Parse()
	fmt.Println("[FLAG]KmsInstanceID: ", KmsInstanceID)
	fmt.Println("[FLAG]WorkerVolumeCRKID: ", WorkerVolumeCRKID)
	fmt.Println("[FLAG]VpcID: ", VpcID)
	fmt.Println("[FLAG]SubnetID: ", SubnetID)
	fmt.Println("[FLAG]Name: ", Name)
	fmt.Println("[FLAG]Zone: ", Zone)

	c := new(bluemix.Config)

	trace.Logger = trace.NewLogger("true")

	var wve *v2.WorkerVolumeEncryption
	if KmsInstanceID != "" && WorkerVolumeCRKID != "" {
		wve = &v2.WorkerVolumeEncryption{
			KmsInstanceID:     KmsInstanceID,
			WorkerVolumeCRKID: WorkerVolumeCRKID,
		}
	}

	var clusterInfo = v2.ClusterCreateRequest{
		DisablePublicServiceEndpoint: true,
		Name:                         Name,
		Provider:                     "vpc-gen2",
		WorkerPools: v2.WorkerPoolConfig{
			DiskEncryption: true,
			Flavor:         "bx2.16x64",
			VpcID:          VpcID,
			WorkerCount:    2,
			Zones: []v2.Zone{
				{
					ID:       Zone,
					SubnetID: SubnetID,
				},
			},
			WorkerVolumeEncryption: wve,
		},
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
	clustersAPI := clusterClient.Clusters()

	out, err := clustersAPI.Create(clusterInfo, target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("out=", out)
}
