package BeaverBot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"strconv"
	"sync"
)

func RunBeaverBot(port string, e *Event) {
	// bot判断缓存池(判断任务是那个机器人来处理的缓存池,4是数量)
	var wg sync.WaitGroup
	var id int
	pool := make(chan whoSign, 4)
	for i := 0; i < cap(pool); i++ {
		go worker(pool, &wg, &id)
	}
	// gin
	g := gin.New()
	var init Event
	g.POST("/", func(c *gin.Context) {
		*e = init
		// body处理
		dataBody := c.Request.Body
		bodyData, err := ioutil.ReadAll(dataBody)
		defer func(dataReader io.ReadCloser) {
			err := dataReader.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(dataBody)
		if err != nil {
			fmt.Println(err)
		}
		err = json.Unmarshal(bodyData, &e)
		if err != nil {
			fmt.Println("command parsing error,please feedback to the developer.error:", err)
		}

		// header处理
		timestamp, err := strconv.ParseInt(c.Request.Header.Get("timestamp"), 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		e.header.timestamp = timestamp
		e.header.sign = c.Request.Header.Get("sign")
		defer c.Request.Header.Clone()

		// 执行事件处理器
		e.explain(&pool, &wg, &id)

	})
	err := g.Run(port)
	if err != nil {
		fmt.Printf("gin:%v", err)
	}
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
