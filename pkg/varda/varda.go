// Copyright 2025 Gregory Nisbet. All rights reserved. This code is MIT Licensed.
package varda

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func IsSymlink(path string) (bool, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return false, fmt.Errorf("is symlink: %w", err)
	}
	return fi.Mode()&os.ModeSymlink != 0, nil
}

func SearchFiles(paths []string, pattern string, symlinkLimit int) error {
	// The symlink counting logic can be cleaned up.
	pat, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("search files: %w", err)
	}
	if symlinkLimit < 0 {
		return nil
	}
	for _, path := range paths {
		fi, err := os.Stat(path)
		if err != nil {
			Eprinterr("calling stat", err)
			continue
		}
		sym, err := IsSymlink(path)
		if err != nil {
			Eprinterr("checking symlink", err)
			continue
		}
		if sym {
			symlinkLimit--
		}
		switch {
		case fi.IsDir():
			items, err := ReadDir(path)
			if err != nil {
				Eprinterr("reading directory", err)
				continue
			}
			// we're compiling the regex over and over again.
			SearchFiles(items, pattern, symlinkLimit)
			continue
		default:
			if err := SearchFile(path, pat, symlinkLimit); err != nil {
				Eprinterr("searching for file", err)
				continue
			}
		}
	}
	return nil
}

func ReadDir(path string) ([]string, error) {
	items, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}
	out := make([]string, len(items))
	// not portable
	for i, item := range items {
		if item.Name() == "." || item.Name() == ".." {
			continue
		}
		// traversing into git or mercurial hurts real world usability.
		if item.Name() == ".git" {
			continue
		}
		if item.Name() == ".hg" {
			continue
		}
		out[i] = fmt.Sprintf("%s/%s", path, item.Name())
	}
	return out, nil
}

func SearchFile(path string, pattern *regexp.Regexp, symlinkLimit int) error {
	// If the file *name* matches the search pattern, return a special
	// line corresponding to the file.
	//
	// TODO: only do this under certain circumstances.
	//
	// Certain types of queries like: "foo.c" should be treated specially.
	// As should "builtin/whatever"
	//
	// For right now, we just check the last portion of the path.
	//
	// I may want to suppress this if the first line of the file also matches.
	// Just a thought.
	if pattern.MatchString(filepath.Base(path)) {
		Oprintf("%s:1:\n", path)
	}

	sym, err := IsSymlink(path)
	if err != nil {
		return fmt.Errorf("search file: %w", err)
	}
	if sym {
		symlinkLimit--
	}
	if symlinkLimit < 0 {
		return nil
	}
	fh, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("search file: %w", err)
	}
	// The first version does not need to be fast.
	// Read the whole thing into memory.
	//
	// Future versions will read a line at a time with a
	// maximum length or something to make the program more
	// robust.
	buf, err := io.ReadAll(fh)
	if err != nil {
		return fmt.Errorf("search file: %w", err)
	}
	lines := strings.Split(string(buf), "\n")
	for i, line := range lines {
		linum := i + 1
		if pattern.MatchString(line) {
			// print the line before and after delimited by @@ if it's short enough.
			// This can provide crucial context when scanning lines.
			//
			// I use some ridiculous heuristics to check whether we're probably looking at text
			// (vs a binary file) and then print a truncated previous line, our line, and the line after.
			//
			// Originally this showed the prev line as well, but that is honestly too confusing.
			nextLine := ""
			if i+1 < len(lines) {
				nextLine = lines[i+1]
			}
			nextNextLine := ""
			if i+2 < len(lines) {
				nextNextLine = lines[i+2]
			}
			if len(line) <= 120 && len(nextLine) <= 120 && len(nextNextLine) <= 120 {
				Oprintf("%s:%d:%.80s @@ %.80s @@ %.80s\n", path, linum, line, nextLine, nextNextLine)
				continue
			}

			// hardcode 80 character limit for line.
			Oprintf("%s:%d:%.80s\n", path, linum, line)
		}
	}
	return nil
}

func Oprintf(content string, args ...any) {
	out := fmt.Sprintf(content, args...)
	// defang all ANSI escape sequences.
	out = Sanitize(out)
	_, _ = fmt.Fprint(os.Stdout, out)
}

func Eprintf(content string, args ...any) {
	out := fmt.Sprintf(content, args...)
	// defang all ANSI escape sequences.
	out = Sanitize(out)
	_, _ = fmt.Fprint(os.Stderr, out)
}

// Sanitize makes ansi escape sequences inert in a really naive way:
// by replacing \x1b with ? .
func Sanitize(content string) string {
	out := []byte(content)
	for i, ch := range out {
		switch int(ch) {
		// Escape \x1b == 27 is definitely part of escape sequences for terminals.
		// But on my machine specifically, another elusive byte sometimes causes
		// us to switch to an alternate character set.
		//
		// Some characters will be removed from this terrible, ad hoc translation later.
		case 0:
			out[i] = '%'
		case 27:
			out[i] = '?'
		}
	}
	return string(out)
}

func Eprintln(args ...any) {
	Eprintf("%s\n", args...)
}

func Eprinterr(wrapmsg string, e error) {
	if e == nil {
		return
	}
	Eprintln(fmt.Errorf("%s: %w", wrapmsg, e))
}
