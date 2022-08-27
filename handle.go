package BeaverBot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
)

type Handle struct {
	AppKey      string
	AppSecret   string
	accessToken string
}

// GetAccessToken 获取token
func (h Handle) GetAccessToken() (err error) {
	body := make(map[string]interface{})
	body["appKey"] = h.AppKey
	body["appSecret"] = h.AppSecret
	bytesData, err := json.Marshal(body)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bytesData)
	url := "https://api.dingtalk.com/v1.0/oauth2/accessToken"
	request, err := http.Post(url, "application/json", reader)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(request.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(request.Body)
	if err != nil {
		return err
	}
	if gjson.Get(string(data), "expireIn").String() != "7200" {
		return errors.New(string(data))
	}
	h.accessToken = gjson.Get(string(data), "accessToken").String()
	return nil
}
