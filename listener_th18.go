package main

// TH18 (东方虹龙洞)
// 4 个角色，每角色 97 张符卡，有符卡练习
// 注意：TH18 的符卡名字段为 0xC0 字节（比其他作品的 0x80 更长）
func newTH18Listener() *roleSpellListener[modernSpellInfoWide] {
	return newRoleSpellListener[modernSpellInfoWide](18, makeExeNames("th18.exe", "th18e.exe"), roleSpellConfig{
		BasePointerOffset: 0xCF41C,
		RoleIDOffset:      20,
		SpellsOffset:      0x8D8,
		RoleStride:        0x130F0,
		RoleCount:         4,
		SpellCount:        97,
		RoleNames:         []string{"Reimu", "Marisa", "Sakuya", "Sanae"},
	}, newModernAccessor[modernSpellInfoWide, *modernSpellInfoWide]())
}
