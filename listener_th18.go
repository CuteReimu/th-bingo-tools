package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type ListenerTh18 struct {
	started   bool
	roleInfos [4]roleInfo
}

func (l *ListenerTh18) Loop() {
	pid, err := getPidByProcessName("th18.exe")
	if err != nil {
		l.started = false
		return
	}
	hand, err := getProcessHandle(uint32(pid))
	if err != nil {
		l.started = false
		return
	}
	baseAddress, err := getModuleBaseAddress(hand, "th18.exe")
	if err != nil {
		l.started = false
		return
	}
	oldInfos := l.roleInfos
	l.roleInfos = [4]roleInfo{}
	for i := range l.roleInfos {
		_, l.roleInfos[i].id = readMemory[uint32](hand, baseAddress, 0xCF41C, 20+0x130F0*uintptr(i))
		_, l.roleInfos[i].spells = readMemory[[97]spellInfo](hand, baseAddress, 0xCF41C, 0x8D8+0x130F0*uintptr(i))
	}
	if !l.started {
		l.started = true
		return
	}
	l.started = true
	for i, role := range l.roleInfos {
		roleName := formatRoleId(role.id)
		for j, info := range role.spells {
			oldInfo := oldInfos[i].spells[j]
			spellPracticeGet, spellPracticeTotal, gameModeGet, gameModeTotal := oldInfo.spellPracticeGet, oldInfo.spellPracticeTotal, oldInfo.gameModeGet, oldInfo.gameModeTotal
			spellPracticeGet2, spellPracticeTotal2, gameModeGet2, gameModeTotal2 := info.spellPracticeGet, info.spellPracticeTotal, info.gameModeGet, info.gameModeTotal
			if spellPracticeTotal2 > spellPracticeTotal {
				buf, _ := json.Marshal(&Message{
					Game:  18,
					Id:    info.id + 1,
					Name:  formatName(bytes.TrimRight(info.name[:], "\000")),
					Event: 0,
					Mode:  1,
					Role:  roleName,
					Rank:  formatRank(info.rank),
				})
				fmt.Println(string(buf))
				chanMap.Range(func(_, value any) bool {
					value.(chan []byte) <- buf
					return true
				})
			}
			if spellPracticeGet2 > spellPracticeGet {
				buf, _ := json.Marshal(&Message{
					Game:  18,
					Id:    info.id + 1,
					Name:  formatName(bytes.TrimRight(info.name[:], "\000")),
					Event: 1,
					Mode:  1,
					Role:  roleName,
					Rank:  formatRank(info.rank),
				})
				fmt.Println(string(buf))
				chanMap.Range(func(_, value any) bool {
					value.(chan []byte) <- buf
					return true
				})
			}
			if gameModeTotal2 > gameModeTotal {
				buf, _ := json.Marshal(&Message{
					Game:  18,
					Id:    info.id + 1,
					Name:  formatName(bytes.TrimRight(info.name[:], "\000")),
					Event: 0,
					Mode:  0,
					Role:  roleName,
					Rank:  formatRank(info.rank),
				})
				fmt.Println(string(buf))
				chanMap.Range(func(_, value any) bool {
					value.(chan []byte) <- buf
					return true
				})
			}
			if gameModeGet2 > gameModeGet {
				buf, _ := json.Marshal(&Message{
					Game:  18,
					Id:    info.id + 1,
					Name:  formatName(bytes.TrimRight(info.name[:], "\000")),
					Event: 1,
					Mode:  0,
					Role:  roleName,
					Rank:  formatRank(info.rank),
				})
				fmt.Println(string(buf))
				chanMap.Range(func(_, value any) bool {
					value.(chan []byte) <- buf
					return true
				})
			}
		}
	}
}

type roleInfo struct {
	id     uint32
	spells [97]spellInfo
}

type spellInfo struct {
	name               [0xC0]byte
	gameModeGet        uint32
	spellPracticeGet   uint32
	gameModeTotal      uint32
	spellPracticeTotal uint32
	id                 uint32
	rank               uint32
	score              uint32 // 这个值乘以10才是分数
}

func formatRoleId(id uint32) string {
	switch id {
	case 0:
		return "Reimu"
	case 1:
		return "Marisa"
	case 2:
		return "Sakuya"
	case 3:
		return "Sanae"
	default:
		return "Unknown"
	}
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
