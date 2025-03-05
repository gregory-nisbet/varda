// Copyright 2025 Gregory Nisbet. All rights reserved. This code is MIT Licensed.
package varda

type RingBuffer struct {
	size   int
	index  int
	hopper []string
}

func NewRingBuffer(n int) *RingBuffer {
	if n < 0 {
		return nil
	}
	return &RingBuffer{
		size:   0,
		index:  0,
		hopper: make([]string, n, n),
	}
}

func (b *RingBuffer) Add(item string) int {
	if b.index+1 == len(b.hopper) {
		b.index = 0
		b.hopper[0] = item
		return len(b.hopper)
	}
	b.size++
	b.index++
	b.hopper[b.index] = item
	if b.size > len(b.hopper) {
		b.size--
		return len(b.hopper)
	}
	return b.size
}

func (b *RingBuffer) ToSlice() []string {
	out := make([]string, len(b.hopper))
	index := b.index
	for ii := 0; ii < b.size; ii++ {
		index++
		if index == len(b.hopper) {
			index -= len(b.hopper)
		}
		out[ii] = b.hopper[index]
	}
	return out
}
