// Copyright 2025 Gregory Nisbet. All rights reserved. This code is MIT Licensed.
package varda

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

func IsSymlink(path string) (bool, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return fi.Mode()&fs.ModeSymlink == 1, nil
}

func SearchFiles(paths []string, pattern string, symlinkLimit int) error {
	// The symlink counting logic can be cleaned up.
	pat, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	if symlinkLimit < 0 {
		return nil
	}
	for _, path := range paths {
		fi, err := os.Stat(path)
		if err != nil {
			Eprintln(err)
			continue
		}
		sym, err := IsSymlink(path)
		if err != nil {
			Eprintln(err)
			continue
		}
		if sym {
			symlinkLimit--
		}
		switch {
		case fi.IsDir():
			items, err := ReadDir(path)
			if err != nil {
				Eprintln(err)
				continue
			}
			// we're compiling the regex over and over again.
			SearchFiles(items, pattern, symlinkLimit)
			continue
		default:
			if err := SearchFile(path, pat, symlinkLimit); err != nil {
				Eprintln(err)
				continue
			}
		}
	}
	return err
}

func ReadDir(path string) ([]string, error) {
	items, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	out := make([]string, len(items))
	// not portable
	for i, item := range items {
		if item.Name() == "." || item.Name() == ".." {
			continue
		}
		out[i] = fmt.Sprintf("%s/%s", path, item.Name())
	}
	return out, nil
}

func SearchFile(path string, pattern *regexp.Regexp, symlinkLimit int) error {
	sym, err := IsSymlink(path)
	if err != nil {
		return err
	}
	if sym {
		symlinkLimit--
	}
	if symlinkLimit < 0 {
		return nil
	}
	fh, err := os.Open(path)
	if err != nil {
		return err
	}
	// The first version does not need to be fast.
	// Read the whole thing into memory.
	//
	// Future versions will read a line at a time with a
	// maximum length or something to make the program more
	// robust.
	buf, err := io.ReadAll(fh)
	if err != nil {
		return err
	}
	lines := strings.Split(string(buf), "\n")
	for i, line := range lines {
		linum := i + 1
		if pattern.MatchString(line) {
			// hardcode 80 character limit for line.
			Oprintf("%s:%d:%.80s\n", path, linum, line)
		}
	}
	return nil
}

func Oprintf(content string, args ...any) {
	_, _ = fmt.Fprintf(os.Stdout, content, args...)
}

func Eprintf(content string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, content, args...)
}

func Eprintln(args ...any) {
	_, _ = fmt.Fprintln(os.Stderr, args...)
}
