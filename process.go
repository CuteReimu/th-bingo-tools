package main

import (
	"golang.org/x/sys/windows"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// 根据进程名获取进程ID
func getPidByProcessName(processName string) (int64, error) {
	task := exec.Command("cmd", "/c", "wmic", "process", "get", "name,", "ProcessId", "|", "findstr", processName)
	data, _ := task.CombinedOutput()
	res := strings.Split(string(data), "\n")[0] //取第一行程序结果
	ss := regexp.MustCompile(`\s+`).Split(res, -1)
	pidStr := ss[1]
	return strconv.ParseInt(pidStr, 10, 64)
}

// 获取进程的句柄
func getProcessHandle(pid uint32) (windows.Handle, error) {
	return windows.OpenProcess(syscall.STANDARD_RIGHTS_ALL|0xFFFF, false, pid)
}

// 获取模块基地址
func getModuleBaseAddress(hand windows.Handle, processName string) (uintptr, error) {
	hModel := [10000]windows.Handle{0}
	var num uint32
	if err := windows.EnumProcessModules(hand, &hModel[0], uint32(len(hModel)), &num); err != nil {
		return 0, err
	}
	for i := uint32(0); i < num; i++ {
		tmp := [50]uint16{0}
		if err := windows.GetModuleBaseName(hand, hModel[i], &tmp[0], uint32(len(tmp))); err != nil {
			log.Println(err)
			continue
		}
		sb := strings.Builder{}
		for _, v := range tmp {
			if v != 0 {
				sb.WriteByte(byte(v))
			}
		}
		if strings.EqualFold(processName, sb.String()) {
			return uintptr(hModel[i]), nil
		}
	}
	return 0, nil
}

// 通过基址+指针链读取到指针地址的值
func readMemory[T any](handle windows.Handle, baseAddress uintptr, addrLinks ...uintptr) (uintptr, T) {
	addr := baseAddress
	for _, x := range addrLinks[:len(addrLinks)-1] {
		var tmp64 uint32
		addr += x
		//fmt.Printf("地址%X\t", addr)
		windows.ReadProcessMemory(handle, addr, (*byte)(unsafe.Pointer(&tmp64)), unsafe.Sizeof(tmp64), nil)
		addr = uintptr(tmp64)
	}
	addr += addrLinks[len(addrLinks)-1]
	//fmt.Printf("地址%X\n", addr)
	var value T
	windows.ReadProcessMemory(handle, addr, (*byte)(unsafe.Pointer(&value)), unsafe.Sizeof(value), nil)
	return addr, value
}
