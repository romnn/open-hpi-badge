package main

import (
	"fmt"

	"github.com/romnnn/openhpibadge"
)

func run() int {
	course, err := openhpibadge.ScrapeMOOCByName("neuralnets2020")
	if err != nil {
		panic(err)
	}
	return course.Participants.Current
}

func main() {
	fmt.Printf("Current participants: %d\n", run())
}
