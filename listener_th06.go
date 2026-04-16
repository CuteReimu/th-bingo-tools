package main

import (
	"bytes"

	"golang.org/x/sys/windows"
)

// TH06 (东方红魔乡) 内存布局
// 基于 thprac 的 thprac_th06.h：
// - GameManager 位于绝对地址 0x69BCA0
// - catk[64] 在 GameManager + 0x30（偏移 0x2C 为 isTimeStopped，前面是其他字段）
// - difficulty 在 GameManager + 0x10
// - character 在 GameManager + 0x181D
// - shotType 在 GameManager + 0x181E
const (
	th06BaseOffset      = 0x29BCA0 // 0x69BCA0 - 0x400000
	th06CatkOffset      = th06BaseOffset + 0x30
	th06DifficultyOff   = th06BaseOffset + 0x10
	th06CharacterOff    = th06BaseOffset + 0x181D
	th06ShotTypeOff     = th06BaseOffset + 0x181E
)

// th06SpellInfo 对应 thprac_th06.h 中的 Catk 结构体（0x40 = 64 字节）
type th06SpellInfo struct {
	_base        [10]byte // Th6k base
	_pad0        [2]byte  // padding for alignment
	CaptureScore int32
	Idx          uint16
	NameCsum     uint8
	CharShot     uint8
	_unk14       uint32
	Name         [32]byte
	_unk38       uint32
	NumAttempts  uint16
	NumSuccess   uint16
}

var th06RoleNames = []string{"ReimuA", "ReimuB", "MarisaA", "MarisaB"}

type listenerTH06 struct {
	started    bool
	spells     [64]th06SpellInfo
	oldSpells  [64]th06SpellInfo
	difficulty uint32
	character  uint8
	shotType   uint8
}

func newTH06Listener() *listenerTH06 {
	return &listenerTH06{}
}

func (l *listenerTH06) Loop() {
	result, err := findGameProcess("th06", makeExeNames("th06.exe", "th06e.exe", "東方紅魔郷.exe"))
	if err != nil {
		l.started = false
		return
	}
	defer windows.CloseHandle(result.Handle)

	l.oldSpells = l.spells
	readMemory(&l.spells, result.Handle, result.BaseAddress, th06CatkOffset)
	readMemory(&l.difficulty, result.Handle, result.BaseAddress, th06DifficultyOff)
	readMemory(&l.character, result.Handle, result.BaseAddress, th06CharacterOff)
	readMemory(&l.shotType, result.Handle, result.BaseAddress, th06ShotTypeOff)

	if !l.started {
		l.started = true
		return
	}

	var message *Message
	for i, info := range l.spells {
		old := l.oldSpells[i]
		msg := &Message{
			Game: 6,
			ID:   uint32(info.Idx) + 1,
			Name: formatName(bytes.TrimRight(info.Name[:], "\000")),
			Role: l.formatRole(),
			Rank: formatRank(l.difficulty),
		}
		if info.NumAttempts > old.NumAttempts {
			if message != nil || info.NumAttempts != old.NumAttempts+1 {
				return // 同一时间只可能改变一张符卡
			}
			msg.Event = EventAttempt
			msg.Mode = ModeGame
			message = msg
		}
		if info.NumSuccess > old.NumSuccess {
			if message != nil || info.NumSuccess != old.NumSuccess+1 {
				return
			}
			msg.Event = EventCapture
			msg.Mode = ModeGame
			message = msg
		}
	}
	if message != nil {
		broadcast(message)
	}
}

func (l *listenerTH06) formatRole() string {
	idx := int(l.character)*2 + int(l.shotType)
	if idx >= 0 && idx < len(th06RoleNames) {
		return th06RoleNames[idx]
	}
	return "Unknown"
}
