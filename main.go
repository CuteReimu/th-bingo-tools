package main

import (
	"fmt"
	"golang.org/x/sys/windows"
	"time"
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
	spellPracticeGet, spellPracticeTotal, gameModeGet, gameModeTotal := getCount(hand, baseAddress)
	fmt.Println("1面1符E难度 Game Mode: ", gameModeGet, "/", gameModeTotal)
	fmt.Println("1面1符E难度 Spell Practice: ", spellPracticeGet, "/", spellPracticeTotal)
	for {
		spellPracticeGet2, spellPracticeTotal2, gameModeGet2, gameModeTotal2 := getCount(hand, baseAddress)
		//fmt.Println("1面1符E难度 Game Mode: ", gameModeGet2, "/", gameModeTotal2)
		//fmt.Println("1面1符E难度 Spell Practice: ", spellPracticeGet2, "/", spellPracticeTotal2)
		if spellPracticeTotal2 > spellPracticeTotal {
			fmt.Println("符卡开始 Spell Practice")
		}
		if spellPracticeGet2 > spellPracticeGet {
			fmt.Println("符卡收取 Spell Practice")
		}
		if gameModeTotal2 > gameModeTotal {
			fmt.Println("符卡开始 Game Mode")
		}
		if gameModeGet2 > gameModeGet {
			fmt.Println("符卡收取 Game Mode")
		}
		spellPracticeGet, spellPracticeTotal, gameModeGet, gameModeTotal = spellPracticeGet2, spellPracticeTotal2, gameModeGet2, gameModeTotal2
		time.Sleep(time.Second / 2)
	}
}

func getCount(hand windows.Handle, baseAddress uintptr) (spellPracticeGet, spellPracticeTotal, gameModeGet, gameModeTotal uint32) {
	_, gameModeGet = readMemory[uint32](hand, baseAddress, 0xCF41C, 0x998)
	_, spellPracticeGet = readMemory[uint32](hand, baseAddress, 0xCF41C, 0x99C)
	_, gameModeTotal = readMemory[uint32](hand, baseAddress, 0xCF41C, 0x9A0)
	_, spellPracticeTotal = readMemory[uint32](hand, baseAddress, 0xCF41C, 0x9A4)
	return
}
