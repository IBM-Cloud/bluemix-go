package main

import (
	"fmt"
	"log"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/session"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"

	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	c := new(bluemix.Config)

	trace.Logger = trace.NewLogger("true")

	var clusterInfo = v2.ClusterCreateRequest{
		DisablePublicServiceEndpoint: true,
		KubeVersion:                  "4.3.23_openshift",
		Name:                         "mycluscretaed123",
		PodSubnet:                    "172.30.0.0/16",
		Provider:                     "vpc-gen2",
		ServiceSubnet:                "172.21.0.0/16",
		CosInstanceCRN:               "crn:v1:bluemix:public:cloud-object-storage:global:a/96fe4b4beb8947bf85223e69dab47878:cf577e01-4095-4b5e-a223-1d515a825cfd::",
		WorkerPools: v2.WorkerPoolConfig{
			DiskEncryption: true,
			Flavor:         "bx2.16x64",
			Isolation:      "shared",
			Name:           "mywork1",
			VpcID:          "r018-b50f22c1-f9c1-4337-8cff-5eb89d53f604",
			WorkerCount:    2,
			Zones: []v2.Zone{
				{
					ID:       "eu-gb-1",
					SubnetID: "0787-5b33830d-616b-4f5c-861a-40849e5203ce",
				},
			},
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
