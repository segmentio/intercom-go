package intercom

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/intercom/intercom-go.v2/interfaces"
)

// ContactRepository defines the interface for working with Contacts through the API.
type ContactRepository interface {
	find(UserIdentifiers) (*Contact, error)
	list(contactListParams) (*ContactList, error)
	scroll(scrollParam string) (*ContactList, error)
	create(*Contact) (Contact, error)
	update(*Contact) (Contact, error)
	convert(*Contact, *User) (User, error)
	delete(id string) (Contact, error)
}

// ContactAPI implements ContactRepository
type ContactAPI struct {
	httpClient interfaces.HTTPClient
}

func (api ContactAPI) find(params UserIdentifiers) (*Contact, error) {
	var contact *Contact
	err := api.getClientForFind(params, &contact)
	return contact, err
}

func (api ContactAPI) getClientForFind(params UserIdentifiers, v interface{}) error {
	switch {
	case params.ID != "":
		return api.httpClient.Get(fmt.Sprintf("/contacts/%s", params.ID), nil, v)
	case params.UserID != "":
		return api.httpClient.Get("/contacts", params, v)
	}
	return errors.New("Missing Contact Identifier")
}

func (api ContactAPI) list(params contactListParams) (*ContactList, error) {
	var list *ContactList
	err := api.httpClient.Get("/contacts", params, &list)
	return list, err
}

func (api ContactAPI) scroll(scrollParam string) (*ContactList, error) {
	var list *ContactList
	params := scrollParams{ScrollParam: scrollParam}
	err := api.httpClient.Get("/contacts/scroll", params, &list)
	return list, err
}

func (api ContactAPI) create(contact *Contact) (Contact, error) {
	requestContact := api.buildRequestContact(contact)
	return unmarshalToContact(api.httpClient.Post("/contacts", &requestContact))
}

func (api ContactAPI) update(contact *Contact) (Contact, error) {
	requestContact := api.buildRequestContact(contact)
	return unmarshalToContact(api.httpClient.Post("/contacts", &requestContact))
}

func (api ContactAPI) convert(contact *Contact, user *User) (User, error) {
	cr := convertRequest{Contact: api.buildRequestContact(contact), User: requestUser{
		ID:         user.ID,
		UserID:     user.UserID,
		Email:      user.Email,
		SignedUpAt: user.SignedUpAt,
	}}
	return unmarshalToUser(api.httpClient.Post("/contacts/convert", &cr))
}

func (api ContactAPI) delete(id string) (Contact, error) {
	contact := Contact{}
	data, err := api.httpClient.Delete(fmt.Sprintf("/contacts/%s", id), nil)
	if err != nil {
		return contact, err
	}
	err = json.Unmarshal(data, &contact)
	return contact, err
}

type convertRequest struct {
	User    requestUser `json:"user"`
	Contact requestUser `json:"contact"`
}

func unmarshalToContact(data []byte, err error) (Contact, error) {
	savedContact := Contact{}
	if err != nil {
		return savedContact, err
	}
	err = json.Unmarshal(data, &savedContact)
	return savedContact, err
}

func (api ContactAPI) buildRequestContact(contact *Contact) requestUser {
	return requestUser{
		ID:                     contact.ID,
		Email:                  contact.Email,
		Phone:                  contact.Phone,
		UserID:                 contact.UserID,
		Name:                   contact.Name,
		LastRequestAt:          contact.LastRequestAt,
		LastSeenIP:             contact.LastSeenIP,
		UnsubscribedFromEmails: contact.UnsubscribedFromEmails,
		Companies:              api.getCompaniesToSendFromContact(contact),
		CustomAttributes:       contact.CustomAttributes,
		UpdateLastRequestAt:    contact.UpdateLastRequestAt,
		NewSession:             contact.NewSession,
	}
}

func (api ContactAPI) getCompaniesToSendFromContact(contact *Contact) []UserCompany {
	if contact.Companies == nil {
		return []UserCompany{}
	}
	return RequestUserMapper{}.MakeUserCompaniesFromCompanies(contact.Companies.Companies)
}
