package registryv1

import (
	"bufio"
	"fmt"
	"testing"
	"os"
	"path"

	ibmcloud "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/api/account/accountv1"
	"github.com/IBM-Cloud/bluemix-go/api/iam/iamv1"
)

func TestCreateCluster(t *testing.T) {
	t.Parallel()
	c := new(ibmcloud.Config)
	c.Region = "us-east"
	c.BluemixAPIKey = ""

	session, _:= session.New(c)
	_, err := accountv1.New(session)
	if err != nil {
		fmt.Println(err)
	}
	iamAPI,err := iamv1.New(session)
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
	registryClient, _:= New(session)

	file, err := os.Open(path.Join("/" , "home", "kuschel", "go","src","github.com","IBM-Cloud","bluemix-go","test","docker.tar.gz"))
	if err != nil {
		fmt.Println(err)
	}

	imagebuildResponse, err :=registryClient.Builds().ImageBuild(*requestStruct, bufio.NewReader(file),*headerStruct)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(imagebuildResponse)
}
