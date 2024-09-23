package main

import (
	"fmt"
	"sync"
)

// First type of worker function
func workerA(id int, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done() // Mark the goroutine as done when it completes

	fmt.Printf("WorkerA %d starting\n", id)
	
	// Simulate some work and send data to the channel
	ch <- id * 2
	
	fmt.Printf("WorkerA %d done\n", id)
}

// Second type of worker function
func workerB(id int, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done() // Mark the goroutine as done when it completes

	fmt.Printf("WorkerB %d starting\n", id)
	
	// Simulate some different work and send data to the channel
	ch <- id * 3
	
	fmt.Printf("WorkerB %d done\n", id)
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan int, 10) // Create a buffered channel with enough capacity

	// Start multiple workerA goroutines
	for i := 1; i <= 3; i++ {
		wg.Add(1) // Increment the WaitGroup counter
		go workerA(i, &wg, ch)
	}

	// Start multiple workerB goroutines
	for i := 1; i <= 3; i++ {
		wg.Add(1) // Increment the WaitGroup counter
		go workerB(i, &wg, ch)
	}

	// Close the channel once all goroutines are done
	go func() {
		wg.Wait()   // Wait for all goroutines to finish
		close(ch)   // Close the channel to signal no more values will be sent
	}()

	// Read from the channel
	for result := range ch {
		fmt.Printf("Received: %d\n", result)
	}
}
