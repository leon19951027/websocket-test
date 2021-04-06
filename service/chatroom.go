package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func (s Service) Chat(c *gin.Context) {
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
	fmt.Println("-----------------------")
	msg, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	Cli := &Client{
		Wsconn: msg,
		RoomID: roomid,
		Ip:     ip,
	}
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(Cli)
	defer msg.Close()
	go s.HandleConn(*Cli, isDone)
	<-isDone
	return
}

func (s Service) HandleConn(cli Client, isDone chan bool) {

	s.Broadcaster.Onlinemap[cli.RoomID] = append(s.Broadcaster.Onlinemap[cli.RoomID], cli)

	defer func() {
		//捕获read抛出的panic
		if err := recover(); err != nil {
			log.Println("read发生错误", err)
		}
	}()

	for {

		_, message, err := cli.Wsconn.ReadMessage()
		log.Println(string(message))

		go Push(s.Broadcaster.Onlinemap[cli.RoomID], message)

		if err != nil {
			log.Println("离线", cli)
			isDone <- true
			return

		}

	}

}

func Push(roomClients []Client, message []byte) {
	log.Println(roomClients)
	for _, v := range roomClients {
		go v.Wsconn.WriteMessage(1, message)
	}
	return
}
