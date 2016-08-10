package anticaptcha

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	ac_client "github.com/FireGM/anti-captcha.com/client"
	"github.com/FireGM/vk/client"
)

type ACClient struct {
	client.Client
	AcClient ac_client.Client
}

func (c *ACClient) ValuesForParse(command string, values url.Values) ([]byte, error) {
	var ner client.ErrorResponse

	req, err := c.MakeRequest(command, values)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.Do(req)
	if err != nil {
		return []byte{}, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.ValuesForParse(command, values)
	}
	json.Unmarshal(b, &ner)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	return b, nil
}

func (c *ACClient) Execute(code string) (b []byte, err error) {
	values := url.Values{}
	values.Set("code", code)
	b, err = c.ValuesForParse("execute", values)
	return b, err
}

func DownloadCaptcha(url string) (string, error) {
	path := "captcha/"
	os.MkdirAll(path, 0777)
	filename := path + "captcha.jpg"
	out, _ := os.Create(filename)
	defer out.Close()
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	io.Copy(out, resp.Body)
	return filename, nil
}

func (c *ACClient) Do(req *http.Request) (res *http.Response, err error) {
	c.LimitRate()
	var ner client.ErrorResponse
	//	request_count++
	//	fmt.Println(request_count)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return res, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, err
	}
	json.Unmarshal(b, &ner)
	if ner.Err.Code == 6 {
		return c.Do(req)
	} else if ner.Err.Code == 14 {
		var captcha ResponseCaptcha
		json.Unmarshal(b, &captcha)
		text := ""
			path, err := DownloadCaptcha(captcha.Error.Captcha_img)
			if err != nil {
				return res, err
			}
			text, err = c.AcClient.SendAndGet(path)
			if err != nil {
				return res, err
			}
			fmt.Println(text)
		values := req.URL.Query()
		values.Add("captcha_sid", captcha.Error.Captcha_sid)
		values.Add("captcha_key", text)
		req.URL.RawQuery = values.Encode()
		return c.Do(req)
	}
	if err != nil {
		fmt.Println(err)
		return res, err
	}
	return res, nil
}

func GetACClient(token, anticaptcha_key string) *ACClient {
	cl := &client.DefaultClient(token)
	acClient := ac_client.GetClient(anticaptcha_key)
	return &ACClient{Client: cl, AcClient: acClient}
}

func GetACClientWithLoginPassword(username, password, api_key string) (*ACClient, bool) {
	token := client.GetTokenWithLoginPassword(username, password)
	if token == "" {
		return &ACClient{}, false
	}
	return GetACClient(token, api_key), true
}