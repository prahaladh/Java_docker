package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// Struct to hold the results of the API calls for each application
type ApiResponse struct {
	ApplicationName string
	API1Response    string
	API2Response    string
	API3Response    string
	API4Response    string
}

// Function to make an HTTP GET request and return the response body as a string
func hitAPI(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func main() {
	// Application names (you can add more or modify these)
	applicationNames := []string{"App1", "App2", "App3", "App4", "App5", "App6"}

	// URLs for the 4 APIs (replace these with actual API endpoints)
	api1 := "https://api1.example.com/endpoint"
	api2 := "https://api2.example.com/endpoint"
	api3 := "https://api3.example.com/endpoint"
	api4 := "https://api4.example.com/endpoint"

	// Slice to store results of all iterations
	results := make([]ApiResponse, len(applicationNames))

	// Iterate through the application names
	for i, appName := range applicationNames {
		var wg sync.WaitGroup
		wg.Add(4) // For 4 API calls

		// Create an instance of ApiResponse to store results for this application
		var response ApiResponse
		response.ApplicationName = appName // Set the application name for this iteration

		// Goroutine for hitting API 1
		go func() {
			defer wg.Done()
			api1Result, err := hitAPI(api1)
			if err != nil {
				fmt.Println("Error hitting API 1:", err)
				return
			}
			response.API1Response = api1Result
		}()

		// Goroutine for hitting API 2
		go func() {
			defer wg.Done()
			api2Result, err := hitAPI(api2)
			if err != nil {
				fmt.Println("Error hitting API 2:", err)
				return
			}
			response.API2Response = api2Result
		}()

		// Goroutine for hitting API 3
		go func() {
			defer wg.Done()
			api3Result, err := hitAPI(api3)
			if err != nil {
				fmt.Println("Error hitting API 3:", err)
				return
			}
			response.API3Response = api3Result
		}()

		// Goroutine for hitting API 4
		go func() {
			defer wg.Done()
			api4Result, err := hitAPI(api4)
			if err != nil {
				fmt.Println("Error hitting API 4:", err)
				return
			}
			response.API4Response = api4Result
		}()

		// Wait for all goroutines to finish
		wg.Wait()

		// Store the response in the results slice
		results[i] = response
	}

	// Output the results of all iterations
	for i, res := range results {
		fmt.Printf("Iteration %d (%s) Results: %+v\n", i+1, res.ApplicationName, res)
	}
}
