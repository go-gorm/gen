package generate

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"gorm.io/gen/internal/diagnostic"
	"gorm.io/gen/internal/parser"
)

func TestBuildDIYMethod_ReturnsDiagnosticOnIncompleteSQL(t *testing.T) {
	dir := t.TempDir()
	src := `package dal

import "gorm.io/gen"

type UserMethods interface {
	// FindByID
	//
	// SELECT * FROM users {{where
	FindByID(id int) gen.T
}
`
	file := filepath.Join(dir, "diy.go")
	if err := os.WriteFile(file, []byte(src), 0640); err != nil {
		t.Fatalf("write file: %v", err)
	}

	var set parser.InterfaceSet
	paths := []*parser.InterfacePath{{Name: "UserMethods", FullName: "dal.UserMethods", Files: []string{file}}}
	if err := set.ParseFile(paths, []string{"User"}); err != nil {
		t.Fatalf("parse file: %v", err)
	}

	meta := &QueryStructMeta{
		ModelStructName: "User",
		QueryStructName: "user",
		S:               "u",
		TableName:       "users",
		StructInfo:      parser.Param{Package: "model", Type: "User"},
	}

	_, err := BuildDIYMethod(&set, meta, nil)
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

