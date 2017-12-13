package intercom

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestAdminAPIList(t *testing.T) {
	http := TestAdminHTTPClient{fixtureFilename: "fixtures/admins.json", expectedURI: "/admins", t: t}
	api := AdminAPI{httpClient: &http}
	adminList, _ := api.list()
	if adminList.Admins[0].ID != "1" {
		t.Errorf("ID was %s, expected 1", adminList.Admins[0].ID)
	}
}

type TestAdminHTTPClient struct {
	TestHTTPClient
	t               *testing.T
	fixtureFilename string
	expectedURI     string
}

func (t TestAdminHTTPClient) Get(uri string, queryParams interface{}, v interface{}) error {
	if t.expectedURI != uri {
		t.t.Errorf("URI was %s, expected %s", uri, t.expectedURI)
	}
	b, err := ioutil.ReadFile(t.fixtureFilename)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
