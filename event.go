package BeaverBot

import (
	"errors"
	"fmt"
	"regexp"
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

// 命令处理器
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
		fmt.Println("[BeaverBot] 机器人识别失败,或签名无效")
		return
	}
	e.HandleID = *id
	*id = 0
	for i := 0; i < len(Tasks); i++ {
		task := Tasks[i]
		status := e.filterStart(task)
		if status == nil {
			return
		}
	}
}

// 触发器过滤
func (e *Event) filterStart(task Task) error {
	for t := 1; t <= len(task.Condition); t++ {
		//fmt.Println(*task.Condition[t-1].Key.(reflect.TypeOf(task.Condition[t-1].Key))
		conditionKey, _ := e.typeAsserts(task.Condition[t-1].Key)
		if t == len(task.Condition) {
			if task.Condition[t-1].Regex == true {
				key, _ := regexp.MatchString(task.Condition[t-1].Value, fmt.Sprint(conditionKey))
				if key {
					task.Run()
					return nil
				}
			}
			if fmt.Sprint(conditionKey) == task.Condition[t-1].Value {
				task.Run()
				return nil
			}
		}
		if task.Condition[t-1].Regex == true {
			key, _ := regexp.MatchString(task.Condition[t-1].Value, fmt.Sprint(conditionKey))
			if key != true {
				return errors.New("1")
			}
		}
		if fmt.Sprint(conditionKey) != task.Condition[t-1].Value {
			return errors.New("1")
		}
	}
	return errors.New("1")
}

// 类型断言
func (e *Event) typeAsserts(key interface{}) (interface{}, error) {
	switch key.(type) {
	case *int64:
		return *key.(*int64), nil
	case *string:
		return *key.(*string), nil
	case *int32:
		return *key.(*int32), nil
	default:
		return nil, errors.New("the current type is not supported. please feedback through issue")
	}
}
