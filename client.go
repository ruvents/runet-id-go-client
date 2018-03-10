package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	httpClient = &http.Client{Timeout: 18 * time.Second}
	itoa       = strconv.Itoa
	sprintf    = fmt.Sprintf
)

type Client struct {
	// private
	apikey  string
	apihash string
	apihost string
	// Настройки режима отладки
	isVerboseShowResponse bool
}

func NewClient(apikey, secret string) *Client {
	hash := md5.Sum([]byte(apikey + secret))
	return &Client{
		apikey:  apikey,
		apihash: hex.EncodeToString(hash[:]),
		apihost: "api.runet-id.com",
	}
}

func (client *Client) SetHost(host string) *Client {
	client.apihost = host
	return client
}

func (client *Client) SetVerboseShowResponse() *Client {
	client.isVerboseShowResponse = true
	return client
}

func (client Client) Request(method string, params RequestParams) (body []byte, err error) {
	var resp *http.Response
	// Подписываем переданный набор параметров запроса реквизитами доступа
	prms := params.ToUrlValues()
	prms.Set("ApiKey", client.apikey)
	prms.Set("Hash", client.apihash)
	// Определяем адрес запроса, который пригодтся для отладочных сообщений
	requestUrl := sprintf("http://%s/%s", client.apihost, method)
	// Отправляем запрос
	if resp, err = httpClient.PostForm(requestUrl, prms); err == nil {
		// Читаем содержимое ответа сервера
		if body, err = ioutil.ReadAll(resp.Body); err != nil {
			return body, mkerr("Ошибка чтения содержания ответа %s: %s", requestUrl, err)
		}
		// Проверка ошибки запроса к api
		if gjson.Get(string(body), "Error.Code").Exists() {
			jsonData := gjson.GetMany(string(body), "Error.Code", "Error.Message")
			return nil, mkerr("Ошибка с кодом %d при обращении к %s: %s",
				uint16(jsonData[0].Num),
				requestUrl,
				jsonData[1].Str,
			)
		}
		// В режиме отладки отображаем возвращённый контент
		if client.isVerboseShowResponse {
			println(string(body))
		}
		return
	} else {
		return nil, mkerr(fmt.Sprintf("Ошибка отправки запроса: %s", err.Error()))
	}
}

func (client Client) CreateUser(schema User, customizers ... RequestParams) (user User, err error) {
	params := RequestParams{
		"Email":      schema.Email,
		"FirstName":  schema.FirstName,
		"LastName":   schema.LastName,
		"FatherName": schema.FatherName,
		"Phone":      schema.Phone,
		"Company":    schema.Company,
	}
	if len(schema.Attributes) != 0 {
		for param, value := range schema.Attributes {
			params["Attributes["+param+"]"] = value
		}
	}
	for _, customizer := range customizers {
		for param, value := range customizer {
			params[param] = value
		}
	}
	var body []byte; /**/ if body, err = client.Request("user/create", params); err == nil {
		err = json.Unmarshal(body, &user)
	}
	return
}

func (client Client) EditUser(schema User, customizers ... RequestParams) (user User, err error) {
	params := RequestParams{
		"RunetId":    strconv.Itoa(schema.RunetId),
		"Email":      schema.Email,
		"FirstName":  schema.FirstName,
		"LastName":   schema.LastName,
		"FatherName": schema.FatherName,
		"Phone":      schema.Phone,
		"Company":    schema.Company,
	}
	if len(schema.Attributes) != 0 {
		for param, value := range schema.Attributes {
			params["Attributes["+param+"]"] = value
		}
	}
	for _, customizer := range customizers {
		for param, value := range customizer {
			params[param] = value
		}
	}
	var body []byte; /**/ if body, err = client.Request("user/edit", params); err == nil {
		err = json.Unmarshal(body, &user)
	}
	return
}


func (client Client) EventRegister(RunetID int, RoleID int) (user User, err error) {
	params := RequestParams{
		"RunetId":    strconv.Itoa(RunetID),
		"RoleId":     strconv.Itoa(RoleID),
	}
	var body []byte; /**/ if body, err = client.Request("event/register", params); err == nil {
		err = json.Unmarshal(body, &user)
	}
	return
}

func (client Client) GetUser(runetid int) (User, error) {
	return client.GetUserByParams(RequestParams{
		"RunetId": itoa(runetid),
	})
}

func (client Client) GetUserByExternalID(eid string) (User, error) {
	return client.GetUserByParams(RequestParams{
		"ExternalId": eid,
	})
}

func (client Client) GetUserByEmail(email string) (User, error) {
	return client.GetUserByParams(RequestParams{
		"Email": email,
	})
}

func (client Client) GetUserByParams(params RequestParams) (user User, err error) {
	var body []byte /**/
	if body, err = client.Request("user/get", params); err == nil {
		err = json.Unmarshal(body, &user)
	}
	return
}

func (client Client) Basket(idPayer int) (basket Basket, err error) {
	body, err := client.Request("pay/list", RequestParams{
		"PayerRunetId": itoa(idPayer),
	})
	if err == nil {
		err = json.Unmarshal(body, &basket)
	}
	return
}

func (client Client) BasketAdd(idProduct, idPayer, idOwner int) (err error) {
	_, err = client.Request("pay/add", RequestParams{
		"ProductId":    itoa(idProduct),
		"PayerRunetId": itoa(idPayer),
		"OwnerRunetId": itoa(idOwner),
	})
	return
}

func (client Client) BasketUrl(idPayer int) (url string, err error) {
	var body []byte
	body, err = client.Request("pay/url", RequestParams{
		"PayerRunetId": itoa(idPayer),
	})
	return gjson.Get(string(body), "Url").String(), err
}

func mkerr(format string, a ...interface{}) error {
	return errors.New(strings.Replace(fmt.Sprintf(format, a...), "\n", "%5Cn", -1))
}

func tojson(v interface{}) string {
	if text, err := json.Marshal(v); err != nil {
		return "Не могу закодировать значение в JSON..."
	} else {
		return string(text)
	}
}
