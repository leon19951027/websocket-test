package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	g := gin.Default()
	g.Use(MiddleJWT())
	g.GET("/ws", StartWs)
	g.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": "200"})
	})

	g.Run(":8080")
}

func MiddleJWT() gin.HandlerFunc {
	log.Println("加载中间件JWT")
	return func(c *gin.Context) {
		if c.GetHeader("AK") != "r7HdGvhpC$dUn3Q" || c.GetHeader("AK") == "" {
			log.Println("--------")
			c.JSON(400, gin.H{"code": "400", "message": "token miss"})
			c.Abort()
			return
		} else {

		}
		c.Next()
	}
}

func StartWs(c *gin.Context) {
	//log.Println(c.Request.RemoteAddr, c.Request.RequestURI)
	var upgrader = websocket.Upgrader{
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	msg, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer msg.Close()
	go read(msg)
	select {}
}

func read(c *websocket.Conn) {

	defer func() {
		//捕获read抛出的panic
		if err := recover(); err != nil {
			log.Println("read发生错误", err)
		}
	}()

	for {
		_, message, err := c.ReadMessage()
		log.Println("client message", string(message), c.RemoteAddr())
		c.WriteMessage(1, []byte("pong"))
		if err != nil { // 离线通知
			//	offline <- c
			log.Println("ReadMessage error1", err)
			break
		}

	}
}
