package main

// TH15 (东方绀珠传)
// 4 个角色，每角色 107 张符卡，有符卡练习
func newTH15Listener() *roleSpellListener[modernSpellInfo] {
	return newRoleSpellListener[modernSpellInfo](15, makeExeNames("th15.exe", "th15e.exe"), roleSpellConfig{
		BasePointerOffset: 0xE9BC0,
		RoleIDOffset:      20,
		SpellsOffset:      0x8D8,
		RoleStride:        0x5318,
		RoleCount:         4,
		SpellCount:        107,
		RoleNames:         []string{"Reimu", "Marisa", "Sanae", "Reisen"},
	}, newModernAccessor[modernSpellInfo, *modernSpellInfo]())
}
