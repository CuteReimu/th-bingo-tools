package main

import (
	"fmt"
	"time"
)

func main() {
	for i := range listeners {
		l := listeners[i]
		go func() {
			for {
				l.Loop()
				time.Sleep(time.Second)
			}
		}()
	}
	fmt.Scanln()
}
