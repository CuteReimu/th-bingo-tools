package main

import "bytes"

// TH07 (东方妖妖梦) 内存布局
// 基于 thprac 的 thprac_th07.h：
// - GameManager 位于绝对地址 0x626270
// - card_data[141] 在 GameManager + 0x18
// - card_data2[141] 在 GameManager + 0x4230
// - difficulty 在 GameManager + 0x10
// - full_shottype 在 GameManager + 0x93D7
const (
	th07BaseOffset    = 0x226270 // 0x626270 - 0x400000
	th07CardDataOff   = th07BaseOffset + 0x18
	th07CardData2Off  = th07BaseOffset + 0x4230
	th07DifficultyOff = th07BaseOffset + 0x10
	th07ShottypeOff   = th07BaseOffset + 0x93D7
)

// th07SpellInfo 对应 thprac_th07.h 中的 ScorefileCatk 结构体（0x78 = 120 字节）
type th07SpellInfo struct {
	Header             [0x0C]byte // ScorefileChapterHeader
	ShottypeMaxBonuses [6]uint32  // +0x0C
	BestMaxBonus       uint32     // +0x24
	SpellNumber        uint16     // +0x28
	NameHash           uint8      // +0x2A
	SpellName          [0x31]byte // +0x2B
	ShottypeAttempts   [6]uint16  // +0x5C
	TotalAttempts      uint16     // +0x68
	ShottypeCaptures   [6]uint16  // +0x6A
	TotalCaptures      uint16     // +0x76
}

var th07RoleNames = []string{"ReimuA", "ReimuB", "MarisaA", "MarisaB", "SakuyaA", "SakuyaB"}

func newTH07Listener() *globalCardListener[th07SpellInfo] {
	return newGlobalCardListener[th07SpellInfo](
		globalCardConfig{
			GameID:        7,
			ExeNames:      makeExeNames("th07.exe", "th07e.exe"),
			CardDataOff:   th07CardDataOff,
			CardData2Off:  th07CardData2Off,
			DifficultyOff: th07DifficultyOff,
			ShottypeOff:   th07ShottypeOff,
			CardCount:     141,
			RoleNames:     th07RoleNames,
		},
		globalCardAccessor[th07SpellInfo]{
			getNumber:   func(s *th07SpellInfo) uint16 { return s.SpellNumber },
			getName:     func(s *th07SpellInfo) string { return formatName(bytes.TrimRight(s.SpellName[:], "\000")) },
			getAttempts: func(s *th07SpellInfo) uint16 { return s.TotalAttempts },
			getCaptures: func(s *th07SpellInfo) uint16 { return s.TotalCaptures },
		},
	)
}
