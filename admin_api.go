package intercom

import (
	"gopkg.in/intercom/intercom-go.v2/interfaces"
)

// AdminRepository defines the interface for working with Admins through the API.
type AdminRepository interface {
	list() (*AdminList, error)
}

// AdminAPI implements AdminRepository
type AdminAPI struct {
	httpClient interfaces.HTTPClient
}

func (api AdminAPI) list() (*AdminList, error) {
	var list *AdminList
	err := api.httpClient.Get("/admins", nil, &list)
	return list, err
}
