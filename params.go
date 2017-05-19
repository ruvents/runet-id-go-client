package client

import (
	"net/url"
	"github.com/apex/log"
)

type RequestParams map[string]string

func (params RequestParams) ToUrlValues() url.Values {
	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}
	return values
}

func (params RequestParams) ToLogFields() log.Fields {
	fields := log.Fields{}
	for key, value := range params {
		fields[key] = value
	}
	return fields
}
