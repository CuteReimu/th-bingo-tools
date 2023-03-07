package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type listenerTh16 struct {
	started      bool
	roleInfos    [4]th16RoleInfo
	oldRoleInfos [4]th16RoleInfo
}

func (l *listenerTh16) Loop() {
	pid, err := getPidByProcessName("th16.exe")
	if err != nil {
		l.started = false
		return
	}
	hand, err := getProcessHandle(uint32(pid))
	if err != nil {
		l.started = false
		return
	}
	baseAddress, err := getModuleBaseAddress(hand, "th16.exe")
	if err != nil {
		l.started = false
		return
	}
	l.oldRoleInfos = l.roleInfos
	for i := range l.roleInfos {
		_ = readMemory(&l.roleInfos[i].id, hand, baseAddress, 0xA6F0C, 20+0x5318*uintptr(i))
		_ = readMemory(&l.roleInfos[i].spells, hand, baseAddress, 0xA6F0C, 0x8D8+0x5318*uintptr(i))
	}
	if !l.started {
		l.started = true
		return
	}
	l.started = true
	var message *Message
	for i, role := range l.roleInfos {
		roleName := role.formatRoleId()
		for j, info := range role.spells {
			oldInfo := l.oldRoleInfos[i].spells[j]
			spellPracticeGet, spellPracticeTotal, gameModeGet, gameModeTotal := oldInfo.spellPracticeGet, oldInfo.spellPracticeTotal, oldInfo.gameModeGet, oldInfo.gameModeTotal
			spellPracticeGet2, spellPracticeTotal2, gameModeGet2, gameModeTotal2 := info.spellPracticeGet, info.spellPracticeTotal, info.gameModeGet, info.gameModeTotal
			msg := &Message{
				Game:  16,
				Id:    info.id + 1,
				Name:  formatName(bytes.TrimRight(info.name[:], "\000")),
				Role:  roleName,
				Rank:  formatRank(info.rank),
				Score: uint64(info.score) * 10,
			}
			if spellPracticeTotal2 > spellPracticeTotal {
				if message != nil || spellPracticeTotal2 != spellPracticeTotal+1 {
					return // 同一时间只可能改变一张符卡
				}
				msg.Event = 0
				msg.Mode = 1
				message = msg
			}
			if spellPracticeGet2 > spellPracticeGet {
				if message != nil || spellPracticeGet2 != spellPracticeGet+1 {
					return // 同一时间只可能改变一张符卡
				}
				msg.Event = 1
				msg.Mode = 1
				message = msg
			}
			if gameModeTotal2 > gameModeTotal {
				if message != nil || gameModeTotal2 != gameModeTotal+1 {
					return // 同一时间只可能改变一张符卡
				}
				msg.Event = 0
				msg.Mode = 0
				message = msg
			}
			if gameModeGet2 > gameModeGet {
				if message != nil || gameModeGet2 != gameModeGet+1 {
					return // 同一时间只可能改变一张符卡
				}
				msg.Event = 1
				msg.Mode = 0
				message = msg
			}
		}
	}
	if message != nil {
		buf, _ := json.Marshal(message)
		fmt.Println(string(buf))
		chanMap.Range(func(_, value any) bool {
			value.(chan []byte) <- buf
			return true
		})
	}
}

type th16RoleInfo struct {
	id     uint32
	spells [119]th16SpellInfo
}

func (info *th16RoleInfo) formatRoleId() string {
	switch info.id {
	case 0:
		return "Reimu"
	case 1:
		return "Cirno"
	case 2:
		return "Aya"
	case 3:
		return "Marisa"
	default:
		return "Unknown"
	}
}

type th16SpellInfo struct {
	name               [0x80]byte
	gameModeGet        uint32
	spellPracticeGet   uint32
	gameModeTotal      uint32
	spellPracticeTotal uint32
	id                 uint32
	rank               uint32
	score              uint32 // 这个值乘以10才是分数
}
