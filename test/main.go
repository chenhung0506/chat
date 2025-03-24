package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	select {
	case msg := <-ch1:
		fmt.Println(msg)
	case msg := <-ch2:
		fmt.Println(msg)
	case <-time.After(3 * time.Second): // 超時處理
		fmt.Println("超時未收到數據")
	}
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "來自 ch1 的數據"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "來自 ch2 的數據"
	}()

}
