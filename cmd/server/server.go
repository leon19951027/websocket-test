package main

import (
	"ws-test/service"

	"github.com/gin-gonic/gin"
)

func initservice() *service.Service {

	httpsvc := gin.New()
	broadcaster := &service.Broadcaster{
		MessageChan: make(chan string),
		Onlinemap:   make(map[string][]service.Client),
	}
	svc := &service.Service{
		HttpSvc:     httpsvc,
		Broadcaster: broadcaster,
	}
	return svc
}

func main() {

	svc := initservice()
	svc.HttpSvc.Use(gin.Logger(), gin.Recovery())
	svc.HttpSvc.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": "200"})
	})
	svc.HttpSvc.GET("/ws", svc.Chat)
	svc.HttpSvc.Run(":8080")
}
