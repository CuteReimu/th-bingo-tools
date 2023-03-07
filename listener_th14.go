package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type listenerTh14 struct {
	started      bool
	roleInfos    [6]th14RoleInfo
	oldRoleInfos [6]th14RoleInfo
}

func (l *listenerTh14) Loop() {
	pid, err := getPidByProcessName("th14.exe")
	if err != nil {
		l.started = false
		return
	}
	hand, err := getProcessHandle(uint32(pid))
	if err != nil {
		l.started = false
		return
	}
	baseAddress, err := getModuleBaseAddress(hand, "th14.exe")
	if err != nil {
		l.started = false
		return
	}
	l.oldRoleInfos = l.roleInfos
	for i := range l.roleInfos {
		_ = readMemory(&l.roleInfos[i].id, hand, baseAddress, 0xDB68C, 20+0x5298*uintptr(i))
		_ = readMemory(&l.roleInfos[i].spells, hand, baseAddress, 0xDB68C, 0xAB8+0x5298*uintptr(i))
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
				Game:  14,
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

type th14RoleInfo struct {
	id     uint32
	spells [120]th14SpellInfo
}

func (info *th14RoleInfo) formatRoleId() string {
	switch info.id {
	case 0:
		return "ReimuA"
	case 1:
		return "ReimuB"
	case 2:
		return "MarisaA"
	case 3:
		return "MarisaB"
	case 4:
		return "SakuyaA"
	case 5:
		return "SakuyaB"
	default:
		return "Unknown"
	}
}

type th14SpellInfo struct {
	name               [0x80]byte
	gameModeGet        uint32
	spellPracticeGet   uint32
	gameModeTotal      uint32
	spellPracticeTotal uint32
	id                 uint32
	rank               uint32
	score              uint32 // 这个值乘以10才是分数
}
