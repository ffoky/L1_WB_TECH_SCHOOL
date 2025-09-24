package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func worker(jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Println(job * job)
	}
}

func main() {
	workerNum := flag.Int("workers", 3, "workerNum")
	flag.Parse()
	if *workerNum < 1 {
		fmt.Println("No worker num provided")
		return
	}
	ch := make(chan int, *workerNum)

	var wg sync.WaitGroup

	for i := 0; i < *workerNum; i++ {
		wg.Add(1)
		go worker(ch, &wg)
	}

	go func() {
		for i := 0; ; i++ {
			ch <- i
			time.Sleep(500 * time.Millisecond)
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	close(ch)

	wg.Wait()
}
