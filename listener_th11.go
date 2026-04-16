package main

// TH11 (东方地灵殿)
// 6 个角色，每角色 175 张符卡，无符卡练习
func newTH11Listener() *roleSpellListener[simpleSpellInfo] {
	return newRoleSpellListener[simpleSpellInfo](11, makeExeNames("th11.exe", "th11e.exe"), roleSpellConfig{
		BasePointerOffset: 0xA8EBC,
		RoleIDOffset:      20,
		SpellsOffset:      0x66C,
		RoleStride:        0x68D4,
		RoleCount:         6,
		SpellCount:        175,
		RoleNames:         []string{"ReimuA", "ReimuB", "ReimuC", "MarisaA", "MarisaB", "MarisaC"},
	}, simpleAccessor)
}
