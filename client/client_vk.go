package client

/*
Package for work with api of vk.com
*/

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
	"time"
)

const (
	paramToken           = "access_token"
	paramVersion         = "v"
	defaultHost          = "api.vk.com"
	defaultPath          = "/method/"
	defaultScheme        = "https"
	defaultVersion       = "5.50"
	defaultMethod        = "POST"
	maxRequestsPerSecond = 1
	minimumRate          = time.Second / maxRequestsPerSecond
)

//client for requesting to vk api
type Client struct {
	lastRequest time.Time
	m           sync.Mutex
	minimumRate time.Duration
	Token       string
}

//check exist user and token
func (c *Client) Check() bool {
	req, _ := c.MakeRequest("users.get", url.Values{})
	b, err := c.DoBytes(req)
	if err != nil {
		return false
	}
	var user userCheckResponse
	json.Unmarshal(b, &user)
	if len(user.Response) < 1 {
		return false
	}
	if user.Response[0].Id < 1 {
		return false
	}
	return true
}

//read and close response from api vk.com
//return slice of byte
func (c *Client) DoBytes(req *http.Request) ([]byte, error) {
	res, err := c.Do(req)
	var ner ErrorResponse
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	json.Unmarshal(b, &ner)
	if ner.Err.Code == 6 {
		return c.DoBytes(req)
	} else if ner.Err.Code == 14 {
		var captcha ResponseCaptcha
		json.Unmarshal(b, &captcha)
		text := ""
		fmt.Println(captcha.Error.Captcha_img)
		open.Run(captcha.Error.Captcha_img[:len(captcha.Error.Captcha_img)-4])
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: \n")
		text, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		text = text[:len(text)-1]
		values := req.URL.Query()
		values.Add("captcha_sid", captcha.Error.Captcha_sid)
		values.Add("captcha_key", text)
		req.URL.RawQuery = values.Encode()
		return c.DoBytes(req)
	}

	return b, nil
}

// limitRate limits request rate
func (c *Client) LimitRate() {
	c.m.Lock()
	defer c.m.Unlock()
	now := time.Now()
	//	fmt.Println(now, c.lastRequest)
	diff := now.Sub(c.lastRequest)
	if diff < minimumRate {
		time.Sleep(minimumRate - diff)
	}
	c.lastRequest = now
}

// Create request
// @name string - method of api
// @params url.Values - params for method
func (c *Client) MakeRequest(name string, params url.Values) (req *http.Request, err error) {
	params.Add(paramVersion, defaultVersion)
	if c.Token != "" {
		params.Add(paramToken, c.Token)
	}
	u := url.URL{}
	u.Host = defaultHost
	u.Scheme = defaultScheme
	u.Path = path.Join(defaultPath, name)
	//	u.RawQuery = parms.Encode()
	req, err = http.NewRequest(defaultMethod, u.String(), bytes.NewBufferString(params.Encode()))
	return req, err
}

//low-level response *http.Response
func (c *Client) Do(req *http.Request) (res *http.Response, err error) {
	c.LimitRate()
	res, err = http.DefaultClient.Do(req)

	return 
}

//Sugar for method "execute"
//more read https://vk.com/dev/execute
func (c *Client) Execute(code string) (b []byte) {
	var ner ErrorResponse
	values := url.Values{}
	values.Set("code", code)
	req, _ := c.MakeRequest("execute", values)
	res, doErr := c.Do(req)
	if doErr != nil {
		return []byte{}
	}
	b, _ = ioutil.ReadAll(res.Body)
	if !res.Close {
		res.Body.Close()
	}
	err := json.Unmarshal(b, &ner)
	if ner.Err.Code == 6 {
		return c.Execute(code)
	}
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	return b
}

//Return default client with RPS 3
func DefaultClient(token string) *Client {
	return &Client{Token: token, minimumRate: minimumRate}
}
