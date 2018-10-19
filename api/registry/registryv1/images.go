package registryv1

import (
	"strconv"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

type ImageTargetHeader struct {
	AccountID string
}

//ToMap ...
func (c ImageTargetHeader) ToMap() map[string]string {
	m := make(map[string]string, 1)
	m[accountIDHeader] = c.AccountID
	return m
}

//Subnets interface
type Images interface {
	GetImages(params GetImageRequest, target ImageTargetHeader) (GetImagesResponse, error)
}

type images struct {
	client *client.Client
}

func newImageAPI(c *client.Client) Images {
	return &images{
		client: c,
	}
}

type Digesttags struct {
	Tags map[string][]string
}

type Labels struct {
	Labels map[string][]string
}

type GetImagesResponse []struct {
	ID                      string              `json:"Id"`
	ParentID                string              `json:"ParentId"`
	DigestTags              map[string][]string `json:"DigestTags"`
	RepoDigests             []string            `json:"RepoDigests"`
	Created                 int                 `json:"Created"`
	Size                    int64               `json:"Size"`
	VirtualSize             int64               `json:"VirtualSize"`
	Labels                  map[string]string   `json:"Labels"`
	Vulnerable              string              `json:"Vulnerable"`
	VulnerabilityCount      int                 `json:"VulnerabilityCount"`
	ConfigurationIssueCount int                 `json:"ConfigurationIssueCount"`
	IssueCount              int                 `json:"IssueCount"`
	ExemptIssueCount        int                 `json:"ExemptIssueCount"`
}

/*GetImageRequest contains all the parameters to send to the API endpoint
for the image list operation typically these are written to a http.Request
*/
type GetImageRequest struct {
	/*IncludeIBM
	  Includes IBM-provided public images in the list of images. If this option is not specified, private images are listed only. If this option is specified more than once, the last parsed setting is the setting that is used.
	*/
	IncludeIBM bool
	/*IncludePrivate
	  Includes private images in the list of images. If this option is not specified, private images are listed. If this option is specified more than once, the last parsed setting is the setting that is used.
	*/
	IncludePrivate bool
	/*Namespace
	  Lists images that are stored in the specified namespace only. Query multiple namespaces by specifying this option for each namespace. If this option is not specified, images from all namespaces in the specified IBM Cloud account are listed.
	*/
	Namespace string
	/*Repository
	  Lists images that are stored in the specified repository, under your namespaces. Query multiple repositories by specifying this option for each repository. If this option is not specified, images from all repos are listed.
	*/
	Repository string
	/*Vulnerabilities
	  Displays Vulnerability Advisor status for the listed images. If this option is specified more than once, the last parsed setting is the setting that is used.
	*/
	Vulnerabilities bool
}

func DefaultGetImageRequest() *GetImageRequest {
	return &GetImageRequest{
		IncludeIBM:      false,
		IncludePrivate:  true,
		Namespace:       "",
		Repository:      "",
		Vulnerabilities: true,
	}
}

func (r *images) GetImages(params GetImageRequest, target ImageTargetHeader) (GetImagesResponse, error) {

	var retVal GetImagesResponse
	req := rest.GetRequest(helpers.GetFullURL(*r.client.Config.Endpoint, "/api/v1/images")).
		Query("includeIBM", strconv.FormatBool(params.IncludeIBM)).
		Query("includePrivate", strconv.FormatBool(params.IncludePrivate)).
		Query("vulnerabilities", strconv.FormatBool(params.Vulnerabilities))
	if params.Namespace != "" {
		req = req.Query("namespace", params.Namespace)
	}
	if params.Repository != "repository" {
		req = req.Query("repository", params.Repository)
	}
	for key, value := range target.ToMap() {
		req.Set(key, value)
	}

	_, err := r.client.SendRequest(req, &retVal)
	return retVal, err
}
