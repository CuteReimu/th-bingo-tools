package main

import (
	"bytes"

	"golang.org/x/sys/windows"
)

// th07 (东方妖妖梦) 的符卡数据结构与 th10+ 不同：
// - 符卡统计数据是全局的（非按角色分组），但每张符卡有各机体的分别统计
// - 有两组数组：card_data（游戏模式）和 card_data2（符卡练习）
// - 使用 uint16 计数
// GameManager 位于绝对地址 0x626270（基于 thprac 的 thprac_th07.h）
// card_data[141] 在 GameManager + 0x18
// card_data2[141] 在 GameManager + 0x4230
// difficulty 在 GameManager + 0x10
// full_shottype 在 GameManager + 0x93D7

const (
	th07GameManagerOffset  = 0x226270 // 0x626270 - 0x400000
	th07CardDataOffset     = th07GameManagerOffset + 0x18
	th07CardData2Offset    = th07GameManagerOffset + 0x4230
	th07DifficultyOffset   = th07GameManagerOffset + 0x10
	th07FullShottypeOffset = th07GameManagerOffset + 0x93D7
)

type listenerTh07 struct {
	started      bool
	cardData     [141]th07SpellInfo
	oldCardData  [141]th07SpellInfo
	cardData2    [141]th07SpellInfo
	oldCardData2 [141]th07SpellInfo
	difficulty   uint32
	fullShottype uint8
}

var th07ExeNames = append([]string{"th07.exe", "th07e.exe"}, chinesePatchExeNames...)

func (l *listenerTh07) Loop() {
	_, _, hand, baseAddress, err := findGameProcess("th07", th07ExeNames)
	if err != nil {
		l.started = false
		return
	}
	defer windows.CloseHandle(hand)
	l.oldCardData = l.cardData
	l.oldCardData2 = l.cardData2
	_ = readMemory(&l.cardData, hand, baseAddress, th07CardDataOffset)
	_ = readMemory(&l.cardData2, hand, baseAddress, th07CardData2Offset)
	_ = readMemory(&l.difficulty, hand, baseAddress, th07DifficultyOffset)
	_ = readMemory(&l.fullShottype, hand, baseAddress, th07FullShottypeOffset)
	if !l.started {
		l.started = true
		return
	}
	l.started = true
	var message *Message
	roleName := th07FormatRole(l.fullShottype)
	// 检测游戏模式符卡变化（card_data）
	for i, info := range l.cardData {
		oldInfo := l.oldCardData[i]
		msg := &Message{
			Game: 7,
			Id:   uint32(info.spellNumber) + 1,
			Name: formatName(bytes.TrimRight(info.spellName[:], "\000")),
			Role: roleName,
			Rank: formatRank(l.difficulty),
		}
		if info.totalAttempts > oldInfo.totalAttempts {
			if message != nil || info.totalAttempts != oldInfo.totalAttempts+1 {
				return
			}
			msg.Event = 0
			msg.Mode = 0
			message = msg
		}
		if info.totalCaptures > oldInfo.totalCaptures {
			if message != nil || info.totalCaptures != oldInfo.totalCaptures+1 {
				return
			}
			msg.Event = 1
			msg.Mode = 0
			message = msg
		}
	}
	// 检测符卡练习模式变化（card_data2）
	for i, info := range l.cardData2 {
		oldInfo := l.oldCardData2[i]
		msg := &Message{
			Game: 7,
			Id:   uint32(info.spellNumber) + 1,
			Name: formatName(bytes.TrimRight(info.spellName[:], "\000")),
			Role: roleName,
			Rank: formatRank(l.difficulty),
		}
		if info.totalAttempts > oldInfo.totalAttempts {
			if message != nil || info.totalAttempts != oldInfo.totalAttempts+1 {
				return
			}
			msg.Event = 0
			msg.Mode = 1
			message = msg
		}
		if info.totalCaptures > oldInfo.totalCaptures {
			if message != nil || info.totalCaptures != oldInfo.totalCaptures+1 {
				return
			}
			msg.Event = 1
			msg.Mode = 1
			message = msg
		}
	}
	if message != nil {
		broadcast(message)
	}
}

func th07FormatRole(fullShottype uint8) string {
	switch fullShottype {
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

// th07SpellInfo 对应 thprac_th07.h 中的 ScorefileCatk 结构体（0x78 = 120 字节）
type th07SpellInfo struct {
	header             [0x0C]byte // ScorefileChapterHeader
	shottypeMaxBonuses [6]uint32  // +0x0C
	bestMaxBonus       uint32     // +0x24
	spellNumber        uint16     // +0x28
	nameHash           uint8      // +0x2A
	spellName          [0x31]byte // +0x2B
	shottypeAttempts   [6]uint16  // +0x5C
	totalAttempts      uint16     // +0x68
	shottypeCaptures   [6]uint16  // +0x6A
	totalCaptures      uint16     // +0x76
}
