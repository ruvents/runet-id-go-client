package api

import (
	"strings"
	"strconv"
	"net/http"
)

type User struct {
	RunetId    uint32 `json:"RunetId"`
	ExternalID string `json:"ExternalId"`
	Email      string
	Phone      string
	LastName   string `json:"LastName"`
	FirstName  string `json:"FirstName"`
	FatherName string `json:"FatherName"`
	Company    string `json:"Company"`
	Position   string `json:"Position"`
	Attributes map[string]string
}

func NewUser(schema ... User) User {
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
	_, err := api.Request(http.MethodPost, "user/create", params)
	return err
}

func (user User) Create(api *Client) error {
	_, err := api.Request(http.MethodPost, "user/create", struct2map(&user))
	return err
}

func (user User) Update(api *Client) error {
	_, err := api.EditUser(user)
	return err
}

func (user User) Register(api *Client, roleid int) error {
	_, err := api.Request(http.MethodPost, "event/register", RequestParams{
		"RoleId":  itoa(roleid),
		"RunetId": strconv.FormatUint(uint64(user.RunetId), 10),
	})
	return err
}
