package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var port = flag.Int("p", 9760, "websocket listening port")

var chanMap sync.Map

func main() {
	flag.Parse()
	if !flag.Parsed() {
		flag.Usage()
		os.Exit(-1)
	}
	for i := range listeners {
		l := listeners[i]
		go func() {
			for {
				l.Loop()
				time.Sleep(100 * time.Millisecond)
			}
		}()
	}
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool { return true },
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		id := c.RemoteAddr().String()
		log.Println("connected:", id)
		ch := make(chan []byte, 64)
		chanMap.Store(id, ch)
		defer func() {
			chanMap.Delete(id)
			_ = c.Close()
		}()
		for {
			err = c.WriteMessage(websocket.TextMessage, <-ch)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	})
	if err := http.ListenAndServe("127.0.0.1:"+strconv.Itoa(*port), nil); err != nil {
		log.Println(err)
	}
}
