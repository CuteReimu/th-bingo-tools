package main

type Message struct {
	Game  uint32 `json:"game"`
	Id    uint32 `json:"id"`
	Name  string `json:"name,omitempty"`
	Event uint8  `json:"event"`
	Mode  uint8  `json:"mode"`
	Role  string `json:"role"`
	Rank  string `json:"rank"`
}

type listener interface {
	Loop()
}

var listeners = []listener{
	&ListenerTh18{},
}
