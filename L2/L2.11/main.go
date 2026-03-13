package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	input := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	fmt.Println(groupAnagrams(input))
}

func groupAnagrams(strs []string) map[string][]string {
	n := len(strs)
	m := make(map[string][]string)
	ans := make(map[string][]string)
	for i := 0; i < n; i++ {
		word := strings.ToLower(strs[i])
		r := []rune(word)
		sort.Slice(r, func(i, j int) bool {
			return r[i] < r[j]
		})
		k := string(r)
		m[k] = append(m[k], word)
	}
	for _, words := range m {
		if len(words) > 1 {
			firstWord := words[0]
			sort.Strings(words)
			ans[firstWord] = words
		}
	}

	return ans
}
