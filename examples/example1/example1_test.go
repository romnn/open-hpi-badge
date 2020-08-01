package main

import (
	"testing"
)

func TestCli(t *testing.T) {
	out := run()
	expected := 9496
	if out < expected {
		t.Errorf("Got %d but expected more or equal %d", out, expected)
	}
}
