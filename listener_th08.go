package main

import "bytes"

// TH08 (东方永夜抄) 内存布局
const (
	th08BaseOffset    = 0x120F528 // 0x160F528 - 0x400000
	th08CardDataOff   = th08BaseOffset + 0x18
	th08CardData2Off  = th08BaseOffset + 0xE688
	th08DifficultyOff = th08BaseOffset + 0x10
	th08ShottypeOff   = 0x124D0B1 // 0x164D0B1 - 0x400000
)

// th08SpellInfo 对应 th08 的 ScorefileCatk 结构体（0xA8 = 168 字节）
// 与 th07 类似但有 12 种机体
type th08SpellInfo struct {
	Header             [0x0C]byte // ScorefileChapterHeader
	ShottypeMaxBonuses [12]uint32 // +0x0C
	BestMaxBonus       uint32     // +0x3C
	SpellNumber        uint16     // +0x40
	NameHash           uint8      // +0x42
	SpellName          [0x31]byte // +0x43
	ShottypeAttempts   [12]uint16 // +0x74
	TotalAttempts      uint16     // +0x8C
	ShottypeCaptures   [12]uint16 // +0x8E
	TotalCaptures      uint16     // +0xA6
}

var th08RoleNames = []string{
	"ReimuYukari", "MarisaAlice", "SakuyaRemilia", "YoumuYuyuko",
	"Reimu", "Yukari", "Marisa", "Alice",
	"Sakuya", "Remilia", "Youmu", "Yuyuko",
}

func newTH08Listener() *globalCardListener[th08SpellInfo] {
	return newGlobalCardListener[th08SpellInfo](
		globalCardConfig{
			GameID:        8,
			ExeNames:      makeExeNames("th08.exe", "th08e.exe"),
			CardDataOff:   th08CardDataOff,
			CardData2Off:  th08CardData2Off,
			DifficultyOff: th08DifficultyOff,
			ShottypeOff:   th08ShottypeOff,
			CardCount:     222,
			RoleNames:     th08RoleNames,
		},
		globalCardAccessor[th08SpellInfo]{
			getNumber:   func(s *th08SpellInfo) uint16 { return s.SpellNumber },
			getName:     func(s *th08SpellInfo) string { return formatName(bytes.TrimRight(s.SpellName[:], "\000")) },
			getAttempts: func(s *th08SpellInfo) uint16 { return s.TotalAttempts },
			getCaptures: func(s *th08SpellInfo) uint16 { return s.TotalCaptures },
		},
	)
}
