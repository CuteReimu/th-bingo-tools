package main

// TH10 (东方风神录)
// 6 个角色，每角色 110 张符卡，无符卡练习
func newTH10Listener() *roleSpellListener[simpleSpellInfo] {
	return newRoleSpellListener[simpleSpellInfo](10, makeExeNames("th10.exe", "th10e.exe"), roleSpellConfig{
		BasePointerOffset: 0x7783C,
		RoleIDOffset:      20,
		SpellsOffset:      0x5A4,
		RoleStride:        0x437C,
		RoleCount:         6,
		SpellCount:        110,
		RoleNames:         []string{"ReimuA", "ReimuB", "ReimuC", "MarisaA", "MarisaB", "MarisaC"},
	}, simpleAccessor)
}
