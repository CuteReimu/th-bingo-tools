package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type listenerTh12 struct {
	started      bool
	roleInfos    [6]th12RoleInfo
	oldRoleInfos [6]th12RoleInfo
}

func (l *listenerTh12) Loop() {
	pid, err := getPidByProcessName("th12.exe")
	if err != nil {
		l.started = false
		return
	}
	hand, err := getProcessHandle(uint32(pid))
	if err != nil {
		l.started = false
		return
	}
	baseAddress, err := getModuleBaseAddress(hand, "th12.exe")
	if err != nil {
		l.started = false
		return
	}
	l.oldRoleInfos = l.roleInfos
	for i := range l.roleInfos {
		_ = readMemory(&l.roleInfos[i].id, hand, baseAddress, 0xB451C, 20+0x45F4*uintptr(i))
		_ = readMemory(&l.roleInfos[i].spells, hand, baseAddress, 0xB451C, 0x66C+0x45F4*uintptr(i))
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
			gameModeGet, gameModeTotal := oldInfo.gameModeGet, oldInfo.gameModeTotal
			gameModeGet2, gameModeTotal2 := info.gameModeGet, info.gameModeTotal
			msg := &Message{
				Game: 12,
				Id:   info.id + 1,
				Name: formatName(bytes.TrimRight(info.name[:], "\000")),
				Role: roleName,
				Rank: formatRank(info.rank),
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

type th12RoleInfo struct {
	id     uint32
	spells [113]th12SpellInfo
}

func (info *th12RoleInfo) formatRoleId() string {
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
		return "SanaeA"
	case 5:
		return "SanaeB"
	default:
		return "Unknown"
	}
}

type th12SpellInfo struct {
	name          [0x80]byte
	gameModeGet   uint32
	gameModeTotal uint32
	id            uint32
	rank          uint32
}
