package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	given := []string{"cat", "cat", "dog", "cat", "tree"}
	set := make(map[string]struct{})
	for _, s := range given {
		set[s] = struct{}{}
	}
	res := make([]string, 0, len(set))
	for v := range set {
		res = append(res, v)
	}
	fmt.Printf("{%s}\n", strings.Join(res, ", "))
}
