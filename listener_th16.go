package main

// TH16 (东方天空璋)
// 4 个角色，每角色 119 张符卡，有符卡练习
func newTH16Listener() *roleSpellListener[modernSpellInfo] {
	return newRoleSpellListener[modernSpellInfo](16, makeExeNames("th16.exe", "th16e.exe"), roleSpellConfig{
		BasePointerOffset: 0xA6F0C,
		RoleIDOffset:      20,
		SpellsOffset:      0x8D8,
		RoleStride:        0x5318,
		RoleCount:         4,
		SpellCount:        119,
		RoleNames:         []string{"Reimu", "Cirno", "Aya", "Marisa"},
	}, newModernAccessor[modernSpellInfo, *modernSpellInfo]())
}
