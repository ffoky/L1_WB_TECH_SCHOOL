package sorting

import (
	"runtime"
	"slices"
	"sync"
)

// чтобы не парсить строку повторно при каждом сравнении в merge
type input struct {
	line string
	key  float64
}

// разбивает на части по числу ядер, сортирует параллельно и объединяет
func sortParallel(lines []string, opts Options, cmpFunc func(a, b string) int) {
	n := runtime.GOMAXPROCS(0)
	chunkSize := len(lines) / n

	if opts.Numeric || opts.Month || opts.HumanReadableSize {
		sortParallelByKey(lines, opts, n, chunkSize)
	} else {
		sortParallelByStr(lines, cmpFunc, n, chunkSize)
	}
}

// ключи вычисляются один раз, дальше merge сравнивает float64 напрямую
func sortParallelByKey(lines []string, opts Options, n, chunkSize int) {
	entries := make([]input, len(lines))
	for i, line := range lines {
		entries[i] = input{line: line, key: parseKey(line, opts)}
	}

	var wg sync.WaitGroup
	chunks := make([][]input, n)
	for i := 0; i < n; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if i == n-1 {
			end = len(entries)
		}
		chunks[i] = entries[start:end]

		wg.Go(func() {
			slices.SortStableFunc(chunks[i], func(a, b input) int {
				if opts.Reverse {
					if a.key > b.key {
						return -1
					} else if a.key < b.key {
						return 1
					}
					return 0
				}
				if a.key < b.key {
					return -1
				} else if a.key > b.key {
					return 1
				}
				return 0
			})
		})
	}
	wg.Wait()

	chunks = mergeChunks(chunks, func(a, b []input) []input {
		return mergeEntries(a, b, opts.Reverse)
	})

	for i, e := range chunks[0] {
		lines[i] = e.line
	}
}

// лексиграфическая сортировка для строк
func sortParallelByStr(lines []string, cmpFunc func(a, b string) int, n, chunkSize int) {
	var wg sync.WaitGroup
	chunks := make([][]string, n)
	for i := 0; i < n; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if i == n-1 {
			end = len(lines)
		}
		chunks[i] = lines[start:end]

		wg.Go(func() {
			slices.SortStableFunc(chunks[i], func(a, b string) int {
				return cmpFunc(a, b)
			})
		})
	}
	wg.Wait()

	chunks = mergeChunks(chunks, func(a, b []string) []string {
		return mergeStrings(a, b, cmpFunc)
	})

	copy(lines, chunks[0])
}

// объединяем отсортированные части попарно параллельно
func mergeChunks[T any](chunks [][]T, mergeFn func(a, b []T) []T) [][]T {
	for len(chunks) > 1 {
		next := make([][]T, 0, (len(chunks)+1)/2)
		results := make([][]T, len(chunks)/2)

		var wg sync.WaitGroup
		resultIndex := 0
		for i := 0; i < len(chunks); i += 2 {
			if i+1 < len(chunks) {
				a, b, idx := chunks[i], chunks[i+1], resultIndex
				wg.Go(func() {
					results[idx] = mergeFn(a, b)
				})
				next = append(next, nil)
				resultIndex++
			} else {
				next = append(next, chunks[i])
			}
		}
		wg.Wait()

		copy(next, results)
		chunks = next
	}
	return chunks
}

// объединяем два отсортированных слайса input, сравнивая по ключу.
func mergeEntries(a, b []input, reverse bool) []input {
	res := make([]input, len(a)+len(b))
	i, j, k := 0, 0, 0
	for i < len(a) && j < len(b) {
		less := a[i].key <= b[j].key
		if reverse {
			less = a[i].key >= b[j].key
		}
		if less {
			res[k] = a[i]
			i++
		} else {
			res[k] = b[j]
			j++
		}
		k++
	}
	copy(res[k:], a[i:])
	copy(res[k:], b[j:])
	return res
}

// просто mergeSort для строк
func mergeStrings(a, b []string, cmpFunc func(a, b string) int) []string {
	res := make([]string, len(a)+len(b))
	i, j, k := 0, 0, 0
	for i < len(a) && j < len(b) {
		if cmpFunc(a[i], b[j]) <= 0 {
			res[k] = a[i]
			i++
		} else {
			res[k] = b[j]
			j++
		}
		k++
	}
	copy(res[k:], a[i:])
	copy(res[k:], b[j:])
	return res
}
