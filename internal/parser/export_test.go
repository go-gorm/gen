package parser_test

import (
	"path/filepath"
	"testing"

	"gorm.io/gen/field"
	"gorm.io/gen/internal/parser"
)

func TestGetInterfacePathImportedPackage(t *testing.T) {
	paths, err := parser.GetInterfacePath(func(expr field.Expr, order field.OrderExpr) {})
	if err != nil {
		t.Fatalf("GetInterfacePath error: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("expected 2 interface paths, got %d", len(paths))
	}

	for _, path := range paths {
		if path.Package != "gorm.io/gen/field" {
			t.Fatalf("unexpected package for %s: %q", path.FullName, path.Package)
		}
		if !containsBaseFile(path.Files, "expr.go") {
			t.Fatalf("expected expr.go in files for %s, got %#v", path.FullName, path.Files)
		}
	}
}

func containsBaseFile(files []string, name string) bool {
	for _, file := range files {
		if filepath.Base(file) == name {
			return true
		}
	}
	return false
}
