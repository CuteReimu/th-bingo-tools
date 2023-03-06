package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type listenerTh17 struct {
	started      bool
	roleInfos    [9]th17RoleInfo
	oldRoleInfos [9]th17RoleInfo
}

func (l *listenerTh17) Loop() {
	pid, err := getPidByProcessName("th17.exe")
	if err != nil {
		l.started = false
		return
	}
	hand, err := getProcessHandle(uint32(pid))
	if err != nil {
		l.started = false
		return
	}
	baseAddress, err := getModuleBaseAddress(hand, "th17.exe")
	if err != nil {
		l.started = false
		return
	}
	l.oldRoleInfos = l.roleInfos
	for i := range l.roleInfos {
		_, l.roleInfos[i].id = readMemory[uint32](hand, baseAddress, 0xB77DC, 20+0x4820*uintptr(i))
		_, l.roleInfos[i].spells = readMemory[[101]th17SpellInfo](hand, baseAddress, 0xB77DC, 0x8D8+0x4820*uintptr(i)) // i = 4  0x12958
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
				Game: 17,
				Id:   info.id + 1,
				Name: formatName(bytes.TrimRight(info.name[:], "\000")),
				Role: roleName,
				Rank: formatRank(info.rank),
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

type th17RoleInfo struct {
	id     uint32
	spells [101]th17SpellInfo
}

func (info *th17RoleInfo) formatRoleId() string {
	switch info.id {
	case 0:
		return "ReimuW"
	case 1:
		return "ReimuO"
	case 2:
		return "ReimuE"
	case 3:
		return "MarisaW"
	case 4:
		return "MarisaO"
	case 5:
		return "MarisaE"
	case 6:
		return "YoumuW"
	case 7:
		return "YoumuO"
	case 8:
		return "YoumuE"
	default:
		return "Unknown"
	}
}

type th17SpellInfo struct {
	name               [0x80]byte
	gameModeGet        uint32
	spellPracticeGet   uint32
	gameModeTotal      uint32
	spellPracticeTotal uint32
	id                 uint32
	rank               uint32
	score              uint32 // 这个值乘以10才是分数
}
