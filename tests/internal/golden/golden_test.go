package golden

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompareDirs(t *testing.T) {
	tests := []struct {
		name      string
		wantFiles map[string]string
		gotFiles  map[string]string
		wantError string
	}{
		{
			name: "identical nested trees and empty files",
			wantFiles: map[string]string{
				"empty.go":       "",
				"nested/user.go": "package model\n",
			},
			gotFiles: map[string]string{
				"empty.go":       "",
				"nested/user.go": "package model\n",
			},
		},
		{
			name:      "content differs",
			wantFiles: map[string]string{"query.go": "line one\nexpected\n"},
			gotFiles:  map[string]string{"query.go": "line one\nactual\n"},
			wantError: "query.go line 2",
		},
		{
			name:      "expected file is missing",
			wantFiles: map[string]string{"model.go": "package model\n"},
			gotFiles:  map[string]string{},
			wantError: "missing path from actual directory: model.go",
		},
		{
			name:      "actual tree has an extra file",
			wantFiles: map[string]string{},
			gotFiles:  map[string]string{"extra.go": "package query\n"},
			wantError: "unexpected path in actual directory: extra.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantDir := t.TempDir()
			gotDir := t.TempDir()
			writeTree(t, wantDir, tt.wantFiles)
			writeTree(t, gotDir, tt.gotFiles)

			err := CompareDirs(wantDir, gotDir)
			if tt.wantError == "" {
				if err != nil {
					t.Fatalf("CompareDirs() error = %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("CompareDirs() error = %v, want substring %q", err, tt.wantError)
			}
		})
	}
}

func TestCompareDirsReportsFileDirectoryMismatch(t *testing.T) {
	wantDir := t.TempDir()
	gotDir := t.TempDir()
	writeTree(t, wantDir, map[string]string{"item": "file"})
	if err := os.Mkdir(filepath.Join(gotDir, "item"), 0750); err != nil {
		t.Fatalf("create directory: %v", err)
	}

	err := CompareDirs(wantDir, gotDir)
	if err == nil || !strings.Contains(err.Error(), "path type differs at item") {
		t.Fatalf("CompareDirs() error = %v, want type mismatch", err)
	}
}

func writeTree(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for relative, content := range files {
		path := filepath.Join(root, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
			t.Fatalf("create parent for %s: %v", relative, err)
		}
		if err := os.WriteFile(path, []byte(content), 0640); err != nil {
			t.Fatalf("write %s: %v", relative, err)
		}
	}
}
