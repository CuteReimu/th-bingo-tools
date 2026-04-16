package main

import (
	"errors"
	"log/slog"
	"strings"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	errProcessNotFound = errors.New("进程未找到")
	errModuleNotFound  = errors.New("模块未找到")
)

// processResult 封装 findGameProcess 的返回结果
type processResult struct {
	PID         uint32
	ExeName     string
	Handle      windows.Handle
	BaseAddress uintptr
}

// gameProcessState 用于跟踪每个游戏的进程检测状态，避免重复打印日志
var gameProcessState sync.Map

// findGameProcess 尝试通过多个可能的进程名查找游戏进程
func findGameProcess(gameID string, names []string) (processResult, error) {
	for _, name := range names {
		pid, err := findProcessByName(name)
		if err != nil {
			continue
		}
		handle, err := openProcess(pid)
		if err != nil {
			continue
		}
		baseAddress, err := getModuleBaseAddress(handle, name)
		if err != nil {
			windows.CloseHandle(handle)
			continue
		}
		if _, loaded := gameProcessState.LoadOrStore(gameID, true); !loaded {
			slog.Info("检测到游戏进程", "game", gameID, "exe", name, "pid", pid)
		}
		return processResult{PID: pid, ExeName: name, Handle: handle, BaseAddress: baseAddress}, nil
	}
	if _, loaded := gameProcessState.LoadAndDelete(gameID); loaded {
		slog.Info("游戏进程已关闭", "game", gameID)
	}
	return processResult{}, errProcessNotFound
}

// findProcessByName 通过进程快照查找指定名称的进程
func findProcessByName(processName string) (uint32, error) {
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

// openProcess 以读取权限打开进程
func openProcess(pid uint32) (windows.Handle, error) {
	return windows.OpenProcess(windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, pid)
}

// getModuleBaseAddress 获取指定模块的基地址
func getModuleBaseAddress(handle windows.Handle, processName string) (uintptr, error) {
	var modules [1024]windows.Handle
	var needed uint32
	if err := windows.EnumProcessModules(handle, &modules[0], uint32(unsafe.Sizeof(modules)), &needed); err != nil {
		return 0, err
	}
	count := needed / uint32(unsafe.Sizeof(modules[0]))
	for i := range count {
		var name [256]uint16
		if err := windows.GetModuleBaseName(handle, modules[i], &name[0], uint32(len(name))); err != nil {
			continue
		}
		if strings.EqualFold(processName, windows.UTF16ToString(name[:])) {
			return uintptr(modules[i]), nil
		}
	}
	return 0, errModuleNotFound
}

// readMemory 通过基址+指针链读取进程内存
// addrLinks 中除最后一个元素外，每个元素都表示一次指针解引用
func readMemory[T any](out *T, handle windows.Handle, baseAddress uintptr, addrLinks ...uintptr) uintptr {
	addr := baseAddress
	for _, offset := range addrLinks[:len(addrLinks)-1] {
		var ptr uint32
		addr += offset
		if err := windows.ReadProcessMemory(handle, addr, (*byte)(unsafe.Pointer(&ptr)), unsafe.Sizeof(ptr), nil); err != nil {
			return 0
		}
		addr = uintptr(ptr)
	}
	addr += addrLinks[len(addrLinks)-1]
	if err := windows.ReadProcessMemory(handle, addr, (*byte)(unsafe.Pointer(out)), unsafe.Sizeof(*out), nil); err != nil {
		return 0
	}
	return addr
}
