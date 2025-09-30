package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

func main() {
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	data := make(map[int]int)
	var concurrentMap sync.Map
	n := 100
	var mu sync.Mutex
	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup
	for i := 0; i < n; i++ {
		wg1.Add(1)
		go func(i int) {
			defer wg1.Done()
			mu.Lock()
			defer mu.Unlock()
			data[i] = i
		}(i)
		wg2.Add(1)
		go func(i int) {
			defer wg2.Done()
			concurrentMap.Store(i, i)
		}(i)
	}
	wg1.Wait()
	wg2.Wait()

	fmt.Fprintln(writer, len(data))
	count := 0
	concurrentMap.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	fmt.Fprintln(writer, count)
}
