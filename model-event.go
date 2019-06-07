package api

import (
	"encoding/json"
	"net/http"
)

type Event struct {
	ID   uint   `json:"Id"`
	Code string `json:"IdName"`
	Name string `json:"Title"`
}

func (client Client) GetEvents() (events []Event, err error) {
	var body []byte; /**/ if body, err = client.Request(http.MethodPost, "event/list", RequestParams{}); err == nil {
		err = json.Unmarshal(body, &events)
	}
	return
}
