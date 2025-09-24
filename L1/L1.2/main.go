package main

import (
	"bufio"
	"fmt"
	"os"
)

func Pow2(a int, ch chan int) {
	ch <- a * a
}

func main() {
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	arr := []int{2, 4, 6, 8, 10}
	length := len(arr)
	ch := make(chan int, length)
	for _, v := range arr {
		go Pow2(v, ch)
	}
	for i := 0; i < length; i++ {
		fmt.Fprint(writer, <-ch, " ")
	}
}
