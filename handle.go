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
func NewHandle(appKey string, appSecret string) (*Handle, error) {
	var h Handle
	h.id = len(AllHandle)
	h.AppKey = appKey
	h.AppSecret = appSecret
	err := h.getAccessToken()
	if err != nil {
		return nil, err
	}
	AllHandle = append(AllHandle, &h)
	return &h, err
}

// GetID 获取当前命令执行器的id
func (h *Handle) GetID() int {
	return h.id
}

// GetAccessToken 获取token
func (h *Handle) getAccessToken() (err error) {
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
func (h *Handle) req(url string, msg *bytes.Reader) (statusCode int, bodyData string, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, msg)
	if err != nil {
		return 0, "", err
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
		return 0, "", err
	}

	return resp.StatusCode, string(body), nil
}

// text 文本类型
type text struct {
	MsgType string `json:"msgType"`
	Content string `json:"content"`
}

func NewTextMsg(content string) string {
	msgInit := text{
		MsgType: "sampleText",
		Content: content,
	}
	msg, _ := json.Marshal(msgInit)
	return string(msg)
}

// Markdown类型
type markdown struct {
	MsgType string `json:"msgType"`
	Title   string `json:"title"`
	Text    string `json:"text"`
}

func NewMarkdownMsg(title string, text string) string {
	msgInit := markdown{
		MsgType: "sampleMarkdown",
		Title:   title,
		Text:    text,
	}
	msg, _ := json.Marshal(msgInit)
	return string(msg)
}

// 图片类型
type image struct {
	MsgType  string `json:"msgType"`
	PhotoURL string `json:"photoURL"`
}

func NewImageMsg(photoURL string) string {
	msgInit := image{
		MsgType:  "sampleImageMsg",
		PhotoURL: photoURL,
	}
	msg, _ := json.Marshal(msgInit)
	return string(msg)
}

// 链接类型
type link struct {
	MsgType    string `json:"msgType"`
	Text       string `json:"text"`
	Title      string `json:"title"`
	PicUrl     string `json:"picUrl"`
	MessageUrl string `json:"messageUrl"`
}

func NewLinkMsg(text string, title string, picUrl string, msgUrl string) string {
	msgInit := link{
		MsgType:    "sampleLink",
		Text:       text,
		Title:      title,
		PicUrl:     picUrl,
		MessageUrl: msgUrl,
	}
	msg, _ := json.Marshal(msgInit)
	return string(msg)
}

//actionCard类型

