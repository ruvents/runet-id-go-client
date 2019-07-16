package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/parnurzeal/gorequest"
	"github.com/ruvents/runet-id-go-client/models"
	"github.com/ruvents/tools"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	itoa          = strconv.Itoa
	sprintf       = fmt.Sprintf
	httpClientNew = gorequest.New()
	httpClient    = &http.Client{
		Timeout:   18 * time.Second,
		Transport: &http.Transport{TLSHandshakeTimeout: 5 * time.Second, Dial: (&net.Dialer{Timeout: 5 * time.Second}).Dial},
	}
)

type Client struct {
	// private
	apikey  string
	apihash string
	apihost string
	// Настройки режима отладки
	isVerboseShowResponse bool
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewClient(apikey, secret string) *Client {
	return &Client{
		apikey:  apikey,
		apihash: tools.MD5(apikey + secret),
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

func (client Client) RequestNew(v interface{}, request *Request) (err error) {
	var clientResponse *http.Response
	var clientRequest *http.Request
	// Готовим подпись реквизитами доступа создаваемый запрос
	queryParams := url.Values{}
	queryParams.Set("ApiKey", client.apikey)
	queryParams.Set("Hash", client.apihash)
	// Определяем адрес запроса, который пригодтся для отладочных сообщений
	requestURL := url.URL{
		Host:     client.apihost + ":9000",
		Path:     string(request.path) + "/",
		Scheme:   "http",
		RawQuery: queryParams.Encode(),
	}
	// Конструируем запрос
	clientRequest, err = http.NewRequest(string(request.kind), requestURL.String(), nil)
	// Отправляем запрос
	if clientResponse, err = httpClient.Do(clientRequest); err == nil {
		// Читаем содержимое ответа сервера
		var body []byte; /**/ if body, err = ioutil.ReadAll(clientResponse.Body); err != nil {
			return mkerr("Ошибка чтения содержания ответа %s: %s", requestURL, err)
		}
		// В режиме отладки отображаем возвращённый контент
		if client.isVerboseShowResponse {
			log.WithField("RequestURL", requestURL.String()).
				WithField("ResponseBODY", string(body)).
				Info("Запрос к API")
		}
		// Пробуем парсить ответ сервера в необходимую структуру данных
		if err = json.Unmarshal(body, &v); err != nil {
			// Не смогли декодировать содержимое ответа, а не стандартная ли ошибка api пришла?
			var errorResponse ErrorResponse; /**/ if err = json.Unmarshal(body, &errorResponse); err == nil {
				log.Errorf("Ошибка на стороне API: %s", errorResponse.Message)
			} else {
				log.Errorf("Ошибка на стороне API (нестандартная): %s", string(body))
			}
		}
		// 	// Проверка ошибки запроса к api
		// 	if gjson.Get(string(body), "Error.Code").Exists() {
		// 		jsonData := gjson.GetMany(string(body), "Error.Code", "Error.Message")
		// 		return nil, mkerr("Ошибка с кодом %d при обращении к %s: %s",
		// 			uint16(jsonData[0].Num),
		// 			requestUrl,
		// 			jsonData[1].Str,
		// 		)
		// 	}
		return
	} else {
		return mkerr(fmt.Sprintf("Ошибка отправки запроса: %s", err.Error()))
	}

	return
}

func (client Client) Request(method string, path string, params RequestParams) (body []byte, err error) {
	var resp *http.Response
	// Подписываем переданный набор параметров запроса реквизитами доступа
	params["ApiKey"] = client.apikey
	params["Hash"] = client.apihash
	// Определяем адрес запроса, который пригодтся для отладочных сообщений
	requestUrl := sprintf("http://%s/%s/", client.apihost, path)
	// Отправляем запрос
	if resp, err = httpClient.PostForm(requestUrl, params.ToUrlValues()); err == nil {
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

func (client Client) CreateUser(schema User, customizers ...RequestParams) (user User, err error) {
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
	var body []byte; /**/ if body, err = client.Request(http.MethodPost, "user/create", params); err == nil {
		err = json.Unmarshal(body, &user)
	}
	return
}

func (client Client) EditUserNew(schema models.User, customizers ...RequestParams) (user User, err error) {
	params := RequestParams{
		"RunetId":    strconv.Itoa(int(schema.ID)),
		"Email":      schema.Email,
		"FirstName":  schema.FirstName,
		"LastName":   schema.LastName,
		"FatherName": schema.FatherName,
		"Phone":      schema.Phone,
		"Company":    schema.Work.Company,
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
	var body []byte; /**/ if body, err = client.Request(http.MethodPost, "user/edit", params); err == nil {
		err = json.Unmarshal(body, &user)
	}
	return
}

func (client Client) EditUser(schema User, customizers ...RequestParams) (user User, err error) {
	params := RequestParams{
		"RunetId":    strconv.Itoa(int(schema.RunetId)),
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
	var body []byte; /**/ if body, err = client.Request(http.MethodPost, "user/edit", params); err == nil {
		err = json.Unmarshal(body, &user)
	}
	return
}

func (client Client) EventHandoutGet(category string, runetid uint32) (handout Handout, err error) {
	params := RequestParams{
		"Category": category,
		"RunetId":  fmt.Sprintf("%d", runetid),
	}
	var body []byte; /**/ if body, err = client.Request(http.MethodPost, "handout/get", params); err == nil {
		err = json.Unmarshal(body, &handout)
	}
	return
}

func (client Client) EventParticipants() (users []models.User, err error) {
	err = client.RequestNew(&users, &Request{
		kind: http.MethodGet,
		path: PathEventParticipants,
	})

	return
}

func (client Client) EventParticipantsUnsafe() (users []models.User) {
	client.RequestNew(&users, &Request{
		kind: http.MethodGet,
		path: PathEventParticipants,
	})

	return
}

func (client Client) EventParticipantsTS18() (users []User, err error) {
	type UserNew struct {
		RunetId    int    `json:"id"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		LastName   string `json:"lastName"`
		FirstName  string `json:"firstName"`
		FatherName string `json:"fatherName"`
		Attributes map[string]string
	}

	var response []UserNew; /**/ httpClientNew.Get("http://runet-id.com:8080/event/participants").EndStruct(&response)
	for _, user := range response {
		users = append(users, User{
			RunetId:    uint32(user.RunetId),
			Email:      user.Email,
			Phone:      user.Phone,
			LastName:   user.LastName,
			FirstName:  user.FirstName,
			FatherName: user.FatherName,
			Attributes: user.Attributes,
		})
	}

	return
}

func (client Client) EventRegister(RunetID uint32, RoleID uint32) (user User, err error) {
	params := RequestParams{
		"RunetId": tools.MakeStringFromUINT32(RunetID),
		"RoleId":  tools.MakeStringFromUINT32(RoleID),
	}
	var body []byte; /**/ if body, err = client.Request(http.MethodPost, "event/register", params); err == nil {
		err = json.Unmarshal(body, &user)
	}
	return
}

func (client Client) EventUnregister(RunetID uint32) (user User, err error) {
	params := RequestParams{
		"RunetId": tools.MakeStringFromUINT32(RunetID),
	}
	var body []byte; /**/ if body, err = client.Request(http.MethodPost, "event/unregister", params); err == nil {
		err = json.Unmarshal(body, &user)
	}
	return
}

func (client Client) GetUser(runetid uint32) (User, error) {
	return client.GetUserByParams(RequestParams{
		"RunetId": tools.MakeStringFromUINT32(runetid),
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
	if body, err = client.Request(http.MethodPost, "user/get", params); err == nil {
		err = json.Unmarshal(body, &user)
	}
	return
}

func (client Client) GetPaperlessOstrovAssociations(deviceNumber uint16, araisedAt int64) (association models.PaperlessOstrovAssociation, err error) {
	body, err := client.Request(http.MethodPost, "paperless/getHall", RequestParams{
		"DeviceID":             tools.MakeStringFromUINT16(deviceNumber),
		"AraisedAt":            tools.MakeStringFromINT64(araisedAt),
	})
	if err == nil {
		err = json.Unmarshal(body, &association)
	}
	return
}

func (client Client) Basket(idPayer uint32) (basket Basket, err error) {
	body, err := client.Request(http.MethodPost, "pay/list", RequestParams{
		"PayerRunetId": tools.MakeStringFromUINT32(idPayer),
	})
	if err == nil {
		err = json.Unmarshal(body, &basket)
	}
	return
}

func (client Client) BasketAdd(idProduct int, idPayer, idOwner uint32) (err error) {
	_, err = client.Request(http.MethodPost, "pay/add", RequestParams{
		"ProductId":    itoa(idProduct),
		"PayerRunetId": tools.MakeStringFromUINT32(idPayer),
		"OwnerRunetId": tools.MakeStringFromUINT32(idOwner),
	})
	return
}

func (client Client) BasketUrl(idPayer uint32) (url string, err error) {
	var body []byte
	body, err = client.Request(http.MethodPost, "pay/url", RequestParams{
		"PayerRunetId": tools.MakeStringFromUINT32(idPayer),
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
