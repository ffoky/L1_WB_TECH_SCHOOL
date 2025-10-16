package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	data := []float64{
		-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5,
	}

	ans := make(map[int][]float64)
	for _, temp := range data {
		group := int(temp) / 10 * 10
		ans[group] = append(ans[group], temp)
	}

	keys := make([]int, 0, len(ans))
	for key := range ans {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	for _, key := range keys {
		fmt.Fprintf(writer, "%d:{", key)
		for i, value := range ans[key] {
			if i > 0 {
				fmt.Fprintf(writer, ", ")
			}
			fmt.Fprintf(writer, "%.1f", value)
		}
		fmt.Fprint(writer, "}, ")
	}
}
