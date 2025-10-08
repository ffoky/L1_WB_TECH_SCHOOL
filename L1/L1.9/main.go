package main

import (
	"fmt"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		defer close(ch2)
		for num := range ch1 {
			ch2 <- num * 2
		}
	}()

	arr := [3]int{1, 2, 3}
	go func() {
		defer close(ch1)
		for _, value := range arr {
			ch1 <- value
		}
	}()

	for value := range ch2 {
		fmt.Println(value)
	}
}
