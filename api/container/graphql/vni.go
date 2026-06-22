package graphql

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

// VNI (Virtual Network Interface) related types

// VNIAttachment represents a VNI attachment to a worker node
type VNIAttachment struct {
	ID                      string                  `json:"id,omitempty"`
	ClusterID               string                  `json:"clusterID,omitempty"`
	WorkerID                string                  `json:"workerID,omitempty"`
	VNIID                   string                  `json:"vniID,omitempty"`
	VNIName                 string                  `json:"vniName,omitempty"`
	VlanID                  *int                    `json:"vlanID,omitempty"`
	CanFloat                *bool                   `json:"canFloat,omitempty"`
	Status                  string                  `json:"status,omitempty"`
	CreatedAt               string                  `json:"createdAt,omitempty"`
	PrimaryIPAddress        string                  `json:"primaryIPAddress,omitempty"`
	MACAddress              string                  `json:"macAddress,omitempty"`
	VirtualNetworkInterface VirtualNetworkInterface `json:"virtualNetworkInterface,omitempty"`
	AttachedTo              NetworkAttachable       `json:"attachedTo,omitempty"`
}

// VirtualNetworkInterface represents a VPC VNI
type VirtualNetworkInterface struct {
	ExternalID       string  `json:"externalID"`
	Name             *string `json:"name,omitempty"`
	AutoDelete       *bool   `json:"autoDelete,omitempty"`
	PrimaryIPAddress *string `json:"primaryIPAddress,omitempty"`
	MACAddress       *string `json:"macAddress,omitempty"`
}

// NetworkAttachable represents a node that can have VNI attachments
type NetworkAttachable struct {
	ID string `json:"id"`
}

// VNIAttachmentConnection represents a paginated list of VNI attachments
type VNIAttachmentConnection struct {
	Edges    []VNIAttachmentEdge `json:"edges"`
	PageInfo PageInfo            `json:"pageInfo"`
}

// VNIAttachmentEdge represents an edge in the VNI attachment connection
type VNIAttachmentEdge struct {
	Node VNIAttachment `json:"node"`
}

// PageInfo contains pagination information
type PageInfo struct {
	HasNextPage bool    `json:"hasNextPage"`
	EndCursor   *string `json:"endCursor,omitempty"`
}

// AddVNIToBareMetalNodeInput is the input for attaching a VNI to a bare metal worker
type AddVNIToBareMetalNodeInput struct {
	VirtualNetworkInterfaceID string  `json:"virtualNetworkInterfaceID"`
	VlanID                    int     `json:"vlanID"`
	Node                      *string `json:"node,omitempty"`
	Cluster                   *string `json:"cluster,omitempty"`
	AutoDelete                bool    `json:"autoDelete,omitempty"`
}

// RemoveVNIFromNodeInput is the input for detaching a VNI from a worker
type RemoveVNIFromNodeInput struct {
	VirtualNetworkInterfaceID string  `json:"virtualNetworkInterfaceID"`
	Node                      *string `json:"node,omitempty"`
	Cluster                   *string `json:"cluster,omitempty"`
	AutoDelete                bool    `json:"autoDelete,omitempty"`
}

// AddVNIToBareMetalNodeResponse is the response from attaching a VNI
type AddVNIToBareMetalNodeResponse struct {
	NetworkAttachment VNIAttachment `json:"networkAttachment"`
}

// RemoveVNIFromNodeResponse is the response from detaching a VNI
type RemoveVNIFromNodeResponse struct {
	Cluster                 NetworkAttachable       `json:"cluster"`
	Node                    NetworkAttachable       `json:"node"`
	VirtualNetworkInterface VirtualNetworkInterface `json:"virtualNetworkInterface"`
}

// ListVNIAttachmentsInput is the input for listing VNI attachments
type ListVNIAttachmentsInput struct {
	NodeID string
	First  *int
	After  *string
}

// ListVNIAttachmentsResponse is the response from listing VNI attachments
type ListVNIAttachmentsResponse struct {
	AttachableType string
	Connection     VNIAttachmentConnection
}

// VNIs interface for VNI operations
type VNIs interface {
	AttachToBareMetalNode(input AddVNIToBareMetalNodeInput, target containerv1.ClusterTargetHeader) (*AddVNIToBareMetalNodeResponse, error)
	DetachFromNode(input RemoveVNIFromNodeInput, target containerv1.ClusterTargetHeader) (*RemoveVNIFromNodeResponse, error)
	ListAttachments(input ListVNIAttachmentsInput, target containerv1.ClusterTargetHeader) (*ListVNIAttachmentsResponse, error)
}

type vni struct {
	client *client.Client
}

// NewVNIAPI creates a new VNI API client
func NewVNIAPI(c *client.Client) VNIs {
	return &vni{
		client: c,
	}
}

type graphqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

