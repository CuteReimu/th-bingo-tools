package main

import (
	"bytes"
	"unsafe"

	"golang.org/x/sys/windows"
)

// modernSpellFields 是 TH13+ 符卡数据结构的公共字段接口
type modernSpellFields interface {
	nameBytes() []byte
	getID() uint32
	getRank() uint32
	getScore() uint32
	getGameModeGet() uint32
	getGameModeTotal() uint32
	getSpellPracticeGet() uint32
	getSpellPracticeTotal() uint32
}

// modernSpellInfo 是 TH13-TH17 共用的符卡数据结构（name 字段为 0x80 字节）
type modernSpellInfo struct {
	Name               [0x80]byte
	GameModeGet        uint32
	SpellPracticeGet   uint32
	GameModeTotal      uint32
	SpellPracticeTotal uint32
	ID                 uint32
	Rank               uint32
	Score              uint32 // 实际分数 = Score * 10
}

func (s *modernSpellInfo) nameBytes() []byte             { return s.Name[:] }
func (s *modernSpellInfo) getID() uint32                 { return s.ID }
func (s *modernSpellInfo) getRank() uint32               { return s.Rank }
func (s *modernSpellInfo) getScore() uint32              { return s.Score }
func (s *modernSpellInfo) getGameModeGet() uint32        { return s.GameModeGet }
func (s *modernSpellInfo) getGameModeTotal() uint32      { return s.GameModeTotal }
func (s *modernSpellInfo) getSpellPracticeGet() uint32   { return s.SpellPracticeGet }
func (s *modernSpellInfo) getSpellPracticeTotal() uint32 { return s.SpellPracticeTotal }

// modernSpellInfoWide 是 TH18 的符卡数据结构（name 字段为 0xC0 字节）
type modernSpellInfoWide struct {
	Name               [0xC0]byte
	GameModeGet        uint32
	SpellPracticeGet   uint32
	GameModeTotal      uint32
	SpellPracticeTotal uint32
	ID                 uint32
	Rank               uint32
	Score              uint32
}

func (s *modernSpellInfoWide) nameBytes() []byte             { return s.Name[:] }
func (s *modernSpellInfoWide) getID() uint32                 { return s.ID }
func (s *modernSpellInfoWide) getRank() uint32               { return s.Rank }
func (s *modernSpellInfoWide) getScore() uint32              { return s.Score }
func (s *modernSpellInfoWide) getGameModeGet() uint32        { return s.GameModeGet }
func (s *modernSpellInfoWide) getGameModeTotal() uint32      { return s.GameModeTotal }
func (s *modernSpellInfoWide) getSpellPracticeGet() uint32   { return s.SpellPracticeGet }
func (s *modernSpellInfoWide) getSpellPracticeTotal() uint32 { return s.SpellPracticeTotal }

// simpleSpellInfo 是 TH10-TH12 的符卡数据结构（无符卡练习模式）
type simpleSpellInfo struct {
	Name          [0x80]byte
	GameModeGet   uint32
	GameModeTotal uint32
	ID            uint32
	Rank          uint32
}

// spellAccessor 提供对符卡数据的通用访问方法
type spellAccessor[S any] struct {
	getID                 func(*S) uint32
	getName               func(*S) string
	getRank               func(*S) uint32
	getScore              func(*S) uint64
	getGameModeGet        func(*S) uint32
	getGameModeTotal      func(*S) uint32
	getSpellPracticeGet   func(*S) uint32
	getSpellPracticeTotal func(*S) uint32
	hasSpellPractice      bool
}

// newModernAccessor 构造 TH13+ 符卡结构的通用访问器（消除重复代码）
func newModernAccessor[S any, P interface {
	*S
	modernSpellFields
}]() spellAccessor[S] {
	return spellAccessor[S]{
		getID:                 func(s *S) uint32 { return P(s).getID() },
		getName:               func(s *S) string { return formatName(bytes.TrimRight(P(s).nameBytes(), "\000")) },
		getRank:               func(s *S) uint32 { return P(s).getRank() },
		getScore:              func(s *S) uint64 { return uint64(P(s).getScore()) * 10 },
		getGameModeGet:        func(s *S) uint32 { return P(s).getGameModeGet() },
		getGameModeTotal:      func(s *S) uint32 { return P(s).getGameModeTotal() },
		getSpellPracticeGet:   func(s *S) uint32 { return P(s).getSpellPracticeGet() },
		getSpellPracticeTotal: func(s *S) uint32 { return P(s).getSpellPracticeTotal() },
		hasSpellPractice:      true,
	}
}

var simpleAccessor = spellAccessor[simpleSpellInfo]{
	getID:            func(s *simpleSpellInfo) uint32 { return s.ID },
	getName:          func(s *simpleSpellInfo) string { return formatName(bytes.TrimRight(s.Name[:], "\000")) },
	getRank:          func(s *simpleSpellInfo) uint32 { return s.Rank },
	getScore:         func(_ *simpleSpellInfo) uint64 { return 0 },
	getGameModeGet:   func(s *simpleSpellInfo) uint32 { return s.GameModeGet },
	getGameModeTotal: func(s *simpleSpellInfo) uint32 { return s.GameModeTotal },
	hasSpellPractice: false,
}

