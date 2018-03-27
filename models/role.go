package models

import "github.com/IBM-Bluemix/bluemix-go/crn"

type Role struct {
	CRN         crn.CRN `json:"crn"`
	Name        string  `json:"displayName"`
	Description string  `json:"description"`
}

func (r Role) ToPolicyRole() PolicyRole {
	return PolicyRole{
		ID:          r.CRN,
		DisplayName: r.Name,
		Description: r.Description,
	}
}
