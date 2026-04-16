package main

// TH13 (东方神灵庙)
// 4 个角色，每角色 119 张符卡，有符卡练习
func newTH13Listener() *roleSpellListener[modernSpellInfo] {
	return newRoleSpellListener[modernSpellInfo](13, makeExeNames("th13.exe", "th13e.exe"), roleSpellConfig{
		BasePointerOffset: 0xC22CC,
		RoleIDOffset:      20,
		SpellsOffset:      0x980,
		RoleStride:        0x56DC,
		RoleCount:         4,
		SpellCount:        119,
		RoleNames:         []string{"Reimu", "Marisa", "Sanae", "Youmu"},
	}, newModernAccessor[modernSpellInfo, *modernSpellInfo]())
}
