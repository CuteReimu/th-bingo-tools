package main

import (
	"encoding/json"
	"log/slog"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// Listener 是游戏监听器的通用接口
type Listener interface {
	Loop()
}

// chinesePatchExeNames 是中文补丁常用的可执行文件名
var chinesePatchExeNames = []string{"custom.exe", "custom_cn.exe", "custom_chs.exe", "custom_cht.exe", "custom_c.exe"}

// listeners 注册所有支持的游戏监听器
var listeners = []Listener{
	newTH06Listener(),
	newTH07Listener(),
	newTH08Listener(),
	newTH10Listener(),
	newTH11Listener(),
	newTH12Listener(),
	newTH13Listener(),
	newTH14Listener(),
	newTH15Listener(),
	newTH16Listener(),
	newTH17Listener(),
	newTH18Listener(),
}

// broadcast 将消息序列化并广播给所有 WebSocket 订阅者
func broadcast(message *Message) {
	buf, err := json.Marshal(message)
	if err != nil {
		slog.Error("序列化消息失败", "error", err)
		return
	}
	slog.Info("广播消息", "data", string(buf))
	chanMap.Range(func(_, value any) bool {
		ch := value.(chan []byte)
		select {
		case ch <- buf:
		default:
			// 通道已满，丢弃消息避免阻塞
		}
		return true
	})
}

// formatRank 将难度数值转换为字符串
func formatRank(rank uint32) string {
	switch rank {
	case 0:
		return "E"
	case 1:
		return "N"
	case 2:
		return "H"
	case 3:
		return "L"
	case 4:
		return "EX"
	default:
		return "Unknown"
	}
}

// formatName 将 Shift-JIS 编码的字节切片解码为 UTF-8 字符串
func formatName(s []byte) string {
	decoded, _, err := transform.String(japanese.ShiftJIS.NewDecoder(), string(s))
	if err != nil {
		return string(s)
	}
	return decoded
}

// makeExeNames 构建游戏的所有可能可执行文件名列表
func makeExeNames(names ...string) []string {
	return append(names, chinesePatchExeNames...)
}
