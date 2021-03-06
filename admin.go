package intercom

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Admin represents an Admin in Intercom.
type Admin struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UnmarshalJSON custom SocialProfile unmarshaller
func (a *Admin) UnmarshalJSON(data []byte) error {
	s := struct {
		ID    interface{} `json:"id"`
		Type  string      `json:"type"`
		Name  string      `json:"name"`
		Email string      `json:"email"`
	}{}
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	a.Type = s.Type
	a.Name = s.Name
	a.Email = s.Email

	switch v := s.ID.(type) {
	case string:
		a.ID = v
	case float64:
		a.ID = strconv.FormatFloat(v, 'f', -1, 64)
	}
	return nil
}

// AdminList represents an object holding list of Admins
type AdminList struct {
	Admins []Admin
}

// AdminService handles interactions with the API through an AdminRepository.
type AdminService struct {
	Repository AdminRepository
}

// List lists the Admins associated with your App.
func (c *AdminService) List() (AdminList, error) {
	return c.Repository.list()
}

// IsNobodyAdmin is a helper function to determine if the Admin is 'Nobody'.
func (a Admin) IsNobodyAdmin() bool {
	return a.Type == "nobody_admin"
}

// Get the address for a Contact in order to message them
func (a Admin) MessageAddress() MessageAddress {
	return MessageAddress{
		Type: "admin",
		ID:   a.ID,
	}
}

func (a Admin) String() string {
	return fmt.Sprintf("[intercom] %s { id: %s name: %s, email: %s }", a.Type, a.ID, a.Name, a.Email)
}
