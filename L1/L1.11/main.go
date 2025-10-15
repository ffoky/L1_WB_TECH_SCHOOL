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
	intersection := make([]int, 0, len(A))
	for _, valueA := range A {
		for _, valueB := range B {
			if valueA == valueB {
				intersection = append(intersection, valueA)
			}
		}
	}
	for _, v := range intersection {
		fmt.Fprint(writer, v)
	}

}
