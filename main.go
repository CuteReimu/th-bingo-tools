package main

import (
	"flag"
	"fmt"
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
	fmt.Println("本程序是为东方Project沙包聚集地的bingo赛设计的自动选卡、收卡小工具。")
	fmt.Println("当你开始一张符卡时，会在bingo赛中自动选择该符卡。")
	fmt.Println("当你收取一张符卡时，会在bingo赛中自动收取该符卡。")
	fmt.Println("由于本程序不可避免的有不少bug，请视情况是否需要使用。若引发了bug导致在比赛中进行了错误的选卡、收卡操作被裁判判罚，由选手本人负责。")
	fmt.Println("目前是测试版本，仅支持th10-th18（不含绀珠传），仅测试了日文版，尚未测试其他版本，请在比赛前自行测试是否能用")
	fmt.Println("本程序并不支持需要全避的符卡，请自行在bingo赛中勾选收取。")
	fmt.Println("目前，若选手把需要全避的符卡收取了，本程序会在bingo赛中自动勾选收取。请勿在比赛中进行这种操作，以免被裁判判罚。")
	fmt.Println()
	fmt.Println("若想要退出本程序，请在本窗口中按Ctrl+C")
	for i := range listeners {
		l := listeners[i]
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("%+v\n", r)
				}
			}()
			for {
				l.Loop()
				time.Sleep(time.Second)
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
