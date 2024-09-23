package main

import (
	"fmt"
	"sync"
)

func workerA(id int, wg *sync.WaitGroup, ch chan<- string) {
	defer wg.Done() // Mark the goroutine as done when it completes

	// Simulate work and return a string
	result := fmt.Sprintf("WorkerA %d: processed", id)
	
	// Send the result to the channel
	ch <- result
}

func workerB(id int, wg *sync.WaitGroup, ch chan<- string) {
	defer wg.Done() // Mark the goroutine as done when it completes

	// Simulate work and return a string
	result := fmt.Sprintf("WorkerB %d: processed", id)
	
	// Send the result to the channel
	ch <- result
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan string, 10) // Buffered channel for results

	// Create slices to store the results from workerA and workerB
	var workerAResults []string
	var workerBResults []string

	// Start multiple workerA goroutines
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go workerA(i, &wg, ch)
	}

	// Start multiple workerB goroutines
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go workerB(i, &wg, ch)
	}

	// Close the channel after all goroutines have finished
	go func() {
		wg.Wait()   // Wait for all goroutines to finish
		close(ch)   // Close the channel
	}()

	// Collect the results from the channel
	for result := range ch {
		// Separate results based on worker type
		if result[:7] == "WorkerA" {
			workerAResults = append(workerAResults, result)
		} else if result[:7] == "WorkerB" {
			workerBResults = append(workerBResults, result)
		}
	}

	// Print the saved results
	fmt.Println("Results from WorkerA:")
	for _, res := range workerAResults {
		fmt.Println(res)
	}

	fmt.Println("\nResults from WorkerB:")
	for _, res := range workerBResults {
		fmt.Println(res)
	}
}
