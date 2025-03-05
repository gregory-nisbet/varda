// Copyright 2025 Gregory Nisbet. All rights reserved. This code is MIT Licensed.
package varda

import (
	"fmt"
	"strings"
	"testing"
)

func TestRingBuffer(t *testing.T) {
	t.Parallel()

	rb := NewRingBuffer(10)
	for i := 1; i <= 100; i++ {
		n := rb.Add(fmt.Sprintf("%d", i))
		switch {
		case i <= 9:
			if n != i {
				t.Errorf("bad size %d", n)
			}
		default:
			if n != 10 {
				t.Errorf("bad size %d", n)
			}
		}

	}

	slice := rb.ToSlice()
	actual := strings.Join(slice, ",")

	if actual != "91,92,93,94,95,96,97,98,99,100" {
		t.Errorf("unexpected slice %s", actual)
	}
}
