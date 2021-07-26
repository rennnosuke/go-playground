package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("start concurrent process...")

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		fmt.Println("add.")

		go func(wg *sync.WaitGroup) {
			time.Sleep(1 * time.Second)
			wg.Done()

			fmt.Println("done.")
		}(&wg)
	}
	
	fmt.Println("waiting...")
	wg.Wait()
	fmt.Println("all done.")
}
