package api

import (
	"strings"
)

type User struct {
	RunetID    int    `json:"RunetId"`
	ExternalID string `json:"ExternalId"`
	Email      string
	LastName   string `json:"LastName"`
	FirstName  string `json:"FirstName"`
	Company    string `json:"Company"`
	Position   string `json:"Position"`
	Attributes map[string]string
}

func NewUser() User {
	return User{
		Attributes: map[string]string{},
	}
}

func (user User) GetFullName() string {
	return strings.Trim(user.LastName+" "+user.FirstName, " ")
}

func (user User) CreateHidden(api *Client) error {
	params := struct2map(&user)
	params["Visible"] = "0"
	_, err := api.Request("user/create", params)
	return err
}

func (user User) Create(api *Client) error {
	_, err := api.Request("user/create", struct2map(&user))
	return err
}

func (user User) Update(api *Client) error {
	_, err := api.Request("user/edit", struct2map(&user))
	return err
}

func (user User) Register(api *Client, roleid int) error {
	_, err := api.Request("event/register", RequestParams{
		"RoleId":  itoa(roleid),
		"RunetId": itoa(user.RunetID),
	})
	return err
}
