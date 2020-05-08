package main

import (
	"fmt"
	"github.com/romnnn/openhpibadge"
)

func run() string {
	return openhpibadge.Shout("This is an example")
}

func main() {
	fmt.Println(run())
}
