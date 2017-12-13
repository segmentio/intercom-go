package intercom

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/intercom/intercom-go.v2/interfaces"
)

// CompanyRepository defines the interface for working with Companies through the API.
type CompanyRepository interface {
	find(CompanyIdentifiers) (*Company, error)
	list(companyListParams) (*CompanyList, error)
	scroll(scrollParam string) (*CompanyList, error)
	save(*Company) (Company, error)
}

// CompanyAPI implements CompanyRepository
type CompanyAPI struct {
	httpClient interfaces.HTTPClient
}

type requestCompany struct {
	ID               string                 `json:"id,omitempty"`
	CompanyID        string                 `json:"company_id,omitempty"`
	Name             string                 `json:"name,omitempty"`
	RemoteCreatedAt  int64                  `json:"remote_created_at,omitempty"`
	MonthlySpend     int64                  `json:"monthly_spend,omitempty"`
	Plan             string                 `json:"plan,omitempty"`
	CustomAttributes map[string]interface{} `json:"custom_attributes,omitempty"`
}

func (api CompanyAPI) find(params CompanyIdentifiers) (*Company, error) {
	var company *Company
	err := api.getClientForFind(params, &company)
	return company, err
}

func (api CompanyAPI) getClientForFind(params CompanyIdentifiers, v interface{}) error {
	switch {
	case params.ID != "":
		return api.httpClient.Get(fmt.Sprintf("/companies/%s", params.ID), nil, v)
	case params.CompanyID != "", params.Name != "":
		return api.httpClient.Get("/companies", params, v)
	}
	return errors.New("Missing Company Identifier")
}

func (api CompanyAPI) list(params companyListParams) (*CompanyList, error) {
	var list *CompanyList
	err := api.httpClient.Get("/companies", params, &list)
	return list, err
}

func (api CompanyAPI) scroll(scrollParam string) (*CompanyList, error) {
	var list *CompanyList
	params := scrollParams{ScrollParam: scrollParam}
	err := api.httpClient.Get("/companies/scroll", params, &list)
	return list, err
}

func (api CompanyAPI) save(company *Company) (Company, error) {
	requestCompany := requestCompany{
		ID:               company.ID,
		Name:             company.Name,
		CompanyID:        company.CompanyID,
		RemoteCreatedAt:  company.RemoteCreatedAt,
		MonthlySpend:     company.MonthlySpend,
		Plan:             api.getPlanName(company),
		CustomAttributes: company.CustomAttributes,
	}

	savedCompany := Company{}
	data, err := api.httpClient.Post("/companies", &requestCompany)
	if err != nil {
		return savedCompany, err
	}
	err = json.Unmarshal(data, &savedCompany)
	return savedCompany, err
}

func (api CompanyAPI) getPlanName(company *Company) string {
	if company.Plan == nil {
		return ""
	}
	return company.Plan.Name
}
