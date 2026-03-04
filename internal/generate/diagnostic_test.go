package generate

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"gorm.io/gen/internal/diagnostic"
	"gorm.io/gen/internal/parser"
)

func parseInterfaceSet(t *testing.T, src string) *parser.InterfaceSet {
	t.Helper()

	dir := t.TempDir()
	file := filepath.Join(dir, "diy.go")
	if err := os.WriteFile(file, []byte(src), 0640); err != nil {
		t.Fatalf("write file: %v", err)
	}

	var set parser.InterfaceSet
	paths := []*parser.InterfacePath{{Name: "UserMethods", FullName: "dal.UserMethods", Files: []string{file}}}
	if err := set.ParseFile(paths, []string{"User"}); err != nil {
		t.Fatalf("parse file: %v", err)
	}
	return &set
}

func testMeta() *QueryStructMeta {
	return &QueryStructMeta{
		ModelStructName: "User",
		QueryStructName: "user",
		S:               "u",
		TableName:       "users",
		StructInfo:      parser.Param{Package: "model", Type: "User"},
	}
}

func TestBuildDIYMethod_ReturnsDiagnosticOnIncompleteSQL(t *testing.T) {
	src := `package dal

import "gorm.io/gen"

type UserMethods interface {
	// FindByID
	//
	// SELECT * FROM users {{where
	FindByID(id int) gen.T
}
`
	set := parseInterfaceSet(t, src)

	_, err := BuildDIYMethod(set, testMeta(), nil)
	if err == nil {
		t.Fatalf("expected error")
	}
	var de *diagnostic.Error
	if !errors.As(err, &de) {
		t.Fatalf("expected diagnostic error, got %T: %v", err, err)
	}
	if de.Diag.Code != "SQL_INCOMPLETE" {
		t.Fatalf("unexpected code: %s", de.Diag.Code)
	}
	if de.Diag.File == "" || de.Diag.Line == 0 {
		t.Fatalf("expected file/line, got file=%q line=%d", de.Diag.File, de.Diag.Line)
	}
	if de.Diag.Interface != "UserMethods" || de.Diag.Method != "FindByID" {
		t.Fatalf("unexpected iface/method: %s.%s", de.Diag.Interface, de.Diag.Method)
	}
}

func TestBuildDIYMethod_DoesNotPanicOnTrailingBackslash(t *testing.T) {
	src := `package dal

import "gorm.io/gen"

type UserMethods interface {
	// FindByName
	//
	// SELECT * FROM users \
	FindByName() gen.T
}
`
	set := parseInterfaceSet(t, src)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	_, _ = BuildDIYMethod(set, testMeta(), nil)
}

func TestBuildDIYMethod_ReturnsDiagnosticOnTemplateParseError(t *testing.T) {
	src := `package dal

import "gorm.io/gen"

type UserMethods interface {
	// FindByName
	//
	// SELECT * FROM users {{bad}}
	FindByName() gen.T
}
`
	set := parseInterfaceSet(t, src)

	_, err := BuildDIYMethod(set, testMeta(), nil)
	if err == nil {
		t.Fatalf("expected error")
	}
	var de *diagnostic.Error
	if !errors.As(err, &de) {
		t.Fatalf("expected diagnostic error, got %T: %v", err, err)
	}
	if de.Diag.Code != "TEMPLATE_PARSE" {
		t.Fatalf("unexpected code: %s", de.Diag.Code)
	}
	if de.Diag.Interface != "UserMethods" || de.Diag.Method != "FindByName" {
		t.Fatalf("unexpected iface/method: %s.%s", de.Diag.Interface, de.Diag.Method)
	}
	if de.Err == nil {
		t.Fatalf("expected cause error")
	}
}
