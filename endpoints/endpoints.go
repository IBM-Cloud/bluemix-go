package endpoints

import (
	"fmt"

	"github.com/IBM-Bluemix/bluemix-go/helpers"
)

//EndpointLocator ...
type EndpointLocator interface {
	AccountManagementEndpoint() (string, error)
	CFAPIEndpoint() (string, error)
	ClusterEndpoint() (string, error)
	IAMEndpoint() (string, error)
	UAAEndpoint() (string, error)
}

var regionToEndpoint = map[string]map[string]string{
	"cf": {
		"us-south": "https://api.ng.bluemix.net",
		"eu-gb":    "https://api.eu-gb.bluemix.net",
		"au-syd":   "https://api.au-syd.bluemix.net",
	},
	"iam": {
		"us-south": "https://iam.ng.bluemix.net",
		"eu-gb":    "https://iam.eu-gb.bluemix.net",
		"au-syd":   "https://iam.au-syd.bluemix.net",
	},

	"uaa": {
		"us-south": "https://login.ng.bluemix.net/UAALoginServerWAR",
		"eu-gb":    "https://login.eu-gb.bluemix.net/UAALoginServerWAR",
		"au-syd":   "https://login.au-syd.bluemix.net/UAALoginServerWAR",
	},
	"account": {
		"us-south": "https://accountmanagement.ng.bluemix.net",
		"eu-gb":    "https://accountmanagement.eu-gb.bluemix.net",
		"au-syd":   "https://accountmanagement.au-syd.bluemix.net",
	},
	"cs": {
		"us-south": "https://us-south.containers.bluemix.net",
	},
}

func init() {
	//TODO populate the endpoints which can be retrieved from given endpoints dynamically
	//Example - UAA can be found from the CF endpoint
}

type endpointLocator struct {
	region string
}

//NewEndpointLocator ...
func NewEndpointLocator(region string) EndpointLocator {
	return &endpointLocator{region: region}
}

func (e *endpointLocator) CFAPIEndpoint() (string, error) {
	if ep, ok := regionToEndpoint["cf"][e.region]; ok {
		//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
		return helpers.EnvFallBack([]string{"IBMCLOUD_CF_API_ENDPOINT"}, ep), nil

	}
	return "", fmt.Errorf("Cloud Foundry endpoint doesn't exist for region: %q", e.region)
}

func (e *endpointLocator) UAAEndpoint() (string, error) {
	if ep, ok := regionToEndpoint["uaa"][e.region]; ok {
		//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
		return helpers.EnvFallBack([]string{"IBMCLOUD_UAA_ENDPOINT"}, ep), nil

	}
	return "", fmt.Errorf("UAA endpoint doesn't exist for region: %q", e.region)
}

func (e *endpointLocator) AccountManagementEndpoint() (string, error) {
	if ep, ok := regionToEndpoint["account"][e.region]; ok {
		//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
		return helpers.EnvFallBack([]string{"IBMCLOUD_ACCOUNT_MANAGEMENT_API_ENDPOINT"}, ep), nil

	}
	return "", fmt.Errorf("Account Management endpoint doesn't exist for region: %q", e.region)
}

func (e *endpointLocator) IAMEndpoint() (string, error) {
	if ep, ok := regionToEndpoint["iam"][e.region]; ok {
		//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
		return helpers.EnvFallBack([]string{"IBMCLOUD_IAM_API_ENDPOINT"}, ep), nil

	}
	return "", fmt.Errorf("IAM  endpoint doesn't exist for region: %q", e.region)
}

func (e *endpointLocator) ClusterEndpoint() (string, error) {
	if ep, ok := regionToEndpoint["cs"][e.region]; ok {
		//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
		return helpers.EnvFallBack([]string{"IBMCLOUD_CS_API_ENDPOINT"}, ep), nil
	}
	return "", fmt.Errorf("Container Service endpoint doesn't exist for region: %q", e.region)
}
