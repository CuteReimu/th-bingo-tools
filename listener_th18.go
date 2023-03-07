package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type listenerTh18 struct {
	started      bool
	roleInfos    [4]th18RoleInfo
	oldRoleInfos [4]th18RoleInfo
}

func (l *listenerTh18) Loop() {
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
	l.oldRoleInfos = l.roleInfos
	for i := range l.roleInfos {
		_, l.roleInfos[i].id = readMemory[uint32](hand, baseAddress, 0xCF41C, 20+0x130F0*uintptr(i))
		_, l.roleInfos[i].spells = readMemory[[97]th18SpellInfo](hand, baseAddress, 0xCF41C, 0x8D8+0x130F0*uintptr(i))
	}
	if !l.started {
		l.started = true
		return
	}
	l.started = true
	for i, role := range l.roleInfos {
		roleName := role.formatRoleId()
		for j, info := range role.spells {
			oldInfo := l.oldRoleInfos[i].spells[j]
			spellPracticeGet, spellPracticeTotal, gameModeGet, gameModeTotal := oldInfo.spellPracticeGet, oldInfo.spellPracticeTotal, oldInfo.gameModeGet, oldInfo.gameModeTotal
			spellPracticeGet2, spellPracticeTotal2, gameModeGet2, gameModeTotal2 := info.spellPracticeGet, info.spellPracticeTotal, info.gameModeGet, info.gameModeTotal
			message := &Message{
				Game:  18,
				Id:    info.id + 1,
				Name:  formatName(bytes.TrimRight(info.name[:], "\000")),
				Role:  roleName,
				Rank:  formatRank(info.rank),
				Score: uint64(info.score) * 10,
			}
			if spellPracticeTotal2 > spellPracticeTotal {
				message.Event = 0
				message.Mode = 1
				buf, _ := json.Marshal(message)
				fmt.Println(string(buf))
				chanMap.Range(func(_, value any) bool {
					value.(chan []byte) <- buf
					return true
				})
			}
			if spellPracticeGet2 > spellPracticeGet {
				message.Event = 1
				message.Mode = 1
				buf, _ := json.Marshal(message)
				fmt.Println(string(buf))
				chanMap.Range(func(_, value any) bool {
					value.(chan []byte) <- buf
					return true
				})
			}
			if gameModeTotal2 > gameModeTotal {
				message.Event = 0
				message.Mode = 0
				buf, _ := json.Marshal(message)
				fmt.Println(string(buf))
				chanMap.Range(func(_, value any) bool {
					value.(chan []byte) <- buf
					return true
				})
			}
			if gameModeGet2 > gameModeGet {
				message.Event = 1
				message.Mode = 0
				buf, _ := json.Marshal(message)
				fmt.Println(string(buf))
				chanMap.Range(func(_, value any) bool {
					value.(chan []byte) <- buf
					return true
				})
			}
		}
	}
}

type th18RoleInfo struct {
	id     uint32
	spells [97]th18SpellInfo
}

func (info *th18RoleInfo) formatRoleId() string {
	switch info.id {
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

type th18SpellInfo struct {
	name               [0xC0]byte
	gameModeGet        uint32
	spellPracticeGet   uint32
	gameModeTotal      uint32
	spellPracticeTotal uint32
	id                 uint32
	rank               uint32
	score              uint32 // 这个值乘以10才是分数
}
