package BeaverBot

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func RunBeaverBot(port string, e *Event) {
	g := gin.New()
	var init Event
	g.POST("/", func(c *gin.Context) {
		dataReader := c.Request.Body
		bodyData, err := ioutil.ReadAll(dataReader)
		if err != nil {
			fmt.Println(err)
		}
		//e.Parse(bodyData)
		*e = init
		err = json.Unmarshal(bodyData, &e)
		if err != nil {
			fmt.Println("command parsing error,please feedback to the developer.error:", err)
		}
		fmt.Println(e)
	})
	err := g.Run(port)
	if err != nil {
		fmt.Printf("gin:%v", err)
	}
}
