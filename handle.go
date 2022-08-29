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
	id          int
	AppKey      string
	AppSecret   string
	accessToken string
}

var AllHandle []*Handle

// NewHandle 创建动作执行器
func NewHandle(appKey string, appSecret string) *Handle {
	var h Handle
	h.id = len(AllHandle)
	h.AppKey = appKey
	h.AppSecret = appSecret
	AllHandle = append(AllHandle, &h)
	return &h
}

// GetID 获取当前命令执行器的id
func (h *Handle) GetID() int {
	return h.id
}

// GetAccessToken 获取token
func (h *Handle) GetAccessToken() (err error) {
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

//请求函数
func (h *Handle) req(url string, msg *bytes.Reader) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, msg)
	if err != nil {
		return "", err
	}

	req.Header.Set("x-acs-dingtalk-access-token", h.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
