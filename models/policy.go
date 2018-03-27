package models

import "github.com/IBM-Bluemix/bluemix-go/crn"

type Policy struct {
	ID        string           `json:"id,omitempty"`
	Roles     []PolicyRole     `json:"roles"`
	Resources []PolicyResource `json:"resources"`
	Version   string           `json:"-"`
}

type PolicyRole struct {
	ID          crn.CRN `json:"id"`
	DisplayName string  `json:"displayName"`
	Description string  `json:"description"`
}

type PolicyResource struct {
	ServiceName     string `json:"serviceName,omitempty"`
	ServiceInstance string `json:"serviceInstance,omitempty"`
	Region          string `json:"region,omitempty"`
	ResourceType    string `json:"resourceType,omitempty"`
	Resource        string `json:"resource,omitempty"`
	SpaceID         string `json:"spaceId,omitempty"`
	AccountID       string `json:"accountId,omitempty"`
	OrganizationID  string `json:"organizationId,omitempty"`
	ResourceGroupID string `json:"resourceGroupId,omitempty"`
}
