package main

import (
	"fmt"
	"regexp"
)

func main() {
	log := `
		INFO: Starting application...
		NullPointerException: Something went wrong
			at com.example.MyClass.method(MyClass.java:23)
			at com.example.AnotherClass.anotherMethod(AnotherClass.java:45)
		Caused by: IllegalArgumentException: Invalid argument passed
			at com.example.Util.helper(Util.java:10)
		AnotherRandomException: Test message
			at example.Demo.run(Demo.java:99)
	`

	pattern := `(?i)(\w+Exception(:.*)?)(\n(\s+at\s+.+))+(\nCaused by: (\w+Exception(:.*)?)(\n(\s+at\s+.+))*)*`

	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}

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
