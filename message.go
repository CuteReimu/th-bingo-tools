package main

// Event 表示符卡事件类型
const (
	EventAttempt uint8 = 0 // 挑战符卡
	EventCapture uint8 = 1 // 收取符卡
)

// Mode 表示游戏模式
const (
	ModeGame          uint8 = 0 // 游戏模式
	ModeSpellPractice uint8 = 1 // 符卡练习模式
)

// Message 是通过 WebSocket 广播的符卡变化消息
type Message struct {
	Game  uint32 `json:"game"`
	ID    uint32 `json:"id"`
	Name  string `json:"name,omitempty"`
	Event uint8  `json:"event"`
	Mode  uint8  `json:"mode"`
	Role  string `json:"role"`
	Rank  string `json:"rank"`
	Score uint64 `json:"score,omitempty"`
}
