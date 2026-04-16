package main

import (
	"encoding/json"
	"fmt"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Message struct {
	Game  uint32 `json:"game"`
	Id    uint32 `json:"id"`
	Name  string `json:"name,omitempty"`
	Event uint8  `json:"event"`
	Mode  uint8  `json:"mode"`
	Role  string `json:"role"`
	Rank  string `json:"rank"`
	Score uint64 `json:"score,omitempty"`
}

type listener interface {
	Loop()
}

var listeners = []listener{
	&listenerTh06{},
	&listenerTh07{},
	&listenerTh08{},
	&listenerTh10{},
	&listenerTh11{},
	&listenerTh12{},
	&listenerTh13{},
	&listenerTh14{},
	&listenerTh15{},
	&listenerTh16{},
	&listenerTh17{},
	&listenerTh18{},
}

// broadcast 将消息序列化并广播给所有 WebSocket 订阅者
func broadcast(message *Message) {
	buf, _ := json.Marshal(message)
	fmt.Println(string(buf))
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
	}
	return "NaN"
}

func formatName(s []byte) string {
	s2, _, _ := transform.String(japanese.ShiftJIS.NewDecoder(), string(s))
	return s2
}

// chinesePatchExeNames 是中文补丁常用的可执行文件名
var chinesePatchExeNames = []string{"custom.exe", "custom_cn.exe", "custom_chs.exe", "custom_cht.exe", "custom_c.exe"}
