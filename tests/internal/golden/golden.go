// Package golden compares generated directory trees without external tools.
package golden

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type entry struct {
	path  string
	isDir bool
}

// CompareDirs compares two directory trees, including their complete file sets.
func CompareDirs(wantDir, gotDir string) error {
	want, err := collect(wantDir)
	if err != nil {
		return fmt.Errorf("read expected directory: %w", err)
	}
	got, err := collect(gotDir)
	if err != nil {
		return fmt.Errorf("read actual directory: %w", err)
	}

	paths := make([]string, 0, len(want)+len(got))
	seen := make(map[string]struct{}, len(want)+len(got))
	for path := range want {
		seen[path] = struct{}{}
		paths = append(paths, path)
	}
	for path := range got {
		if _, ok := seen[path]; ok {
			continue
		}
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		wantEntry, inWant := want[path]
		gotEntry, inGot := got[path]
		switch {
		case !inWant:
			return fmt.Errorf("unexpected path in actual directory: %s", path)
		case !inGot:
			return fmt.Errorf("missing path from actual directory: %s", path)
		case wantEntry.isDir != gotEntry.isDir:
			return fmt.Errorf("path type differs at %s: expected %s, got %s", path, entryType(wantEntry), entryType(gotEntry))
		case wantEntry.isDir:
			continue
		}

		wantContent, err := os.ReadFile(filepath.Join(wantDir, filepath.FromSlash(path)))
		if err != nil {
			return fmt.Errorf("read expected file %s: %w", path, err)
		}
		gotContent, err := os.ReadFile(filepath.Join(gotDir, filepath.FromSlash(path)))
		if err != nil {
			return fmt.Errorf("read actual file %s: %w", path, err)
		}
		if !bytes.Equal(wantContent, gotContent) {
			line, wantLine, gotLine := firstDifferentLine(wantContent, gotContent)
			return fmt.Errorf("content differs at %s line %d: expected %q, got %q", path, line, wantLine, gotLine)
		}
	}
	return nil
}

func collect(root string) (map[string]entry, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", root)
	}

	entries := make(map[string]entry)
	err = filepath.WalkDir(root, func(path string, dirEntry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == root {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		entries[relative] = entry{path: relative, isDir: dirEntry.IsDir()}
		return nil
	})
	return entries, err
}

func entryType(value entry) string {
	if value.isDir {
		return "directory"
	}
	return "file"
}

func firstDifferentLine(want, got []byte) (line int, wantLine, gotLine string) {
	wantLines := strings.Split(string(want), "\n")
	gotLines := strings.Split(string(got), "\n")
	count := len(wantLines)
	if len(gotLines) > count {
		count = len(gotLines)
	}
	for index := 0; index < count; index++ {
		var expected, actual string
		if index < len(wantLines) {
			expected = wantLines[index]
		}
		if index < len(gotLines) {
			actual = gotLines[index]
		}
		if expected != actual {
			return index + 1, expected, actual
		}
	}
	return 1, "", ""
}
