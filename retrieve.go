package pkt

import (
	"sort"
)

func (c *Client) Retrieve(after string, offset int) ([]Item, error) {
	action := "/get"
	request := retrieveRequest{
		Auth: *c.auth,
		retrieveOptions: retrieveOptions{
			Since:      after,
			Offset:     offset,
			DetailType: "complete",
			Sort:       "newest",
			Count:      PageCount,
		},
	}
	response := retrieveResponse{}

	err := c.post(action, request, &response)

	return response.Items(), err
}

func (c *Client) RetrieveAll(after string) ([]Item, error) {
	items := []Item{}
	offset := 0
	for {
		retrieved, err := c.Retrieve(after, offset)
		if err != nil {
			return nil, err
		}

		if len(retrieved) == 0 {
			break
		}

		items = append(items, retrieved...)
		offset = offset + PageCount
	}

	return items, nil
}

type retrieveOptions struct {
	DetailType string `json:"detailType"`
	Since      string `json:"since"`
	Offset     int    `json:"offset"`
	Sort       string `json:"sort"`
	Count      int    `json:"count"`
}

type retrieveRequest struct {
	Auth
	retrieveOptions
}

type retrieveResponse struct {
	Status int             `json:"status"`
	List   map[string]Item `json:"list"`
}

func (r retrieveResponse) Items() []Item {
	var items []Item

	for _, item := range r.List {
		items = append(items, item)
	}

	sort.Slice(items[:], func(i, j int) bool {
		return items[i].AddedAt >= items[j].AddedAt
	})

	return items
}
