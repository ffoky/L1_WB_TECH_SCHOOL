package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func StopByCondition(condition *bool) {
	counter := 0
	for !*condition {
		fmt.Printf("StopByCondition функция вычисляет что-то %d\n", counter)
		counter++
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("В функции StopByCondition условие выполнилось, горутина останавливается")
}

func StopByDoneChannel(done <-chan struct{}) {

	for i := 0; ; i++ {
		select {
		case <-done:
			fmt.Println("StopByDoneChannel: получен сигнал, завершаю работу")
			return
		default:
			fmt.Printf("StopByDoneChannel: вычисляю что-то %d\n", i)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func StopByContext(ctx context.Context) {
	counter := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("StopByContext функция завершает работу через отмену контекста: %v\n", ctx.Err())
			return
		default:
			fmt.Printf("StopByContext что-то делает: %d\n", counter)
			counter++
			time.Sleep(300 * time.Millisecond)
		}
	}
}

func StopByContextTimeout(ctx context.Context) {
	counter := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("В StopByContextTimeout функции истек таймаут: %v\n", ctx.Err())
			return
		default:
			fmt.Printf("StopByContextTimeout функция вычисляет: %d\n", counter)
			counter++
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func StopByRuntimeGoexit() {

	for i := 0; ; i++ {
		fmt.Printf("StopByRuntimeGoexit функция вычисляет %d\n", i)
		time.Sleep(500 * time.Millisecond)

		if i >= 3 {
			fmt.Println("StopByRuntimeGoexit вызов runtime.Goexit()")
			runtime.Goexit()
		}
	}
}

func StopByPanic() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic восстановление: %v\n", r)
		}
	}()

	counter := 0
	for {
		fmt.Printf("StopByPanic функция выполняет вычисления: %d\n", counter)
		counter++
		time.Sleep(200 * time.Millisecond)

		if counter >= 3 {
			panic("StopByPanic функция паникует")
		}
	}
}

func StopByClosingDataChannel(data <-chan int) {
	for {
		select {
		case val, ok := <-data:
			if !ok {
				fmt.Println("Завершение горутины StopByClosingDataChannel")
				return
			}
			fmt.Printf("DataChannel: получено значение %d\n", val)
		}
	}
}

func WorkerWithWaitGroup(id int, stop *bool, wg *sync.WaitGroup) {
	defer wg.Done()
	counter := 0
	for !*stop {
		fmt.Printf("Worker-%d: работаю %d\n", id, counter)
		counter++
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Printf("Worker-%d: завершил работу\n", id)
}

func main() {
	fmt.Println("Остановка по условию")
	var condition bool
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		StopByCondition(&condition)
	}()
	time.Sleep(1 * time.Second)
	condition = true
	wg.Wait()

	fmt.Println("\nОстановка через канал")
	done := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		StopByDoneChannel(done)
	}()
	time.Sleep(1 * time.Second)
	fmt.Println("Закрытие канала")
	close(done)
	wg.Wait()

	fmt.Println("\nОстановка горутины отменой контекста")
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		StopByContext(ctx)
	}()
	time.Sleep(1 * time.Second)
	fmt.Println("Отмена контекста")
	cancel()
	wg.Wait()

	fmt.Println("\nОстановка по таймауту")
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelTimeout()
	wg.Add(1)
	go func() {
		defer wg.Done()
		StopByContextTimeout(ctxTimeout)
	}()
	wg.Wait()

	fmt.Println("\nОстановка через runtime.Goexit")
	wg.Add(1)
	go func() {
		defer wg.Done()
		StopByRuntimeGoexit()
	}()
	wg.Wait()

	fmt.Println("\nОстановка через панику")
	wg.Add(1)
	go func() {
		defer wg.Done()
		StopByPanic()
	}()
	wg.Wait()

	fmt.Println("\nОстановка при закрытии канала")
	data := make(chan int)
	wg.Add(1)
	go func() {
		defer wg.Done()
		StopByClosingDataChannel(data)
	}()

	go func() {
		for i := 1; i <= 3; i++ {
			data <- i
			time.Sleep(400 * time.Millisecond)
		}
		fmt.Println("Закрываю канал данных")
		close(data)
	}()
	wg.Wait()
}
