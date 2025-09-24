package main

import (
	"fmt"
	"time"
)

const N = 5

func reader(ch chan int) {
	for data := range ch {
		fmt.Println(data)
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	ch := make(chan int)

	deadline := time.After(N * time.Second)
	go reader(ch)
	for i := 0; ; i++ {
		select {
		case ch <- i:
		case <-deadline:
			close(ch)
			time.Sleep(100 * time.Millisecond)
			return
		}
	}
}
