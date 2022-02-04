package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	myMutex := sync.Mutex{}

	myMutex.Lock()

	go func() {
		myMutex.Lock()
		fmt.Println("Printing from the goroutine 1")
		myMutex.Unlock()
	}()

	go func() {
		myMutex.Lock()
		fmt.Println("Printing from the goroutine 2")
		myMutex.Unlock()
	}()

	fmt.Println("Printing from the main routine.")
	myMutex.Unlock()
	time.Sleep(time.Second * 1)

}
