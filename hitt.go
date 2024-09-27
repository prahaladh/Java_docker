package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// Struct to hold the results of the API calls for each application
type ApiResponse struct {
	ApplicationName string `json:"application_name"`
	API1Response    string `json:"api1_response"`
	API2Response    string `json:"api2_response"`
	API3Response    string `json:"api3_response"`
	API4Response    string `json:"api4_response"`
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

// Function to handle the API calls and send the result back via a channel
func fetchAPIResults(appName, api1, api2, api3, api4 string, ch chan<- ApiResponse, wg *sync.WaitGroup) {
	defer wg.Done()

	var response ApiResponse
	response.ApplicationName = appName

	var wgAPIs sync.WaitGroup
	wgAPIs.Add(4) // For 4 API calls

	// Goroutine for hitting API 1
	go func() {
		defer wgAPIs.Done()
		api1Result, err := hitAPI(api1)
		if err != nil {
			fmt.Println("Error hitting API 1:", err)
			return
		}
		response.API1Response = api1Result
	}()

	// Goroutine for hitting API 2
	go func() {
		defer wgAPIs.Done()
		api2Result, err := hitAPI(api2)
		if err != nil {
			fmt.Println("Error hitting API 2:", err)
			return
		}
		response.API2Response = api2Result
	}()

	// Goroutine for hitting API 3
	go func() {
		defer wgAPIs.Done()
		api3Result, err := hitAPI(api3)
		if err != nil {
			fmt.Println("Error hitting API 3:", err)
			return
		}
		response.API3Response = api3Result
	}()

	// Goroutine for hitting API 4
	go func() {
		defer wgAPIs.Done()
		api4Result, err := hitAPI(api4)
		if err != nil {
			fmt.Println("Error hitting API 4:", err)
			return
		}
		response.API4Response = api4Result
	}()

	// Wait for all API calls to finish
	wgAPIs.Wait()

	// Send the response to the buffered channel
	ch <- response
}

func main() {
	// Application names (you can add more or modify these)
	applicationNames := []string{"App1", "App2", "App3", "App4", "App5", "App6"}

	// URLs for the 4 APIs (replace these with actual API endpoints)
	api1 := "https://api1.example.com/endpoint"
	api2 := "https://api2.example.com/endpoint"
	api3 := "https://api3.example.com/endpoint"
	api4 := "https://api4.example.com/endpoint"

	// Buffered channel to hold results for all applications
	ch := make(chan ApiResponse, len(applicationNames))

	var wg sync.WaitGroup
	wg.Add(len(applicationNames)) // One goroutine per application

	// Iterate through the application names and start the goroutines
	for _, appName := range applicationNames {
		go fetchAPIResults(appName, api1, api2, api3, api4, ch, &wg)
	}

	// Wait for all fetch routines to complete
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collect results from the buffered channel
	var results []ApiResponse
	for result := range ch {
		results = append(results, result)
	}

	// Convert results to JSON
	jsonData, err := json.Marshal(results)
	if err != nil {
		fmt.Println("Error converting results to JSON:", err)
		return
	}

	// Print the JSON data
	fmt.Println(string(jsonData))
}
