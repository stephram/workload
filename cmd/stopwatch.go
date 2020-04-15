package main

import (
	"fmt"
	"time"
)

func main() {
	defer stopWatch(time.Now(), "Waited %s")
}

func stopWatch(start time.Time, message string) {
	fmt.Printf(message, time.Since(start))
}
