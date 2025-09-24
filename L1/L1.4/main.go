package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"time"

	"syscall"
)

const workerNum = 10

func worker(ctx context.Context, jobs chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			fmt.Println(job * job)
			time.Sleep(1000 * time.Millisecond)
		}

	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()
	var wg sync.WaitGroup

	jobs := make(chan int, workerNum)
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go worker(ctx, jobs, &wg)
	}

	go func() {
		defer close(jobs)
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				return
			case jobs <- i:
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	<-ctx.Done()
	wg.Wait()
}
