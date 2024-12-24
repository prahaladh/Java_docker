package main

import (
	"fmt"
	"regexp"
)

func main() {
	// Define the regex pattern for Java exceptions
	pattern := `(?i)(\b\w+Exception\b|Caused by: \w+Exception)`

	// Compile the regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}

	// Input string
	log := `
		Some log text here
		NullPointerException occurred while accessing a null object.
		Caused by: IllegalArgumentException: Invalid argument passed.
	`

	// Check if the pattern matches
	if re.MatchString(log) {
		fmt.Println("Java exception found!")
	} else {
		fmt.Println("No Java exceptions found.")
	}

	// Find all exceptions
	matches := re.FindAllString(log, -1)
	fmt.Println("Exceptions found:", matches)
}
