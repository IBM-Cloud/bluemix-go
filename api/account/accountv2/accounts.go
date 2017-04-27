package accountv2

import (
	"fmt"

	bluemix "github.com/IBM-Bluemix/bluemix-go"
	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
)

//Account Model ...
type Account struct {
	GUID          string
	Name          string
	Type          string
	State         string
	OwnerGUID     string
	OwnerUserID   string
	OwnerUniqueID string
	CustomerID    string
	CountryCode   string
	CurrencyCode  string
	Organizations []AccountOrganization
	Members       []AccountMember `json:"members"`
}

//AccountOrganization ...
type AccountOrganization struct {
	GUID   string `json:"guid"`
	Region string `json:"region"`
}

//AccountMember ...
type AccountMember struct {
	GUID     string `json:"guid"`
	UserID   string `json:"user_id"`
	UniqueID string `json:"unique_id"`
}

//AccountResource ...
type AccountResource struct {
	Resource
	Entity AccountEntity
}

//AccountEntity ...
type AccountEntity struct {
	Name          string                `json:"name"`
	Type          string                `json:"type"`
	State         string                `json:"state"`
	OwnerGUID     string                `json:"owner"`
	OwnerUserID   string                `json:"owner_userid"`
	OwnerUniqueID string                `json:"owner_unique_id"`
	CustomerID    string                `json:"customer_id"`
	CountryCode   string                `json:"country_code"`
	CurrencyCode  string                `json:"currency_code"`
	Organizations []AccountOrganization `json:"organizations_region"`
	Members       []AccountMember       `json:"members"`
}

//ToModel ...
func (resource AccountResource) ToModel() Account {
	entity := resource.Entity

	return Account{
		GUID:          resource.Metadata.GUID,
		Name:          entity.Name,
		Type:          entity.Type,
		State:         entity.State,
		OwnerGUID:     entity.OwnerGUID,
		OwnerUserID:   entity.OwnerUserID,
		OwnerUniqueID: entity.OwnerUniqueID,
		CustomerID:    entity.CustomerID,
		CountryCode:   entity.CountryCode,
		CurrencyCode:  entity.CurrencyCode,
		Organizations: entity.Organizations,
		Members:       entity.Members,
	}
}

//AccountQueryResponse ...
type AccountQueryResponse struct {
	Metadata Metadata
	Accounts []AccountResource `json:"resources"`
}

//Accounts ...
type Accounts interface {
	FindByOrg(orgGUID string) (*Account, error)
}

type account struct {
	client *AccountClient
	config *bluemix.Config
}

func newAccountAPI(c *AccountClient) Accounts {
	return &account{
		client: c,
		config: c.config,
	}
}

//FindByOrg ...
func (r *account) FindByOrg(orgGUID string) (*Account, error) {
	region := r.config.Region
	type organizationRegion struct {
		GUID   string `json:"guid"`
		Region string `json:"region"`
	}

	payLoad := struct {
		OrganizationsRegion []organizationRegion `json:"organizations_region"`
	}{
		OrganizationsRegion: []organizationRegion{
			{
				GUID:   orgGUID,
				Region: region,
			},
		},
	}

	queryResp := AccountQueryResponse{}
	response, err := r.client.post("/coe/v2/getaccounts", payLoad, &queryResp)
	if err != nil {

		if response.StatusCode == 404 {
			return nil, bmxerror.New(ErrCodeNoAccountExists,
				fmt.Sprintf("No account exists in the given region: %q and the given org: %q", region, orgGUID))
		}
		return nil, err

	} else if len(queryResp.Accounts) > 0 {
		account := queryResp.Accounts[0].ToModel()
		return &account, nil
	}
	return nil, nil
}