// AttachToBareMetalNode attaches a Virtual Network Interface to a bare metal worker node
func (r *vni) AttachToBareMetalNode(input AddVNIToBareMetalNodeInput, target containerv1.ClusterTargetHeader) (*AddVNIToBareMetalNodeResponse, error) {
	u, err := url.Parse(*r.client.Config.Endpoint)
	if err != nil {
		return nil, err
	}
	u.Path = "/graphql"

	body, err := json.Marshal(graphqlRequest{
		Query: `
mutation($input: AddVirtualNetworkInterfaceToBareMetalNodeInput!) {
	addVirtualNetworkInterfaceToBareMetalNode(input: $input) {
		networkAttachment {
			attachedTo { id }
			virtualNetworkInterface { 
				externalID 
				name 
				primaryIPAddress
				macAddress
			}
			... on BareMetalNetworkAttachmentByVLAN {
				vlanID
			}
		}
	}
}
`,
		Variables: map[string]any{
			"input": input,
		},
	})
	if err != nil {
		return nil, err
	}

	req := rest.PostRequest(u.String()).Body(body)
	for header, value := range target.ToMap() {
		req.Set(header, value)
	}

	var response struct {
		Data struct {
			AddVirtualNetworkInterfaceToBareMetalNode AddVNIToBareMetalNodeResponse `json:"addVirtualNetworkInterfaceToBareMetalNode"`
		} `json:"data"`
		Errors []struct {
			Message    string `json:"message"`
			Extensions struct {
				Code string `json:"code"`
			} `json:"extensions"`
		} `json:"errors"`
	}

	_, err = r.client.SendRequest(req, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s (code: %s)", response.Errors[0].Message, response.Errors[0].Extensions.Code)
	}

	return &response.Data.AddVirtualNetworkInterfaceToBareMetalNode, nil
}

// DetachFromNode detaches a Virtual Network Interface from a worker node
func (r *vni) DetachFromNode(input RemoveVNIFromNodeInput, target containerv1.ClusterTargetHeader) (*RemoveVNIFromNodeResponse, error) {
	u, err := url.Parse(*r.client.Config.Endpoint)
	if err != nil {
		return nil, err
	}
	u.Path = "/graphql"

	body, err := json.Marshal(graphqlRequest{
		Query: `
mutation($input: RemoveVirtualNetworkInterfaceFromNodeInput!) {
	removeVirtualNetworkInterfaceFromNode(input: $input) {
		cluster { id }
		node { id }
		virtualNetworkInterface { 
			externalID 
			name 
		}
	}
}
`,
		Variables: map[string]any{
			"input": input,
		},
	})
	if err != nil {
		return nil, err
	}

	req := rest.PostRequest(u.String()).Body(body)
	for header, value := range target.ToMap() {
		req.Set(header, value)
	}

	var response struct {
		Data struct {
			RemoveVirtualNetworkInterfaceFromNode RemoveVNIFromNodeResponse `json:"removeVirtualNetworkInterfaceFromNode"`
		} `json:"data"`
		Errors []struct {
			Message    string `json:"message"`
			Extensions struct {
				Code string `json:"code"`
			} `json:"extensions"`
		} `json:"errors"`
	}

	_, err = r.client.SendRequest(req, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s (code: %s)", response.Errors[0].Message, response.Errors[0].Extensions.Code)
	}

	return &response.Data.RemoveVirtualNetworkInterfaceFromNode, nil
}

// ListAttachments lists all Virtual Network Interface attachments for a cluster or worker node
func (r *vni) ListAttachments(input ListVNIAttachmentsInput, target containerv1.ClusterTargetHeader) (*ListVNIAttachmentsResponse, error) {
	u, err := url.Parse(*r.client.Config.Endpoint)
	if err != nil {
		return nil, err
	}
	u.Path = "/graphql"

	variables := map[string]any{
		"id": input.NodeID,
	}
	if input.First != nil {
		variables["first"] = *input.First
	}
	if input.After != nil {
		variables["after"] = *input.After
	}

	body, err := json.Marshal(graphqlRequest{
		Query: `
query($id: ID!, $after: String, $first: Int) {
	   node(id: $id) {
	       ... on KubernetesCluster {
			__typename
			networkAttachments(after: $after, first: $first) {
				pageInfo {
					hasNextPage
					endCursor
				}
				edges {
					node {
						attachedTo {
							id
						}
						virtualNetworkInterface {
							externalID
							name
							autoDelete
							macAddress
							primaryIPAddress
						}
						... on BareMetalNetworkAttachmentByVLAN {
							vlanID
							canFloat
						}
					}
				}
			}
	       }
	       ... on VPCBareMetalKubernetesNode {
			__typename
			networkAttachments(after: $after, first: $first) {
				pageInfo {
					hasNextPage
					endCursor
				}
				edges {
					node {
						attachedTo {
							id
						}
						virtualNetworkInterface {
							externalID
							name
							autoDelete
							macAddress
							primaryIPAddress
						}
						... on BareMetalNetworkAttachmentByVLAN {
							vlanID
							canFloat
						}
					}
				}
			}
	       }
	   }
}
`,
		Variables: variables,
	})
	if err != nil {
		return nil, err
	}

	req := rest.PostRequest(u.String()).Body(body)
	for header, value := range target.ToMap() {
		req.Set(header, value)
	}

	var response struct {
		Data struct {
			Node struct {
				TypeName           string                  `json:"__typename"`
				NetworkAttachments VNIAttachmentConnection `json:"networkAttachments"`
			} `json:"node"`
		} `json:"data"`
		Errors []struct {
			Message    string `json:"message"`
			Extensions struct {
				Code string `json:"code"`
			} `json:"extensions"`
		} `json:"errors"`
	}

	_, err = r.client.SendRequest(req, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s (code: %s)", response.Errors[0].Message, response.Errors[0].Extensions.Code)
	}

	return &ListVNIAttachmentsResponse{
		AttachableType: response.Data.Node.TypeName,
		Connection:     response.Data.Node.NetworkAttachments,
	}, nil
}
