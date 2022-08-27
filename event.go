package BeaverBot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

// Event 整个事件信息
type Event struct {
	header struct {
		timestamp int64
		sign      string
	}
	HandleID       int
	ConversationId string `json:"conversationId"`
	AtUsers        []struct {
		DingtalkId string `json:"dingtalkId"`
		StaffId    string `json:"staffId"`
	} `json:"atUsers"`
	ChatbotCorpId             string `json:"chatbotCorpId"`
	ChatbotUserId             string `json:"chatbotUserId"`
	MsgId                     string `json:"msgId"`
	SenderNick                string `json:"senderNick"`
	IsAdmin                   bool   `json:"isAdmin"`
	SenderStaffId             string `json:"senderStaffId"`
	SessionWebhookExpiredTime int64  `json:"sessionWebhookExpiredTime"`
	CreateAt                  int64  `json:"createAt"`
	SenderCorpId              string `json:"senderCorpId"`
	ConversationType          string `json:"conversationType"`
	SenderId                  string `json:"senderId"`
	ConversationTitle         string `json:"conversationTitle"`
	IsInAtList                bool   `json:"isInAtList"`
	SessionWebhook            string `json:"sessionWebhook"`
	Text                      struct {
		Content string `json:"content"`
	} `json:"text"`
	Msgtype string `json:"msgtype"`
}

// 消息签名信息
type whoSign struct {
	timestamp int64
	sign      string
}

func worker(pool chan whoSign, wg *sync.WaitGroup, id *int) {
	for signInfo := range pool {
		for i, handle := range AllHandle {
			secStr := fmt.Sprintf("%d\n%s", signInfo.timestamp, handle.AppSecret)
			hmac256 := hmac.New(sha256.New, []byte(handle.AppSecret))
			hmac256.Write([]byte(secStr))
			result := hmac256.Sum(nil)
			sign := base64.StdEncoding.EncodeToString(result)
			if sign == signInfo.sign {
				*id = i + 1
			}
		}
		wg.Done()
	}

}

//func (e *Event)sign(t int64, secret string) string {
//	secStr := fmt.Sprintf("%d\n%s", t, secret)
//	hmac256 := hmac.New(sha256.New, []byte(secret))
//	hmac256.Write([]byte(secStr))
//	result := hmac256.Sum(nil)
//	return base64.StdEncoding.EncodeToString(result)
//}
func (e *Event) explain(pool *chan whoSign, wg *sync.WaitGroup, id *int) {
	// 获取当前时间戳(毫秒)
	now := time.Now()
	nowTime := now.UnixNano() / 1e6
	// 时间判断是否合法
	if (e.header.timestamp-nowTime)/3600000 >= 1 {
		return
	}
	wg.Add(1)

	*pool <- whoSign{
		timestamp: e.header.timestamp,
		sign:      e.header.sign,
	}
	wg.Wait()
	if *id == 0 {
		fmt.Println("判断失败")
	}
	e.HandleID = *id
	*id = 0
	fmt.Println("签名有效,机器人ID为:", e.HandleID)
}
