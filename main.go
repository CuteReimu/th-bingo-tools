package main

import (
	"bytes"
	"fmt"
	"golang.org/x/sys/windows"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	pid, err := getPidByProcessName("th18.exe")
	if err != nil {
		panic(err)
	}
	hand, err := getProcessHandle(uint32(pid))
	if err != nil {
		panic(err)
	}
	baseAddress, err := getModuleBaseAddress(hand, "th18.exe")
	if err != nil {
		panic(err)
	}
	getCount(hand, baseAddress)
	//spellPracticeGet, spellPracticeTotal, gameModeGet, gameModeTotal := getCount(hand, baseAddress)
	//fmt.Println("1面1符E难度 Game Mode: ", gameModeGet, "/", gameModeTotal)
	//fmt.Println("1面1符E难度 Spell Practice: ", spellPracticeGet, "/", spellPracticeTotal)
	//for {
	//	spellPracticeGet2, spellPracticeTotal2, gameModeGet2, gameModeTotal2 := getCount(hand, baseAddress)
	//	//fmt.Println("1面1符E难度 Game Mode: ", gameModeGet2, "/", gameModeTotal2)
	//	//fmt.Println("1面1符E难度 Spell Practice: ", spellPracticeGet2, "/", spellPracticeTotal2)
	//	if spellPracticeTotal2 > spellPracticeTotal {
	//		fmt.Println("符卡开始 Spell Practice")
	//	}
	//	if spellPracticeGet2 > spellPracticeGet {
	//		fmt.Println("符卡收取 Spell Practice")
	//	}
	//	if gameModeTotal2 > gameModeTotal {
	//		fmt.Println("符卡开始 Game Mode")
	//	}
	//	if gameModeGet2 > gameModeGet {
	//		fmt.Println("符卡收取 Game Mode")
	//	}
	//	spellPracticeGet, spellPracticeTotal, gameModeGet, gameModeTotal = spellPracticeGet2, spellPracticeTotal2, gameModeGet2, gameModeTotal2
	//	time.Sleep(time.Second / 2)
	//}
	fmt.Scanln()
}

func getCount(hand windows.Handle, baseAddress uintptr) {
	_, allInfo := readMemory[allInfo](hand, baseAddress, 0xCF41C, 0)
	for _, role := range allInfo.roles {
		for _, info := range role.spells {
			fmt.Println("Spell", info.id+1, "role:", formatRoleId(role.id),
				"name:", formatName(bytes.TrimRight(info.name[:], "\000")), "rank:", formatRank(info.rank),
				"scprac:", info.spellPracticeGet, "/", info.spellPracticeTotal,
				"game:", info.gameModeGet, "/", info.gameModeTotal,
				"score:", uint64(info.score)*10)
		}
	}
}

type allInfo struct {
	_     [8]byte
	roles [5]roleInfo
}

type roleInfo struct {
	_         uint32
	totalTime uint16
	_         uint16
	_         uint32
	id        uint32
	_         [0x8C0]byte
	spells    [97]spellInfo
	_         [0xD4C4]byte
}

type spellInfo struct {
	name               [0xC0]byte
	gameModeGet        uint32
	spellPracticeGet   uint32
	gameModeTotal      uint32
	spellPracticeTotal uint32
	id                 uint32
	rank               uint32
	score              uint32 // 这个值乘以10才是分数
}

func formatRoleId(id uint32) string {
	switch id {
	case 0:
		return "灵梦"
	case 1:
		return "魔理沙"
	case 2:
		return "咲夜"
	case 3:
		return "早苗"
	default:
		return "未知"
	}
}

func formatRank(rank uint32) string {
	switch rank {
	case 0:
		return "E"
	case 1:
		return "N"
	case 2:
		return "H"
	case 3:
		return "L"
	case 4:
		return "EX"
	}
	return "NaN"
}

func formatName(s []byte) string {
	s2, _, _ := transform.String(japanese.ShiftJIS.NewDecoder(), string(s))
	return s2
}
