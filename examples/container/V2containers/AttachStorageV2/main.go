package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	bluemix "github.com/Mavrickk3/bluemix-go"
	"github.com/Mavrickk3/bluemix-go/session"
	"github.com/Mavrickk3/bluemix-go/trace"

	v2 "github.com/Mavrickk3/bluemix-go/api/container/containerv2"
)

func main() {

	var clusterID string
	flag.StringVar(&clusterID, "clusterID", "", "Cluster ID or Name")

	var workerID string
	flag.StringVar(&workerID, "workerID", "", "worker ID of the worker node in the cluster")

	var volumeID string
	flag.StringVar(&volumeID, "volumeID", "", "volumeID of the volume to be attched to the worker")

	flag.Parse()

	if clusterID == "" || workerID == "" || volumeID == "" {
		flag.Usage()
		os.Exit(1)
	}
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
	workersAPI := clusterClient.Workers()

	attachVolumeRequest := v2.VolumeRequest{
		Cluster:  clusterID,
		VolumeID: volumeID,
		Worker:   workerID,
	}
	volumeattached, errC := workersAPI.CreateStorageAttachment(attachVolumeRequest, target)
	if errC != nil {
		fmt.Println(errC)
		return
	}
	fmt.Println("Volume attached: ", volumeattached)
	volumeAttachmentID := volumeattached.Id
	time.Sleep(10 * time.Second)

	volumeAttachment, errG := workersAPI.GetStorageAttachment(clusterID, workerID, volumeAttachmentID, target)
	if errG != nil {
		fmt.Println(errG)
		return
	}
	fmt.Println("Volume attachment with worker nodes: ", volumeAttachment)
	time.Sleep(5 * time.Second)

	detachVolumeRequest := v2.VolumeRequest{
		Cluster:            clusterID,
		VolumeAttachmentID: volumeAttachmentID,
		Worker:             workerID,
	}
	out, errD := workersAPI.DeleteStorageAttachment(detachVolumeRequest, target)
	if errD != nil {
		fmt.Println(errD)
		return
	}
	fmt.Println("Volume attachment removed", out)
}
