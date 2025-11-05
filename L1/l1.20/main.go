package main

import "strings"

func reverseWords(s string) string {
	runes := []rune(s)
	var result strings.Builder

	for i := len(runes) - 1; i >= 0; {
		if runes[i] == ' ' {
			i--
			continue
		}

		end := i + 1
		for i >= 0 && runes[i] != ' ' {
			i--
		}
		start := i + 1

		result.WriteString(string(runes[start:end]))
		if i >= 0 {
			result.WriteRune(' ')
		}
	}

	return result.String()
}

func main() {
	input := "snow dog sun"
	print(reverseWords(input))
}
