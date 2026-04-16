package main

import (
	"errors"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

var errProcessNotFound = errors.New("没有找到进程")
var errModuleNotFound = errors.New("没有找到模块")

// 根据进程名获取进程ID，使用 Windows 原生 API 而非 wmic 外部命令
func getPidByProcessName(processName string) (uint32, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, err
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	if err = windows.Process32First(snapshot, &entry); err != nil {
		return 0, err
	}
	for {
		name := windows.UTF16ToString(entry.ExeFile[:])
		if strings.EqualFold(name, processName) {
			return entry.ProcessID, nil
		}
		if err = windows.Process32Next(snapshot, &entry); err != nil {
			break
		}
	}
	return 0, errProcessNotFound
}

// 获取进程的句柄
func getProcessHandle(pid uint32) (windows.Handle, error) {
	return windows.OpenProcess(windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, pid)
}

// 获取模块基地址
func getModuleBaseAddress(hand windows.Handle, processName string) (uintptr, error) {
	var hModel [10000]windows.Handle
	var num uint32
	if err := windows.EnumProcessModules(hand, &hModel[0], uint32(len(hModel)), &num); err != nil {
		return 0, err
	}
	for i := range num {
		var tmp [50]uint16
		if err := windows.GetModuleBaseName(hand, hModel[i], &tmp[0], uint32(len(tmp))); err != nil {
			continue
		}
		if strings.EqualFold(processName, windows.UTF16ToString(tmp[:])) {
			return uintptr(hModel[i]), nil
		}
	}
	return 0, errModuleNotFound
}

// 通过基址+指针链读取到指针地址的值
func readMemory[T any](out *T, handle windows.Handle, baseAddress uintptr, addrLinks ...uintptr) uintptr {
	addr := baseAddress
	for _, x := range addrLinks[:len(addrLinks)-1] {
		var tmp64 uint32
		addr += x
		if err := windows.ReadProcessMemory(handle, addr, (*byte)(unsafe.Pointer(&tmp64)), unsafe.Sizeof(tmp64), nil); err != nil {
			return 0
		}
		addr = uintptr(tmp64)
	}
	addr += addrLinks[len(addrLinks)-1]
	if err := windows.ReadProcessMemory(handle, addr, (*byte)(unsafe.Pointer(out)), unsafe.Sizeof(*out), nil); err != nil {
		return 0
	}
	return addr
}
