// Copyright 2025 Gregory Nisbet. All rights reserved. This code is MIT Licensed.
package main

import (
	"os"

	varda "github.com/gregory-nisbet/varda/pkg/varda"
)

const symlinkLimit = 2

func main() {
	switch len(os.Args) {
	case 0:
		panic("impossible")
	case 1:
		varda.Eprintln("nothing to do")
		os.Exit(0)
	}

	pattern := os.Args[1]
	paths := os.Args[2:]
	if len(paths) == 0 {
		paths = []string{"."}
	}

	if err := varda.SearchFiles(paths, pattern, symlinkLimit); err != nil {
		varda.Eprintln(err)
		os.Exit(1)
	}
}
