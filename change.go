package main

import (
	"fmt"
)

// Example struct with three fields
type Example struct {
	Field1 string
	Field2 string
	Status string
}

// Function to compare Field1 and Field2 for each struct in an array
// and update Status based on the comparison
func CompareAndUpdateStatus(examples []Example) {
	for i := range examples {
		if examples[i].Field1 == examples[i].Field2 {
			examples[i].Status = "Same"
		} else {
			examples[i].Status = "Different"
		}
	}
}

func main() {
	// Array of Example structs with different values for testing
	examples := []Example{
		{Field1: "apple", Field2: "apple"},
		{Field1: "apple", Field2: "orange"},
		{Field1: "Apple", Field2: "apple"}, // Case-sensitive comparison
		{Field1: "banana", Field2: "banana"},
	}

	// Perform comparison and update status for each struct in the array
	CompareAndUpdateStatus(examples)

	// Print results
	for i, example := range examples {
		fmt.Printf("Example %d - Field1: %s, Field2: %s, Status: %s\n", i+1, example.Field1, example.Field2, example.Status)
	}
}
