package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	A := []int{1, 2, 3}
	B := []int{2, 3, 4}

	mapA := make(map[int]struct{})
	for _, v := range A {
		mapA[v] = struct{}{}
	}

	intersection := make(map[int]struct{})
	for _, v := range B {
		if _, ok := mapA[v]; ok {
			intersection[v] = struct{}{}
		}
	}

	for v := range intersection {
		fmt.Fprint(writer, v)
	}
}
