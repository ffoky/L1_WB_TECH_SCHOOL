package main

import "fmt"

func deleteIthElemFromSlice(i int, slice []int) []int {
	if i < 0 || i >= len(slice) {
		return slice
	}

	res := make([]int, len(slice)-1)
	copy(res, append(slice[:i], slice[i+1:]...))
	return res
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8}
	fmt.Print(deleteIthElemFromSlice(2, nums))
}
