package main

import (
	"fmt"

	openhpibadge "github.com/romnnn/open-hpi-badge"
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
