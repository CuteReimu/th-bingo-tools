package main

import (
	"golang.org/x/sys/windows"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

// 根据进程名获取进程ID
func getPidByProcessName(name string) (int64, error) {
	task := exec.Command("cmd", "/c", "wmic", "process", "get", "name,", "ProcessId", "|", "findstr", name)
	data, _ := task.CombinedOutput()
	res := strings.Split(string(data), "\n")[0] //取第一行程序结果
	ss := regexp.MustCompile(`\s+`).Split(res, -1)
	pidStr := ss[1]
	return strconv.ParseInt(pidStr, 10, 64)
}

// 获取进程的句柄
func getProcessHandle(pid uint32) windows.Handle {
	hand, err := windows.OpenProcess(syscall.STANDARD_RIGHTS_ALL|0xFFFF, false, pid)
	if err != nil {
		if err.Error() != "The operation completed successfully." {
			log.Println("报错：", err.Error())
			syscall.Exit(1)
		}
	}
	if hand <= 0 {
		log.Println("打开句柄失败", hand)
		syscall.Exit(1)
	}
	log.Println("[+] 打开目标进程成功", hand)
	return hand
}
