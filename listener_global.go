package main

import (
	"golang.org/x/sys/windows"
)

// globalCardAccessor 提供对 TH07/TH08 符卡结构的通用访问
type globalCardAccessor[S any] struct {
	getNumber   func(*S) uint16
	getName     func(*S) string
	getAttempts func(*S) uint16
	getCaptures func(*S) uint16
}

// globalCardConfig 描述全局符卡数组的内存布局
type globalCardConfig struct {
	GameID        uint32
	ExeNames      []string
	CardDataOff   uintptr
	CardData2Off  uintptr
	DifficultyOff uintptr
	ShottypeOff   uintptr
	CardCount     int
	RoleNames     []string
}

// globalCardListener 是 TH07/TH08 共用的全局符卡数组监听器
type globalCardListener[S any] struct {
	config   globalCardConfig
	accessor globalCardAccessor[S]

	started      bool
	cards        []S
	oldCards     []S
	cards2       []S
	oldCards2    []S
	difficulty   uint32
	fullShottype uint8
}

func newGlobalCardListener[S any](cfg globalCardConfig, acc globalCardAccessor[S]) *globalCardListener[S] {
	return &globalCardListener[S]{
		config:    cfg,
		accessor:  acc,
		cards:     make([]S, cfg.CardCount),
		oldCards:  make([]S, cfg.CardCount),
		cards2:    make([]S, cfg.CardCount),
		oldCards2: make([]S, cfg.CardCount),
	}
}

func (l *globalCardListener[S]) Loop() {
	gameTag := "th0" + string(rune('0'+l.config.GameID))
	result, err := findGameProcess(gameTag, l.config.ExeNames)
	if err != nil {
		l.started = false
		return
	}
	defer windows.CloseHandle(result.Handle)

	copy(l.oldCards, l.cards)
	copy(l.oldCards2, l.cards2)
	readMemory(&l.cards[0], result.Handle, result.BaseAddress, l.config.CardDataOff)
	readMemory(&l.cards2[0], result.Handle, result.BaseAddress, l.config.CardData2Off)
	readMemory(&l.difficulty, result.Handle, result.BaseAddress, l.config.DifficultyOff)
	readMemory(&l.fullShottype, result.Handle, result.BaseAddress, l.config.ShottypeOff)

	if !l.started {
		l.started = true
		return
	}

	roleName := l.formatRole()
	var message *Message

	var aborted bool
	message, aborted = l.detectChanges(l.cards, l.oldCards, roleName, ModeGame, message)
	if !aborted {
		message, aborted = l.detectChanges(l.cards2, l.oldCards2, roleName, ModeSpellPractice, message)
	}

	if !aborted && message != nil {
		broadcast(message)
	}
}

func (l *globalCardListener[S]) detectChanges(cards, oldCards []S, roleName string, mode uint8, message *Message) (*Message, bool) {
	acc := &l.accessor
	for i := range cards {
		cur := &cards[i]
		old := &oldCards[i]
		msg := &Message{
			Game: l.config.GameID,
			ID:   uint32(acc.getNumber(cur)) + 1,
			Name: acc.getName(cur),
			Role: roleName,
			Rank: formatRank(l.difficulty),
		}
		if acc.getAttempts(cur) > acc.getAttempts(old) {
			if message != nil || acc.getAttempts(cur) != acc.getAttempts(old)+1 {
				return nil, true // 同一时间多张符卡变化，中止
			}
			msg.Event = EventAttempt
			msg.Mode = mode
			message = msg
		}
		if acc.getCaptures(cur) > acc.getCaptures(old) {
			if message != nil || acc.getCaptures(cur) != acc.getCaptures(old)+1 {
				return nil, true
			}
			msg.Event = EventCapture
			msg.Mode = mode
			message = msg
		}
	}
	return message, false
}

func (l *globalCardListener[S]) formatRole() string {
	idx := int(l.fullShottype)
	if idx >= 0 && idx < len(l.config.RoleNames) {
		return l.config.RoleNames[idx]
	}
	return "Unknown"
}
