package main

import (
	"bytes"

	"golang.org/x/sys/windows"
)

// th08 (东方永夜抄) 的符卡数据结构与 th07 类似：
// - 符卡统计数据是全局的（非按角色分组），但每张符卡有各机体的分别统计
// - 有两组数组：card_data（游戏模式）和 card_data2（符卡练习）
// - 使用 uint16 计数
// - 有 12 种机体（4 组队 + 8 单人）
// GameManager 位于绝对地址 0x160f528（基于 thprac 的 DIFF_ADDR 0x160f538 - 0x10 推算）
// card_data[222] 在 GameManager + 0x18
// card_data2[222] 在 GameManager + 0x18 + 222*0xA8 = GameManager + 0xE688
// difficulty 在 GameManager + 0x10
// full_shottype 在 SHOTTYPE_ADDR 0x164d0b1

const (
	th08GameManagerOffset  = 0x120F528 // 0x160f528 - 0x400000
	th08CardDataOffset     = th08GameManagerOffset + 0x18
	th08CardData2Offset    = th08GameManagerOffset + 0xE688
	th08DifficultyOffset   = th08GameManagerOffset + 0x10
	th08FullShottypeOffset = 0x124D0B1 // 0x164d0b1 - 0x400000
)

type listenerTh08 struct {
	started      bool
	cardData     [222]th08SpellInfo
	oldCardData  [222]th08SpellInfo
	cardData2    [222]th08SpellInfo
	oldCardData2 [222]th08SpellInfo
	difficulty   uint32
	fullShottype uint8
}

var th08ExeNames = append([]string{"th08.exe", "th08e.exe"}, chinesePatchExeNames...)

func (l *listenerTh08) Loop() {
	_, _, hand, baseAddress, err := findGameProcess("th08", th08ExeNames)
	if err != nil {
		l.started = false
		return
	}
	defer windows.CloseHandle(hand)
	l.oldCardData = l.cardData
	l.oldCardData2 = l.cardData2
	_ = readMemory(&l.cardData, hand, baseAddress, th08CardDataOffset)
	_ = readMemory(&l.cardData2, hand, baseAddress, th08CardData2Offset)
	_ = readMemory(&l.difficulty, hand, baseAddress, th08DifficultyOffset)
	_ = readMemory(&l.fullShottype, hand, baseAddress, th08FullShottypeOffset)
	if !l.started {
		l.started = true
		return
	}
	l.started = true
	var message *Message
	roleName := th08FormatRole(l.fullShottype)
	// 检测游戏模式符卡变化（card_data）
	for i, info := range l.cardData {
		oldInfo := l.oldCardData[i]
		msg := &Message{
			Game: 8,
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
			Game: 8,
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

func th08FormatRole(fullShottype uint8) string {
	switch fullShottype {
	case 0:
		return "ReimuYukari" // 结界组
	case 1:
		return "MarisaAlice" // 咏唱组
	case 2:
		return "SakuyaRemilia" // 红魔组
	case 3:
		return "YoumuYuyuko" // 幽冥组
	case 4:
		return "Reimu"
	case 5:
		return "Yukari"
	case 6:
		return "Marisa"
	case 7:
		return "Alice"
	case 8:
		return "Sakuya"
	case 9:
		return "Remilia"
	case 10:
		return "Youmu"
	case 11:
		return "Yuyuko"
	default:
		return "Unknown"
	}
}

// th08SpellInfo 对应 th08 的 ScorefileCatk 结构体
// 与 th07 类似但有 12 种机体（0xA8 = 168 字节）
type th08SpellInfo struct {
	header             [0x0C]byte // ScorefileChapterHeader
	shottypeMaxBonuses [12]uint32 // +0x0C
	bestMaxBonus       uint32     // +0x3C
	spellNumber        uint16     // +0x40
	nameHash           uint8      // +0x42
	spellName          [0x31]byte // +0x43
	shottypeAttempts   [12]uint16 // +0x74
	totalAttempts      uint16     // +0x8C
	shottypeCaptures   [12]uint16 // +0x8E
	totalCaptures      uint16     // +0xA6
}
