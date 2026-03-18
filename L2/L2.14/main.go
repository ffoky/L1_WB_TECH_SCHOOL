package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v\n", time.Since(start))
}

func or(channels ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	once := &sync.Once{}

	for _, ch := range channels {
		go func() {
			<-ch
			//если на закрою через sync.Once
			// то вомзожно второе закрытие канала
			once.Do(func() { close(out) })
		}()
	}

	return out
}
