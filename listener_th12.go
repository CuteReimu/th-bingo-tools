package main

// TH12 (东方星莲船)
// 6 个角色，每角色 113 张符卡，无符卡练习
func newTH12Listener() *roleSpellListener[simpleSpellInfo] {
	return newRoleSpellListener[simpleSpellInfo](12, makeExeNames("th12.exe", "th12e.exe"), roleSpellConfig{
		BasePointerOffset: 0xB451C,
		RoleIDOffset:      20,
		SpellsOffset:      0x66C,
		RoleStride:        0x45F4,
		RoleCount:         6,
		SpellCount:        113,
		RoleNames:         []string{"ReimuA", "ReimuB", "MarisaA", "MarisaB", "SanaeA", "SanaeB"},
	}, simpleAccessor)
}
