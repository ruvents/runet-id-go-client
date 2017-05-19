package client

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

const cfgApiHost = "api.runet-id.com"

type Client struct {
	// private
	apikey  string
	apihash string
	apihost string
}

func NewClient(apikey, secret string) *Client {
	hash := md5.Sum([]byte(apikey + secret))
	return &Client{
		apikey:  apikey,
		apihash: hex.EncodeToString(hash[:]),
		apihost: "runet-id.com",
	}
}

func (client *Client) SetHost(host string) {
	client.apihost = host
}

func (client Client) Request(method string, params RequestParams) (body []byte, err error) {
	var resp *http.Response
	// Подписываем переданный набор параметров запроса реквизитами доступа
	prms := params.ToUrlValues()
	prms.Set("ApiKey", client.apikey)
	prms.Set("Hash", client.apihash)
	// Определяем адрес запроса, который пригодтся для отладочных сообщений
	requestUrl := sprintf("http://%s/%s", cfgApiHost, method)
	// Отправляем запрос
	if resp, err = httpClient.PostForm(requestUrl, prms); err == nil {
		// Читаем содержимое ответа сервера
		if body, err = ioutil.ReadAll(resp.Body); err != nil {
			return body, mkerr("Ошибка чтения содержания ответа %s: %s", requestUrl, err)
		}
		// Проверка ошибки запроса к api
		if gjson.Get(string(body), "Error.Code").Exists() {
			jsonData := gjson.GetMany(string(body), "Error.Code", "Error.Message")
			return nil, mkerr("Ошибка с кодом %d при обращении к %s: %s /////////////// %s",
				uint16(jsonData[0].Num),
				requestUrl,
				jsonData[1].Str,
				tojson(params),
			)
		}
		return
	} else {
		return nil, mkerr("Ошибка отправки запроса")
	}
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
	var body []byte
	if body, err = client.Request("user/get", params); err == nil {
		err = json.Unmarshal(body, &user)
	}
	return
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
