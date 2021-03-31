package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "192.168.1.180:8080", "http service address")

func getKeyborad(message chan string) {
	for {
		keyborad := bufio.NewReader(os.Stdin)
		s1, _ := keyborad.ReadString('\n')
		message <- s1
	}

}

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws", ForceQuery: true, RawQuery: "roomid=1"}

	log.Printf("connecting to %s", u.String())

	header := make(http.Header)
	fmt.Println("--------------------")
	header.Add("AK", "r7HdGvhpC$dUn3Q")

	log.Println(header)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("error:", err)
	}

	defer c.Close()

	message := make(chan string)
	done := make(chan struct{})
	go getKeyborad(message)

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg := <-message:
			err := c.WriteMessage(2, []byte(msg))
			if err != nil {
				time.Sleep(1 * time.Second)
				log.Println("write:", err)
				//	return
			}
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}
}
