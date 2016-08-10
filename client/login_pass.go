package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

//Данные приложения для windows клиента
//https://new.vk.com/windows_app
//Это приложение имеет права на прямую авторизацию с логином и паролем
const (
	client_id     = "3697615"
	client_secret = "AlVXZFMUqyrnABp8ncuU"
	authHost      = "oauth.vk.com"
	authPath      = "/token"
)

//Client by login/password. Uses client_id of windows_app
//Return client and bool. If login/password good - bool==true
func GetClientWithLoginPassword(username, password string) (*Client, bool) {
	token := GetTokenWithLoginPassword(username, password)
	if token == "" {
		return &Client{}, false
	}
	return DefaultClient(token), true
}


//get token for windows_app by login/password
func GetTokenWithLoginPassword(username, password string) string {
	params := paramsMake(username, password)
	req, err := requestmake(params)
	if err != nil {
		return ""
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	var response ResponseAuthLoginPass
	json.Unmarshal(b, &response)
	return response.AccessToken
}

func paramsMake(username, password string) url.Values {
	params := url.Values{}
	params.Add("client_id", client_id)
	params.Add("client_secret", client_secret)
	params.Add("grant_type", "password")
	params.Add("v", defaultVersion)
	params.Add("username", username)
	params.Add("password", password)
	return params
}

func requestmake(params url.Values) (*http.Request, error) {
	u := url.URL{}
	u.Host = authHost
	u.Scheme = defaultScheme
	u.Path = authPath
	u.RawQuery = params.Encode()
	req, err := http.NewRequest(defaultMethod, u.String(), nil)
	return req, err
}
