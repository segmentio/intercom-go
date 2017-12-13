package intercom

import (
	"fmt"

	"gopkg.in/intercom/intercom-go.v2/interfaces"
)

// SegmentRepository defines the interface for working with Segments through the API.
type SegmentRepository interface {
	list() (*SegmentList, error)
	find(id string) (*Segment, error)
}

// SegmentAPI implements SegmentRepository
type SegmentAPI struct {
	httpClient interfaces.HTTPClient
}

func (api SegmentAPI) list() (*SegmentList, error) {
	var list *SegmentList
	err := api.httpClient.Get("/segments", nil, &list)
	return list, err
}

func (api SegmentAPI) find(id string) (*Segment, error) {
	var segment *Segment
	err := api.httpClient.Get(fmt.Sprintf("/segments/%s", id), nil, &segment)
	return segment, err
}
