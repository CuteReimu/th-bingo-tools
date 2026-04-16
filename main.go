package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var port = flag.Int("p", 9760, "websocket listening port")

var chanMap sync.Map

func main() {
	flag.Parse()

	fmt.Println("本程序是为东方Project沙包聚集地的bingo赛设计的自动选卡、收卡小工具。")
	fmt.Println("当你开始一张符卡时，会在bingo赛中自动选择该符卡。")
	fmt.Println("当你收取一张符卡时，会在bingo赛中自动收取该符卡。")
	fmt.Println("由于本程序不可避免的有不少bug，请视情况是否需要使用。若引发了bug导致在比赛中进行了错误的选卡、收卡操作被裁判判罚，由选手本人负责。")
	fmt.Println("目前是测试版本，支持th06-th08、th10-th18，支持日文版和中文版，请在比赛前自行测试是否能用")
	fmt.Println("本程序并不支持需要全避的符卡，请自行在bingo赛中勾选收取。")
	fmt.Println("目前，若选手把需要全避的符卡收取了，本程序会在bingo赛中自动勾选收取。请勿在比赛中进行这种操作，以免被裁判判罚。")
	fmt.Println()
	fmt.Println("若想要退出本程序，请在本窗口中按Ctrl+C")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// 启动所有游戏监听器
	for _, l := range listeners {
		go runListener(ctx, l)
	}

	// 配置 WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool { return true },
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("WebSocket 升级失败", "error", err)
			return
		}
		handleWebSocket(conn)
	})

	server := &http.Server{
		Addr:    "127.0.0.1:" + strconv.Itoa(*port),
		Handler: mux,
	}

	// 优雅关闭
	go func() {
		<-ctx.Done()
		slog.Info("正在关闭服务器...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	slog.Info("WebSocket 服务启动", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("服务器启动失败", "error", err)
	}
}

// runListener 在上下文取消前持续轮询游戏状态
func runListener(ctx context.Context, l Listener) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("监听器异常恢复", "error", r)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.Loop()
		}
	}
}

// handleWebSocket 管理单个 WebSocket 连接的生命周期
func handleWebSocket(conn *websocket.Conn) {
	id := conn.RemoteAddr()
	slog.Info("客户端连接", "remote_addr", id)

	ch := make(chan []byte, 64)
	chanMap.Store(id, ch)
	defer func() {
		chanMap.Delete(id)
		_ = conn.Close()
		slog.Info("客户端断开", "remote_addr", id)
	}()

	// 启动读取协程以检测客户端断开
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	for {
		select {
		case msg := <-ch:
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				slog.Error("WebSocket 写入失败", "error", err)
				return
			}
		case <-done:
			return
		}
	}
}
