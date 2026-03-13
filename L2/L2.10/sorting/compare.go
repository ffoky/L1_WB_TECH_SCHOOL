package sorting

import (
	"cmp"
	"strconv"
	"strings"
)

// monthIndex содержит порядковые номера месяцев.
var monthIndex = map[string]int{
	"Jan": 1, "Feb": 2, "Mar": 3,
	"Apr": 4, "May": 5, "Jun": 6,
	"Jul": 7, "Aug": 8, "Sep": 9,
	"Oct": 10, "Nov": 11, "Dec": 12,
}

// sizeMap содержит множители для человекочитаемых суффиксов размера.
//
//	K —килобайт  1 << 10 = 1 024
//	M  мегабайт  1 << 20 = 1 048 576
//	G  гигабайт  1 << 30 = 1 073 741 824
//	T  терабайт  1 << 40 = 1 099 511 627 776
var sizeMap = map[byte]int{
	'K': 1 << 10, 'k': 1 << 10,
	'M': 1 << 20, 'm': 1 << 20,
	'G': 1 << 30, 'g': 1 << 30,
	'T': 1 << 40, 't': 1 << 40,
}

// makeCmpFunc возвращает функцию для сравнения.
// Учитывает столбец -k, обрезку пробелов -b, режим сравнения -n/-M/-h и обратную -r.
func makeCmpFunc(opts Options) func(a, b string) int {
	return func(a, b string) int {
		fa := extractField(a, opts.Column)
		fb := extractField(b, opts.Column)

		if opts.TrimTailSpace {
			fa = strings.TrimRight(fa, " \t")
			fb = strings.TrimRight(fb, " \t")
		}

		var result int
		switch {
		case opts.Numeric:
			result = compareNumeric(fa, fb)
		case opts.Month:
			result = compareMonth(fa, fb)
		case opts.HumanReadableSize:
			result = compareHuman(fa, fb)
		default:
			result = strings.Compare(fa, fb)
		}

		if opts.Reverse {
			return -result
		}
		return result
	}
}

// extractField возвращает n-ю колонку.
// Если col == 0, возвращает всю строку. Если столбца нет — пустую строку.
func extractField(line string, col int) string {
	if col == 0 {
		return line
	}
	cols := strings.Split(line, "\t")
	if col > len(cols) {
		return ""
	}
	return cols[col-1]
}

// compareNumeric сравнивает два элемента как float64, т.е. -n.
// Непарсируемые значения считаются равными 0.
func compareNumeric(a, b string) int {
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		af = 0
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		bf = 0
	}
	return cmpOrdered(af, bf)
}

// compareMonth сравнивает числа месяцев -M.
// Не месяц = 0 .
func compareMonth(a, b string) int {
	return cmpOrdered(monthIndex[a], monthIndex[b])
}

// compareHuman сравнивает размеры  двух элементов.
func compareHuman(a, b string) int {
	return cmpOrdered(parseSize(a), parseSize(b))
}

// parseSize парсит человекочитаемый размер в байты.
// Возвращает 0 если строка не распознана как в GNU.
func parseSize(s string) int {
	if s == "" {
		return 0
	}
	last := s[len(s)-1]
	if sizeInBytes, ok := sizeMap[last]; ok {
		n, err := strconv.Atoi(s[:len(s)-1])
		if err != nil {
			return 0
		}
		return n * sizeInBytes
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}

// cmpOrdered сравнивает два сравниваемых типа в стиле GNU.
func cmpOrdered[T cmp.Ordered](a, b T) int {
	if a > b {
		return 1
	} else if a < b {
		return -1
	}
	return 0
}
