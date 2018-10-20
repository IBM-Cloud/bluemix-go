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
func TestBuild(t *testing.T) {
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
func TestNamespaces(t *testing.T) {
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

	headerStruct := &NamespaceTargetHeader{
		AccountID: userInfo.Account.Bss,
	}

	registryClient, _ := New(session)
	retval, err := registryClient.Namespaces().GetNamespaces(*headerStruct)
	fmt.Printf("%v", retval)
	namespaces := "devtest"
	nameback, err := registryClient.Namespaces().AddNamespace(namespaces, *headerStruct)
	fmt.Printf("%v", nameback)
	err = registryClient.Namespaces().DeleteNamespace(namespaces, *headerStruct)
	fmt.Printf("%v", err)

}

func TestTokens(t *testing.T) {
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

	headerStruct := &TokenTargetHeader{
		AccountID: userInfo.Account.Bss,
	}

	registryClient, _ := New(session)
	params := DefaultIssueTokenRequest()
	params.Description = "TTTEEEEESSSSSSSSSSSSSSSSSSSSST"
	retval, err := registryClient.Tokens().IssueToken(*params, *headerStruct)
	fmt.Printf("%v", retval)

	retval2, err := registryClient.Tokens().GetTokens(*headerStruct)
	fmt.Printf("%v", retval2)
	fmt.Printf("--------------------------------------------")
	retval1, err := registryClient.Tokens().GetToken(retval.ID, *headerStruct)
	fmt.Printf("%v", retval1)
	err = registryClient.Tokens().DeleteToken(retval.ID, *headerStruct)

}

func TestImages(t *testing.T) {
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

	headerStruct := &ImageTargetHeader{
		AccountID: userInfo.Account.Bss,
	}

	registryClient, _ := New(session)
	params := DefaultGetImageRequest()
	retval, err := registryClient.Images().GetImages(*params, *headerStruct)
	fmt.Printf("%v", retval)
	retval1, err := registryClient.Images().InspectImage("registry.ng.bluemix.net/gpfs/sklm:latest", *headerStruct)
	fmt.Printf("%v", retval1)
	retval3, err := registryClient.Images().ImageVulnerabilities("registry.ng.bluemix.net/gpfs/sklm:latest", *DefaultImageVulnerabilitiesRequest(), *headerStruct)
	//retval3, err := registryClient.Images().DeleteImage("registry.ng.bluemix.net/bkuschel/testimage:latest", *headerStruct)
	fmt.Printf("%v", retval3)
}
