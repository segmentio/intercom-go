package intercom

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/intercom/intercom-go.v2/interfaces"
)

// UserRepository defines the interface for working with Users through the API.
type UserRepository interface {
	find(UserIdentifiers) (*User, error)
	list(userListParams) (*UserList, error)
	scroll(scrollParam string) (*UserList, error)
	save(*User) (User, error)
	delete(id string) (User, error)
}

// UserAPI implements UserRepository
type UserAPI struct {
	httpClient interfaces.HTTPClient
}

type requestScroll struct {
	ScrollParam string `json:"scroll_param,omitempty"`
}
type requestUser struct {
	ID                     string                 `json:"id,omitempty"`
	Email                  string                 `json:"email,omitempty"`
	Phone                  string                 `json:"phone,omitempty"`
	UserID                 string                 `json:"user_id,omitempty"`
	Name                   string                 `json:"name,omitempty"`
	SignedUpAt             int64                  `json:"signed_up_at,omitempty"`
	RemoteCreatedAt        int64                  `json:"remote_created_at,omitempty"`
	LastRequestAt          int64                  `json:"last_request_at,omitempty"`
	LastSeenIP             string                 `json:"last_seen_ip,omitempty"`
	UnsubscribedFromEmails *bool                  `json:"unsubscribed_from_emails,omitempty"`
	Companies              []UserCompany          `json:"companies,omitempty"`
	CustomAttributes       map[string]interface{} `json:"custom_attributes,omitempty"`
	UpdateLastRequestAt    *bool                  `json:"update_last_request_at,omitempty"`
	NewSession             *bool                  `json:"new_session,omitempty"`
	LastSeenUserAgent      string                 `json:"last_seen_user_agent,omitempty"`
}

func (api UserAPI) find(params UserIdentifiers) (*User, error) {
	var user *User
	err := api.getClientForFind(params, &user)
	return user, err
}

func (api UserAPI) getClientForFind(params UserIdentifiers, v interface{}) error {
	switch {
	case params.ID != "":
		return api.httpClient.Get(fmt.Sprintf("/users/%s", params.ID), nil, v)
	case params.UserID != "", params.Email != "":
		return api.httpClient.Get("/users", params, v)
	}
	return errors.New("Missing User Identifier")
}

func (api UserAPI) list(params userListParams) (*UserList, error) {
	var list *UserList
	err := api.httpClient.Get("/users", params, &list)
	return list, err
}

func (api UserAPI) scroll(scrollParam string) (*UserList, error) {
	var list *UserList
	params := scrollParams{ScrollParam: scrollParam}
	err := api.httpClient.Get("/users/scroll", params, &list)
	return list, err
}

func (api UserAPI) save(user *User) (User, error) {
	return unmarshalToUser(api.httpClient.Post("/users", RequestUserMapper{}.ConvertUser(user)))
}

func unmarshalToUser(data []byte, err error) (User, error) {
	savedUser := User{}
	if err != nil {
		return savedUser, err
	}
	err = json.Unmarshal(data, &savedUser)
	return savedUser, err
}

func (api UserAPI) delete(id string) (User, error) {
	user := User{}
	data, err := api.httpClient.Delete(fmt.Sprintf("/users/%s", id), nil)
	if err != nil {
		return user, err
	}
	err = json.Unmarshal(data, &user)
	return user, err
}
