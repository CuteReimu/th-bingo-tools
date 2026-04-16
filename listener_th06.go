package main

import (
	"bytes"

	"golang.org/x/sys/windows"
)

// th06 (东方红魔乡) 的符卡数据结构与 th10+ 不同：
// - 不是按角色分组，而是全局 64 张符卡的统一数组
// - 使用 uint16 计数
// - 没有符卡练习模式
// GameManager 位于绝对地址 0x69BCA0（基于 thprac 的 thprac_th06.h）
// catk[64] 在 GameManager + 0x30
// difficulty 在 GameManager + 0x10
// character 在 GameManager + 0x181D
// shotType 在 GameManager + 0x181E

const (
	th06GameManagerOffset = 0x29BCA0 // 0x69BCA0 - 0x400000
	th06CatkOffset        = th06GameManagerOffset + 0x30
	th06DifficultyOffset  = th06GameManagerOffset + 0x10
	th06CharacterOffset   = th06GameManagerOffset + 0x181D
	th06ShotTypeOffset    = th06GameManagerOffset + 0x181E
)

type listenerTh06 struct {
	started    bool
	spells     [64]th06SpellInfo
	oldSpells  [64]th06SpellInfo
	difficulty uint32
	character  uint8
	shotType   uint8
}

var th06ExeNames = append([]string{"th06.exe", "th06e.exe", "東方紅魔郷.exe"}, chinesePatchExeNames...)

func (l *listenerTh06) Loop() {
	_, _, hand, baseAddress, err := findGameProcess("th06", th06ExeNames)
	if err != nil {
		l.started = false
		return
	}
	defer windows.CloseHandle(hand)
	l.oldSpells = l.spells
	_ = readMemory(&l.spells, hand, baseAddress, th06CatkOffset)
	_ = readMemory(&l.difficulty, hand, baseAddress, th06DifficultyOffset)
	_ = readMemory(&l.character, hand, baseAddress, th06CharacterOffset)
	_ = readMemory(&l.shotType, hand, baseAddress, th06ShotTypeOffset)
	if !l.started {
		l.started = true
		return
	}
	l.started = true
	var message *Message
	for i, info := range l.spells {
		oldInfo := l.oldSpells[i]
		msg := &Message{
			Game: 6,
			Id:   uint32(info.idx) + 1,
			Name: formatName(bytes.TrimRight(info.name[:], "\000")),
			Role: th06FormatRole(l.character, l.shotType),
			Rank: formatRank(l.difficulty),
		}
		if info.numAttempts > oldInfo.numAttempts {
			if message != nil || info.numAttempts != oldInfo.numAttempts+1 {
				return // 同一时间只可能改变一张符卡
			}
			msg.Event = 0
			msg.Mode = 0
			message = msg
		}
		if info.numSuccess > oldInfo.numSuccess {
			if message != nil || info.numSuccess != oldInfo.numSuccess+1 {
				return // 同一时间只可能改变一张符卡
			}
			msg.Event = 1
			msg.Mode = 0
			message = msg
		}
	}
	if message != nil {
		broadcast(message)
	}
}

func th06FormatRole(character, shotType uint8) string {
	switch character*2 + shotType {
	case 0:
		return "ReimuA"
	case 1:
		return "ReimuB"
	case 2:
		return "MarisaA"
	case 3:
		return "MarisaB"
	default:
		return "Unknown"
	}
}

// th06SpellInfo 对应 thprac_th06.h 中的 Catk 结构体（0x40 = 64 字节）
type th06SpellInfo struct {
	_base        [10]byte // Th6k base
	_pad0        [2]byte  // padding for alignment
	captureScore int32
	idx          uint16
	nameCsum     uint8
	charShot     uint8
	_unk14       uint32
	name         [32]byte
	_unk38       uint32
	numAttempts  uint16
	numSuccess   uint16
}
