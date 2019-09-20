package containerv2

import (
	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/client"
)

//Clusters interface
type Clusters interface {
	List(target v1.ClusterTargetHeader) ([]v1.ClusterInfo, error)
	//TODO Add other opertaions
}
type clusters struct {
	client     *client.Client
	pathPrefix string
}

func newClusterAPI(c *client.Client) Clusters {
	return &clusters{
		client:     c,
		pathPrefix: "/v2/vpc/",
	}
}

//List ...
func (r *clusters) List(target v1.ClusterTargetHeader) ([]v1.ClusterInfo, error) {
	clusters := []v1.ClusterInfo{}
	_, err := r.client.Get(r.pathPrefix+"getClusters", &clusters, target.ToMap())
	if err != nil {
		return nil, err
	}

	return clusters, err
}
