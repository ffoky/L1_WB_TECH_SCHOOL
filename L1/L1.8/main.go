package main

import "fmt"

func setBit(num int64, i uint, value int) int64 {
	if value == 1 {
		return num | (1 << i)
	}
	return num &^ (1 << i)
}

func main() {
	var num int64 = 5
	fmt.Printf("исходное число: %d (%b)\n", num, num)

	result := setBit(num, 0, 0)
	fmt.Printf("после установки 0-го бита в 0: %d (%b)\n", result, result)

	result = setBit(num, 1, 1)
	fmt.Printf("после установки 1-го бита в 1: %d (%b)\n", result, result)

	result = setBit(num, 3, 1)
	fmt.Printf("после установки 3-го бита в 1: %d (%b)\n", result, result)
}
