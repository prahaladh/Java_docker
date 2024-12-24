package main

import (
	"fmt"
	"regexp"
)

func main() {
	// Example input log containing a Java exception stack trace
	log := `
		INFO: Starting application...
		ERROR: NullPointerException: Something went wrong
			at com.example.MyClass.method(MyClass.java:23)
			at com.example.AnotherClass.anotherMethod(AnotherClass.java:45)
		Caused by: IllegalArgumentException: Invalid argument passed
			at com.example.Util.helper(Util.java:10)
			at com.example.Main.main(Main.java:5)
		END OF LOG.
	`

	// Define regex pattern for Java stack trace
	pattern := `(?i)(\w+Exception(:.*?)?(\n\s+at .+)+(\nCaused by: \w+Exception(:.*?)?(\n\s+at .+)+)?)`

	// Compile the regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}

	// Find the entire stack trace
	stackTraces := re.FindAllString(log, -1)

	if len(stackTraces) > 0 {
		fmt.Println("Exception Stack Traces Found:")
		for i, trace := range stackTraces {
			fmt.Printf("Stack Trace %d:\n%s\n", i+1, trace)
		}
	} else {
		fmt.Println("No exception stack traces found.")
	}
}
