package main

// TH17 (东方鬼形兽)
// 9 个角色（3 角色 × 3 灵兽），每角色 101 张符卡，有符卡练习
func newTH17Listener() *roleSpellListener[modernSpellInfo] {
	return newRoleSpellListener[modernSpellInfo](17, makeExeNames("th17.exe", "th17e.exe"), roleSpellConfig{
		BasePointerOffset: 0xB77DC,
		RoleIDOffset:      20,
		SpellsOffset:      0x8D8,
		RoleStride:        0x4820,
		RoleCount:         9,
		SpellCount:        101,
		RoleNames: []string{
			"ReimuW", "ReimuO", "ReimuE",
			"MarisaW", "MarisaO", "MarisaE",
			"YoumuW", "YoumuO", "YoumuE",
		},
	}, newModernAccessor[modernSpellInfo, *modernSpellInfo]())
}
