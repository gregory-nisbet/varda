// Copyright 2025 Gregory Nisbet. All rights reserved. This code is MIT Licensed.
package varda

import (
	"fmt"
	"os"
	"testing"
)

func TestIsSymlink_NotASymlink(t *testing.T) {
	t.Parallel()
	td := t.TempDir()

	path := fmt.Sprintf("%s/%s", td, "a.txt")
	fh, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	if _, err := fh.WriteString("aaaa"); err != nil {
		panic(err)
	}
	if err := fh.Close(); err != nil {
		panic(err)
	}
	sym, err := IsSymlink(path)
	if err != nil {
		t.Error(err)
	}
	if sym {
		t.Error("erroneously reported as symlink")
	}
}

func TestIsSymlink_Symlink(t *testing.T) {
	t.Parallel()
	td := t.TempDir()

	path := fmt.Sprintf("%s/%s", td, "a.txt")
	otherPath := fmt.Sprintf("%s/%s", td, "b.txt")
	fh, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	if _, err := fh.WriteString("aaaa"); err != nil {
		panic(err)
	}
	if err := fh.Close(); err != nil {
		panic(err)
	}
	if err := os.Symlink(path, otherPath); err != nil {
		panic(err)
	}

	sym, err := IsSymlink(otherPath)
	if err != nil {
		t.Error(err)
	}
	if !sym {
		t.Error("symlink should be symlink")
	}
}
