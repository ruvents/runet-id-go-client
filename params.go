package api

import (
	"github.com/apex/log"
	"net/url"
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

func (params RequestParams) ToArray() []interface{} {
	values := make([]interface{}, len(params)*2)
	for key, value := range params {
		values = append(values, key, value)
	}
	return values
}
