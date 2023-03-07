package main

import (
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
	&listenerTh17{},
	&listenerTh18{},
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
