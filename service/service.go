package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Service struct {
	HttpSvc     *gin.Engine
	Broadcaster *Broadcaster
}

type Broadcaster struct {
	MessageChan chan string
	Onlinemap   map[string][]Client
}

type Client struct {
	Wsconn *websocket.Conn
	RoomID string
	Ip     string
}
