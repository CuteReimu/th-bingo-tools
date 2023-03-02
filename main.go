package main

func main() {
	pid, err := getPidByProcessName("chrome.exe")
	if err != nil {
		panic(err)
	}
	getProcessHandle(uint32(pid))
}
