package generate

import (
	"errors"
	"strings"
	"testing"

	"gorm.io/gen/internal/diagnostic"
)

func TestParseDocString_WrapperAndQuoteValidation(t *testing.T) {
	cases := []struct {
		name    string
		doc     string
		wantSQL string
		wantErr bool
	}{
		{name: "sql wrapper", doc: "sql(select 1)", wantSQL: "select 1"},
		{name: "sql wrapper uppercase", doc: "SQL(select 1)", wantSQL: "select 1"},
		{name: "where wrapper", doc: "where(age > 1)", wantSQL: "age > 1"},
		{name: "quoted sql", doc: `"select 1"`, wantSQL: "select 1"},
		{name: "plain sql", doc: "select 1", wantSQL: "select 1"},
		{name: "empty doc", doc: "", wantSQL: ""},
		{name: "sql open paren only", doc: "sql(", wantErr: true},
		{name: "where open paren only", doc: "where(", wantErr: true},
		{name: "sql missing close paren", doc: "sql(select 1", wantErr: true},
		{name: "where missing close paren", doc: "where(age > 1", wantErr: true},
		{name: "lone double quote", doc: `"`, wantSQL: `"`},
		{name: "empty quoted sql", doc: `""`, wantSQL: ""},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			m := &InterfaceMethod{
				MethodName:    "FindByID",
				Doc:           tt.doc,
				File:          "diy.go",
				DocLine:       3,
				DocColumn:     4,
				InterfaceName: "UserMethods",
			}
			sql, _, _, err := m.parseDocString()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got sql=%q", sql)
				}
				var de *diagnostic.Error
				if !errors.As(err, &de) || de.Diag.Code != diagnostic.CodeSQLIncomplete {
					t.Fatalf("expected %s diagnostic, got %v", diagnostic.CodeSQLIncomplete, err)
				}
				if !strings.Contains(err.Error(), tt.doc) {
					t.Fatalf("error should quote the raw comment: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if sql != tt.wantSQL {
				t.Fatalf("sql = %q, want %q", sql, tt.wantSQL)
			}
		})
	}
}

func TestCheckSQL_DegenerateDocCommentsFailCleanly(t *testing.T) {
	cases := []struct {
		name     string
		doc      string
		wantCode string
	}{
		{name: "sql open paren only", doc: "sql(", wantCode: diagnostic.CodeSQLIncomplete},
		{name: "where open paren only", doc: "where(", wantCode: diagnostic.CodeSQLIncomplete},
		{name: "sql missing close paren", doc: "sql(select 1", wantCode: diagnostic.CodeSQLIncomplete},
		{name: "lone double quote", doc: `"`, wantCode: diagnostic.CodeSQLIncomplete},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			m := &InterfaceMethod{
				MethodName:    "FindByID",
				Doc:           tt.doc,
				File:          "diy.go",
				DocLine:       3,
				DocColumn:     4,
				InterfaceName: "UserMethods",
			}
			err := m.checkSQL()
			if err == nil {
				t.Fatalf("expected error, got none")
			}
			var de *diagnostic.Error
			if !errors.As(err, &de) {
				t.Fatalf("expected diagnostic error, got %T: %v", err, err)
			}
			if de.Diag.Code != tt.wantCode {
				t.Fatalf("code = %s, want %s", de.Diag.Code, tt.wantCode)
			}
			if de.Diag.Interface != "UserMethods" || de.Diag.Method != "FindByID" {
				t.Fatalf("unexpected iface/method: %s.%s", de.Diag.Interface, de.Diag.Method)
			}
			if de.Diag.File != "diy.go" || de.Diag.Line != 3 {
				t.Fatalf("unexpected location: file=%q line=%d", de.Diag.File, de.Diag.Line)
			}
		})
	}
}

func TestBuildDIYMethod_RejectsUnterminatedDocWrapper(t *testing.T) {
	src := `package dal

import "gorm.io/gen"

type UserMethods interface {
	// sql(
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
	if de.Diag.Code != diagnostic.CodeSQLIncomplete {
		t.Fatalf("unexpected code: %s", de.Diag.Code)
	}
	if de.Diag.Interface != "UserMethods" || de.Diag.Method != "FindByID" {
		t.Fatalf("unexpected iface/method: %s.%s", de.Diag.Interface, de.Diag.Method)
	}
	if !strings.Contains(de.Diag.Message, "sql(") {
		t.Fatalf("message should quote the raw comment: %s", de.Diag.Message)
	}
}

func TestBuildDIYMethod_AcceptsWrappedDocComments(t *testing.T) {
	src := `package dal

import "gorm.io/gen"

type UserMethods interface {
	// sql(select * from users where id = @id)
	FindByID(id int) gen.T

	// where(age > @age)
	FindByAge(age int) gen.T
}
`
	set := parseInterfaceSet(t, src)

	methods, err := BuildDIYMethod(set, testMeta(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	byName := map[string]*InterfaceMethod{}
	for _, m := range methods {
		byName[m.MethodName] = m
	}
	if m := byName["FindByID"]; m == nil || m.SQLString != "select * from users where id = @id" || m.GormOption != "Raw" {
		t.Fatalf("unexpected FindByID result: %+v", m)
	}
	if m := byName["FindByAge"]; m == nil || m.SQLString != "age > @age" || m.GormOption != "Where" {
		t.Fatalf("unexpected FindByAge result: %+v", m)
	}
}
