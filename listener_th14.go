package main

// TH14 (东方辉针城)
// 6 个角色，每角色 120 张符卡，有符卡练习
func newTH14Listener() *roleSpellListener[modernSpellInfo] {
	return newRoleSpellListener[modernSpellInfo](14, makeExeNames("th14.exe", "th14e.exe"), roleSpellConfig{
		BasePointerOffset: 0xDB68C,
		RoleIDOffset:      20,
		SpellsOffset:      0xAB8,
		RoleStride:        0x5298,
		RoleCount:         6,
		SpellCount:        120,
		RoleNames:         []string{"ReimuA", "ReimuB", "MarisaA", "MarisaB", "SakuyaA", "SakuyaB"},
	}, newModernAccessor[modernSpellInfo, *modernSpellInfo]())
}