// roleSpellConfig 描述角色符卡数据在内存中的布局
type roleSpellConfig struct {
	BasePointerOffset uintptr  // 基础指针地址偏移（相对于模块基地址）
	RoleIDOffset      uintptr  // 角色 ID 的偏移
	SpellsOffset      uintptr  // 符卡数组的偏移
	RoleStride        uintptr  // 每个角色的数据步长
	RoleCount         int      // 角色数量
	SpellCount        int      // 每个角色的符卡数量
	RoleNames         []string // 角色 ID 到名称的映射
}

// roleSpellListener 是 TH10+ 基于角色的通用监听器
type roleSpellListener[S any] struct {
	gameID   uint32
	exeNames []string
	config   roleSpellConfig
	accessor spellAccessor[S]
	started  bool
	roles    []roleData[S]
	oldRoles []roleData[S]
}

type roleData[S any] struct {
	id     uint32
	spells []S
}

func newRoleSpellListener[S any](gameID uint32, exeNames []string, cfg roleSpellConfig, acc spellAccessor[S]) *roleSpellListener[S] {
	l := &roleSpellListener[S]{
		gameID:   gameID,
		exeNames: exeNames,
		config:   cfg,
		accessor: acc,
		roles:    make([]roleData[S], cfg.RoleCount),
		oldRoles: make([]roleData[S], cfg.RoleCount),
	}
	for i := range l.roles {
		l.roles[i].spells = make([]S, cfg.SpellCount)
		l.oldRoles[i].spells = make([]S, cfg.SpellCount)
	}
	return l
}

func (l *roleSpellListener[S]) Loop() {
	gameTag := "th" + itoa(int(l.gameID))
	result, err := findGameProcess(gameTag, l.exeNames)
	if err != nil {
		l.started = false
		return
	}
	defer windows.CloseHandle(result.Handle)

	for i := range l.config.RoleCount {
		l.oldRoles[i].id = l.roles[i].id
		copy(l.oldRoles[i].spells, l.roles[i].spells)

		offset := l.config.RoleStride * uintptr(i)
		readMemory(&l.roles[i].id, result.Handle, result.BaseAddress, l.config.BasePointerOffset, l.config.RoleIDOffset+offset)
		// 读取整个符卡数组（利用连续内存布局）
		spellSize := unsafe.Sizeof(l.roles[i].spells[0])
		readMemoryRaw(result.Handle, result.BaseAddress,
			unsafe.Pointer(&l.roles[i].spells[0]),
			spellSize*uintptr(l.config.SpellCount),
			l.config.BasePointerOffset, l.config.SpellsOffset+offset)
	}

	if !l.started {
		l.started = true
		return
	}

	acc := &l.accessor
	var message *Message
	for i, role := range l.roles {
		roleName := l.formatRole(role.id)
		for j := range role.spells {
			cur := &l.roles[i].spells[j]
			old := &l.oldRoles[i].spells[j]

			msg := &Message{
				Game:  l.gameID,
				ID:    acc.getID(cur) + 1,
				Name:  acc.getName(cur),
				Role:  roleName,
				Rank:  formatRank(acc.getRank(cur)),
				Score: acc.getScore(cur),
			}

			if acc.hasSpellPractice {
				if acc.getSpellPracticeTotal(cur) > acc.getSpellPracticeTotal(old) {
					if message != nil || acc.getSpellPracticeTotal(cur) != acc.getSpellPracticeTotal(old)+1 {
						return
					}
					msg.Event = EventAttempt
					msg.Mode = ModeSpellPractice
					message = msg
				}
				if acc.getSpellPracticeGet(cur) > acc.getSpellPracticeGet(old) {
					if message != nil || acc.getSpellPracticeGet(cur) != acc.getSpellPracticeGet(old)+1 {
						return
					}
					msg.Event = EventCapture
					msg.Mode = ModeSpellPractice
					message = msg
				}
			}

			if acc.getGameModeTotal(cur) > acc.getGameModeTotal(old) {
				if message != nil || acc.getGameModeTotal(cur) != acc.getGameModeTotal(old)+1 {
					return
				}
				msg.Event = EventAttempt
				msg.Mode = ModeGame
				message = msg
			}
			if acc.getGameModeGet(cur) > acc.getGameModeGet(old) {
				if message != nil || acc.getGameModeGet(cur) != acc.getGameModeGet(old)+1 {
					return
				}
				msg.Event = EventCapture
				msg.Mode = ModeGame
				message = msg
			}
		}
	}

	if message != nil {
		broadcast(message)
	}
}

func (l *roleSpellListener[S]) formatRole(id uint32) string {
	if int(id) < len(l.config.RoleNames) {
		return l.config.RoleNames[id]
	}
	return "Unknown"
}

// readMemoryRaw 通过指针链读取指定大小的原始内存
func readMemoryRaw(handle windows.Handle, baseAddress uintptr, dst unsafe.Pointer, size uintptr, addrLinks ...uintptr) {
	addr := baseAddress
	for _, offset := range addrLinks[:len(addrLinks)-1] {
		var ptr uint32
		addr += offset
		if err := windows.ReadProcessMemory(handle, addr, (*byte)(unsafe.Pointer(&ptr)), unsafe.Sizeof(ptr), nil); err != nil {
			return
		}
		addr = uintptr(ptr)
	}
	addr += addrLinks[len(addrLinks)-1]
	_ = windows.ReadProcessMemory(handle, addr, (*byte)(dst), size, nil)
}

func itoa(i int) string {
	if i < 10 {
		return string(rune('0' + i))
	}
	return itoa(i/10) + string(rune('0'+i%10))
}
