package parser_test

import (
	"path/filepath"
	"testing"

	"gorm.io/gen/internal/parser"
	"gorm.io/gen/internal/parser/fixtures/ifaces"
)

func TestGetInterfacePath(t *testing.T) {
	paths, err := parser.GetInterfacePath(func(testIF ifaces.TestIF, m ifaces.InsertMethod) {})
	if err != nil {
		t.Fatalf("GetInterfacePath error: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(paths))
	}
	for _, p := range paths {
		if len(p.Files) == 0 {
			t.Fatalf("expected files for %s, got 0", p.FullName)
		}
		found := false
		for _, f := range p.Files {
			if filepath.Base(f) == "ifaces.go" {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected to find ifaces.go in files for %s", p.FullName)
		}
	}
}