type actionCard struct {
	MsgType     string `json:"msgType"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	SingleTitle string `json:"singleTitle"`
	SingleURL   string `json:"singleURL"`
}

func NewActionCardMsg(text string, title string, singleTitle string, singleURL string) string {
	msgInit := actionCard{
		MsgType:     "sampleActionCard",
		Text:        text,
		Title:       title,
		SingleURL:   singleURL,
		SingleTitle: singleTitle,
	}
	msg, _ := json.Marshal(msgInit)
	return string(msg)
}

type actionCard2 struct {
	MsgType      string `json:"msgType"`
	Title        string `json:"title"`
	Text         string `json:"text"`
	ActionTitle1 string `json:"actionTitle1"`
	ActionURL1   string `json:"actionURL1"`
	ActionTitle2 string `json:"actionTitle2"`
	ActionURL2   string `json:"actionURL2"`
}

func NewActionCard2Msg(text string, title string, actionTitle1 string, actionURL1 string, actionTitle2 string, actionURL2 string) string {
	msgInit := actionCard2{
		MsgType:      "sampleActionCard2",
		Text:         text,
		Title:        title,
		ActionTitle1: actionTitle1,
		ActionURL1:   actionURL1,
		ActionTitle2: actionTitle2,
		ActionURL2:   actionURL2,
	}
	msg, _ := json.Marshal(msgInit)
	return string(msg)
}

type actionCard3 struct {
	MsgType      string `json:"msgType"`
	Title        string `json:"title"`
	Text         string `json:"text"`
	ActionTitle1 string `json:"actionTitle1"`
	ActionURL1   string `json:"actionURL1"`
	ActionTitle2 string `json:"actionTitle2"`
	ActionURL2   string `json:"actionURL2"`
	ActionTitle3 string `json:"actionTitle3"`
	ActionURL3   string `json:"actionURL3"`
}

func NewActionCard3Msg(text string, title string, actionTitle1 string, actionURL1 string, actionTitle2 string, actionURL2 string, actionTitle3 string, actionURL3 string) string {
	msgInit := actionCard3{
		MsgType:      "sampleActionCard3",
		Text:         text,
		Title:        title,
		ActionTitle1: actionTitle1,
		ActionURL1:   actionURL1,
		ActionTitle2: actionTitle2,
		ActionURL2:   actionURL2,
		ActionTitle3: actionTitle3,
		ActionURL3:   actionURL3,
	}
	msg, _ := json.Marshal(msgInit)
	return string(msg)
}

type actionCard4 struct {
	MsgType      string `json:"msgType"`
	Title        string `json:"title"`
	Text         string `json:"text"`
	ActionTitle1 string `json:"actionTitle1"`
	ActionURL1   string `json:"actionURL1"`
	ActionTitle2 string `json:"actionTitle2"`
	ActionURL2   string `json:"actionURL2"`
	ActionTitle3 string `json:"actionTitle3"`
	ActionURL3   string `json:"actionURL3"`
	ActionTitle4 string `json:"actionTitle4"`
	ActionURL4   string `json:"actionURL4"`
}

func NewActionCard4Msg(text string, title string, actionTitle1 string, actionURL1 string, actionTitle2 string, actionURL2 string, actionTitle3 string, actionURL3 string, actionTitle4 string, actionURL4 string) string {
	msgInit := actionCard4{
		MsgType:      "sampleActionCard4",
		Text:         text,
		Title:        title,
		ActionTitle1: actionTitle1,
		ActionURL1:   actionURL1,
		ActionTitle2: actionTitle2,
		ActionURL2:   actionURL2,
		ActionTitle3: actionTitle3,
		ActionURL3:   actionURL3,
		ActionTitle4: actionTitle4,
		ActionURL4:   actionURL4,
	}
	msg, _ := json.Marshal(msgInit)
	return string(msg)
}

type actionCard5 struct {
	MsgType      string `json:"msgType"`
	Title        string `json:"title"`
	Text         string `json:"text"`
	ActionTitle1 string `json:"actionTitle1"`
	ActionURL1   string `json:"actionURL1"`
	ActionTitle2 string `json:"actionTitle2"`
	ActionURL2   string `json:"actionURL2"`
	ActionTitle3 string `json:"actionTitle3"`
	ActionURL3   string `json:"actionURL3"`
	ActionTitle4 string `json:"actionTitle4"`
	ActionURL4   string `json:"actionURL4"`
	ActionTitle5 string `json:"actionTitle5"`
	ActionURL5   string `json:"actionURL5"`
}

func NewActionCard5Msg(text string, title string, actionTitle1 string, actionURL1 string, actionTitle2 string, actionURL2 string, actionTitle3 string, actionURL3 string, actionTitle4 string, actionURL4 string, actionTitle5 string, actionURL5 string) string {
	msgInit := actionCard5{
		MsgType:      "sampleActionCard5",
		Text:         text,
		Title:        title,
		ActionTitle1: actionTitle1,
		ActionURL1:   actionURL1,
		ActionTitle2: actionTitle2,
		ActionURL2:   actionURL2,
		ActionTitle3: actionTitle3,
		ActionURL3:   actionURL3,
		ActionTitle4: actionTitle4,
		ActionURL4:   actionURL4,
		ActionTitle5: actionTitle5,
		ActionURL5:   actionURL5,
	}
	msg, _ := json.Marshal(msgInit)
	return string(msg)
}

type actionCard6 struct {
	MsgType      string `json:"msgType"`
	Title        string `json:"title"`
	Text         string `json:"text"`
	ButtonTitle1 string `json:"buttonTitle1"`
	ButtonUrl1   string `json:"buttonUrl1"`
	ButtonTitle2 string `json:"buttonTitle2"`
	ButtonUrl2   string `json:"buttonUrl2"`
}

func NewActionCard6Msg(text string, title string, buttonTitle1 string, buttonUrl1 string, buttonTitle2 string, buttonUrl2 string) string {
	msgInit := actionCard6{
		MsgType:      "sampleActionCard6",
		Text:         text,
		Title:        title,
		ButtonTitle1: buttonTitle1,
		ButtonUrl1:   buttonUrl1,
		ButtonTitle2: buttonTitle2,
		ButtonUrl2:   buttonUrl2,
	}
	msg, _ := json.Marshal(msgInit)
	return string(msg)
}

// 企业内机器人发送消息权限封装

// SendGroupMessages 企业机器人向内部群发消息
func (h *Handle) SendGroupMessages(msg string, conversationId string) (statusCode int, bodyData string, err error) {
	// 数据拼接
	body := make(map[string]interface{})
	body["msgParam"] = msg
	body["msgKey"] = gjson.Get(msg, "msgType").String()
	body["openConversationId"] = conversationId
	body["robotCode"] = h.AppKey
	bytesData, err := json.Marshal(body)
	if err != nil {
		return 0, "", err
	}

	reader := bytes.NewReader(bytesData)

	return h.req("https://api.dingtalk.com/v1.0/robot/groupMessages/send", reader)
}
