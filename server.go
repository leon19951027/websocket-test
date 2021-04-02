package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	wsconn *websocket.Conn
	roomID string
	ip     string
}

var CliSlice []Client
var messagechan = make(chan []byte)
var onlinemap = make(map[string][]Client)

func main() {
	g := gin.Default()
	//	g.Use(MiddleJWT())
	g.GET("/ws", Chat)
	g.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": "200"})
	})
	//	go Push()
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

func Chat(c *gin.Context) {
	isDone := make(chan bool)
	ip := c.Request.RemoteAddr
	roomid := c.Query("roomid")

	log.Println(roomid)
	//welcome := "欢迎用户:  " + ip + " roomid:" + roomid
	//	messagechan <- []byte(welcome)
	var upgrader = websocket.Upgrader{
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	msg, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	Cli := &Client{
		wsconn: msg,
		roomID: roomid,
		ip:     ip,
	}
	if err != nil {
		log.Println(err)
		return
	}
	defer msg.Close()
	go HandleConn(*Cli, isDone)
	<-isDone
	return
}

func HandleConn(cli Client, isDone chan bool) {

	fmt.Println("----------------")

	onlinemap[cli.roomID] = append(onlinemap[cli.roomID], cli)

	defer func() {
		//捕获read抛出的panic
		if err := recover(); err != nil {
			log.Println("read发生错误", err)
		}
	}()

	for {

		_, message, err := cli.wsconn.ReadMessage()
		log.Println(string(message))

		go Push(onlinemap[cli.roomID], message)

		if err != nil {
			log.Println("离线", cli)
			isDone <- true

		}

	}
	return
}

func Push(CliSlice []Client, message []byte) {

	for k, _ := range CliSlice {
		go CliSlice[k].wsconn.WriteMessage(1, message)
	}
	return
}
