package main

import (
	"errors"
	"log/slog"
	"strings"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

var errProcessNotFound = errors.New("没有找到进程")
var errModuleNotFound = errors.New("没有找到模块")

// gameProcessState 用于跟踪每个游戏的进程检测状态，避免重复打印日志
var gameProcessState sync.Map

// findGameProcess 尝试通过多个可能的进程名查找游戏进程，返回 pid、匹配的进程名和句柄
// 如果找到进程，会打日志；如果之前找到过但现在找不到了，也会打日志
func findGameProcess(gameId string, names []string) (pid uint32, matchedName string, handle windows.Handle, baseAddress uintptr, err error) {
	for _, name := range names {
		pid, err = getPidByProcessName(name)
		if err != nil {
			continue
		}
		handle, err = getProcessHandle(pid)
		if err != nil {
			continue
		}
		baseAddress, err = getModuleBaseAddress(handle, name)
		if err != nil {
			windows.CloseHandle(handle)
			continue
		}
		matchedName = name
		// 检测到进程时打日志（仅首次）
		if _, loaded := gameProcessState.LoadOrStore(gameId, true); !loaded {
			slog.Info("检测到游戏进程", "game", gameId, "exe", name, "pid", pid)
		}
		return pid, matchedName, handle, baseAddress, nil
	}
	// 如果之前检测到过，现在找不到了
	if _, loaded := gameProcessState.LoadAndDelete(gameId); loaded {
		slog.Info("游戏进程已关闭", "game", gameId)
	}
	return 0, "", 0, 0, errProcessNotFound
}

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
