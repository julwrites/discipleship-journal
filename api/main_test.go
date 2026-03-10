package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestSimple(t *testing.T) {
	if 1+1 != 2 {
		t.Error("Math is broken")
	}
}
