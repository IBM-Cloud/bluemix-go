package registryv1

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"testing"

	ibmcloud "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/api/account/accountv1"
	"github.com/IBM-Cloud/bluemix-go/api/iam/iamv1"
	"github.com/IBM-Cloud/bluemix-go/session"
)

const (
	dockerfile = `FROM golang:1.7-alpine3.6

ARG JQ_VERSION=jq-1.5
RUN apk update
RUN apk add build-base git bash
ADD https://github.com/stedolan/jq/releases/download/${JQ_VERSION}/jq-linux64 /tmp`

	dockerfileName = "Dockerfile"
)

func createTestTar() *bytes.Buffer {
	// Create and add some files to the archive.
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	hdr := &tar.Header{
		Name: dockerfileName,
		Mode: 0600,
		Size: int64(len(dockerfile)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		log.Fatal(err)
	}
	if _, err := tw.Write([]byte(dockerfile)); err != nil {
		log.Fatal(err)
	}
	return &buf
}
func TestCreateCluster(t *testing.T) {
	t.Parallel()
	c := new(ibmcloud.Config)
	c.Region = "us-east"
	c.BluemixAPIKey = ""

	session, _ := session.New(c)
	_, err := accountv1.New(session)
	if err != nil {
		fmt.Println(err)
	}
	iamAPI, err := iamv1.New(session)
	identityAPI := iamAPI.Identity()
	userInfo, err := identityAPI.UserInfo()
	if err != nil {
		fmt.Println(err)
	}
	requestStruct := DefaultImageBuildRequest()
	requestStruct.T = "registry.ng.bluemix.net/bkuschel/testimage"

	headerStruct := &BuildTargetHeader{
		AccountID: userInfo.Account.Bss,
	}
	registryClient, _ := New(session)

	buffer := createTestTar()
	err = registryClient.Builds().ImageBuild(*requestStruct, buffer, *headerStruct, func(respV ImageBuildResponse) bool {
		fmt.Println(respV)
		return true
	})
}
