package main

import (
	"github.com/BoyChai/BeaverBot"
)

func main() {
	h, _ := BeaverBot.NewHandle("你的appKey", "你的AppSecret")
	var e BeaverBot.Event
	c := []BeaverBot.Condition{{
		Key:   &e.Text.Content,
		Value: "hello",
		Regex: true,
	}}
	BeaverBot.NewTask(BeaverBot.Task{
		Condition: c,
		Run: func() {
			msg := BeaverBot.NewTextMsg("hello")
			h.SendGroupMessages(h.AppKey, msg, e.ConversationId)
		},
	})

	BeaverBot.RunBeaverBot(":8080", &e)
}
