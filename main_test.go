package main

import (
	"testing"
)

func TestMain(t *testing.T) {

	assert(t, true, true)
}

func assert(t *testing.T, val1 any, val2 any) {
	if val1 != val2 {
		t.Fail()
	}
}
