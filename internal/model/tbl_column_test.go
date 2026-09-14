package model

import (
	"bytes"
	"database/sql"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
	"testing"
	"text/template"

	"gorm.io/gen/field"
	gentemplate "gorm.io/gen/internal/template"
	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

type testColumnType struct {
	name          string
	databaseType  string
	columnType    string
	primaryKey    bool
	autoIncrement bool
	nullable      bool
	unique        bool
	defaultValue  string
	defaultValid  bool
	comment       string
	commentValid  bool
	scanType      reflect.Type
}

func (t testColumnType) Name() string                { return t.name }
func (t testColumnType) DatabaseTypeName() string    { return t.databaseType }
func (t testColumnType) ColumnType() (string, bool)  { return t.columnType, t.columnType != "" }
func (t testColumnType) PrimaryKey() (bool, bool)    { return t.primaryKey, true }
func (t testColumnType) AutoIncrement() (bool, bool) { return t.autoIncrement, true }
func (t testColumnType) Nullable() (bool, bool)      { return t.nullable, true }
func (t testColumnType) Unique() (bool, bool)        { return t.unique, true }
func (t testColumnType) DefaultValue() (string, bool) {
	return t.defaultValue, t.defaultValid || t.defaultValue != ""
}
func (t testColumnType) Comment() (string, bool)           { return t.comment, t.commentValid || t.comment != "" }
func (t testColumnType) ScanType() reflect.Type            { return t.scanType }
func (t testColumnType) Length() (int64, bool)             { return 0, false }
func (t testColumnType) DecimalSize() (int64, int64, bool) { return 0, 0, false }

var _ gorm.ColumnType = testColumnType{}

func TestBuildGormTagIndexOrderDeterministic(t *testing.T) {
	col := Column{
		ColumnType: testColumnType{
			name:         "tenant_id",
			databaseType: "uuid",
			columnType:   "uuid",
		},
		Indexes: []*Index{
			{
				Index: migrator.Index{
					NameValue:   "idx_payout_tenant_payment_at",
					ColumnList:  []string{"tenant_id"},
					UniqueValue: sql.NullBool{Bool: false, Valid: true},
				},
				Priority: 2,
			},
			{
				Index: migrator.Index{
					NameValue:   "idx_payout_tenant_id",
					ColumnList:  []string{"tenant_id"},
					UniqueValue: sql.NullBool{Bool: false, Valid: true},
				},
				Priority: 1,
			},
		},
	}

	tag := col.buildGormTag(false)
	got := tag[field.TagKeyGormIndex]
	want := []string{
		"idx_payout_tenant_id,priority:1",
		"idx_payout_tenant_payment_at,priority:2",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected index order: got=%v want=%v", got, want)
	}
}

func TestColumnGetDataType(t *testing.T) {
	tests := []struct {
		name   string
		column Column
		want   string
	}{
		{
			name: "custom mapping",
			column: Column{
				ColumnType: testColumnType{databaseType: "uuid"},
				dataTypeMap: map[string]func(gorm.ColumnType) string{
					"uuid": func(gorm.ColumnType) string { return "custom.UUID" },
				},
			},
			want: "custom.UUID",
		},
		{
			name: "scan type",
			column: Column{
				ColumnType:  testColumnType{databaseType: "varchar", scanType: reflect.TypeOf(sql.NullString{})},
				UseScanType: true,
			},
			want: "sql.NullString",
		},
		{
			name:   "tinyint bool",
			column: Column{ColumnType: testColumnType{databaseType: "tinyint", columnType: "tinyint(1)"}},
			want:   "bool",
		},
		{
			name:   "unknown defaults to string",
			column: Column{ColumnType: testColumnType{databaseType: "uuid"}},
			want:   "string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.column.GetDataType(); got != tt.want {
				t.Fatalf("GetDataType() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestColumnToFieldTypeRules(t *testing.T) {
	tests := []struct {
		name       string
		columnType testColumnType
		nullable   bool
		coverable  bool
		signable   bool
		wantType   string
	}{
		{
			name:       "nullable",
			columnType: testColumnType{name: "name", databaseType: "varchar", nullable: true},
			nullable:   true,
			wantType:   "*string",
		},
		{
			name:       "coverable default",
			columnType: testColumnType{name: "age", databaseType: "int", defaultValue: "1"},
			coverable:  true,
			wantType:   "*int32",
		},
		{
			name:       "unsigned",
			columnType: testColumnType{name: "id", databaseType: "bigint", columnType: "bigint unsigned"},
			signable:   true,
			wantType:   "uint64",
		},
		{
			name:       "soft delete",
			columnType: testColumnType{name: "deleted_at", databaseType: "datetime"},
			wantType:   "gorm.DeletedAt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			column := Column{ColumnType: tt.columnType}
			column.WithNS(nil)
			got := column.ToField(tt.nullable, tt.coverable, tt.signable, false)
			if got.Type != tt.wantType {
				t.Fatalf("ToField().Type = %q, want %q", got.Type, tt.wantType)
			}
			if got.Tag[field.TagKeyJson] != tt.columnType.name {
				t.Fatalf("json tag = %q, want %q", got.Tag[field.TagKeyJson], tt.columnType.name)
			}
		})
	}
}

func TestColumnDefaultAndCommentSanitization(t *testing.T) {
	blankDefault := Column{ColumnType: testColumnType{name: "title", defaultValue: "   ", defaultValid: true}}
	if got, want := blankDefault.defaultTagValue(), "'   '"; got != want {
		t.Fatalf("defaultTagValue() = %q, want %q", got, want)
	}

	column := Column{ColumnType: testColumnType{
		name:         "note",
		databaseType: "text",
		comment:      "line1\n*/ line2: value; `quoted` \"text\"",
	}}
	column.WithNS(func(name string) string { return "json_" + name })
	got := column.ToField(false, false, false, false)
	if !got.MultilineComment {
		t.Fatal("expected multiline comment")
	}
	if strings.Contains(got.ColumnComment, "*/") {
		t.Fatalf("column comment still closes comment: %q", got.ColumnComment)
	}
	if got.Tag[field.TagKeyJson] != "json_note" {
		t.Fatalf("json tag = %q, want json_note", got.Tag[field.TagKeyJson])
	}
	commentTag := got.GORMTag[field.TagKeyGormComment]
	if len(commentTag) != 1 || strings.ContainsAny(commentTag[0], "\n`;:") {
		t.Fatalf("comment tag was not sanitized: %v", commentTag)
	}
}

// TestColumnTagValueEscaping verifies that single-line comments and default
// values are escaped before entering the gorm tag: a double quote used to
// truncate the quoted tag value at the reflect.StructTag level and silently
// swallow the following directives, a backtick terminated the raw-string
// literal of the generated struct tag, and a semicolon started a bogus gorm
// directive.
func TestColumnTagValueEscaping(t *testing.T) {
	tests := []struct {
		name         string
		comment      string
		defaultValue string
		withDefault  bool
		wantComment  string
		wantDefault  string
	}{
		{
			name:        "single-line double quote",
			comment:     `say "hi" ok`,
			wantComment: `say "hi" ok`,
		},
		{
			name:        "single-line backtick",
			comment:     "has `tick` inside",
			wantComment: "has 'tick' inside",
		},
		{
			name:        "single-line semicolon",
			comment:     "a;b",
			wantComment: "a b",
		},
		{
			name:        "single-line colon preserved",
			comment:     "start: 10",
			wantComment: "start: 10",
		},
		{
			name:         "default value with double quote",
			defaultValue: `"none"`,
			withDefault:  true,
			wantDefault:  `"none"`,
		},
		{
			name:         "default value with backtick",
			defaultValue: "`v`",
			withDefault:  true,
			wantDefault:  "'v'",
		},
		{
			name:         "comment and default combined",
			comment:      `x";y`,
			defaultValue: `"z"`,
			withDefault:  true,
			wantComment:  `x" y`,
			wantDefault:  `"z"`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			column := Column{ColumnType: testColumnType{
				name:         "note",
				databaseType: "text",
				comment:      tt.comment,
				defaultValue: tt.defaultValue,
				defaultValid: true,
			}}
			column.WithNS(nil)
			got := column.ToField(false, false, false, tt.withDefault)

			tag := got.Tags()
			value, ok := reflect.StructTag(tag).Lookup("gorm")
			if !ok {
				t.Fatalf("gorm tag lookup failed for %s", tag)
			}
			settings := schema.ParseTagSetting(value, ";")
			if settings["COLUMN"] != "note" {
				t.Errorf("COLUMN = %q, want %q (raw gorm value %q)", settings["COLUMN"], "note", value)
			}
			if tt.wantComment != "" && settings["COMMENT"] != tt.wantComment {
				t.Errorf("COMMENT = %q, want %q (raw gorm value %q)", settings["COMMENT"], tt.wantComment, value)
			}
			if tt.wantDefault != "" && settings["DEFAULT"] != tt.wantDefault {
				t.Errorf("DEFAULT = %q, want %q (raw gorm value %q)", settings["DEFAULT"], tt.wantDefault, value)
			}
		})
	}
}

// TestColumnTagEscapingRendersCompilableModel renders the model template with
// a field whose comment contains a backtick and a double quote and verifies
// the generated file parses as Go source.
func TestColumnTagEscapingRendersCompilableModel(t *testing.T) {
	column := Column{ColumnType: testColumnType{
		name:         "note",
		databaseType: "text",
		comment:      "has `tick` and \"quote\"",
	}}
	column.WithNS(nil)
	field := column.ToField(false, false, false, false)

	tmpl, err := template.New("model").Parse(gentemplate.Model)
	if err != nil {
		t.Fatalf("parse model template: %v", err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"StructInfo":            map[string]string{"Package": "model", "Type": "Note"},
		"ModelStructName":       "Note",
		"TableName":             "notes",
		"MultilineTableComment": false,
		"StructComment":         "",
		"ImportPkgPaths":        []string{},
		"Fields":                []*Field{field},
	})
	if err != nil {
		t.Fatalf("execute model template: %v", err)
	}

	if _, err := parser.ParseFile(token.NewFileSet(), "note.gen.go", buf.String(), 0); err != nil {
		t.Fatalf("generated model does not parse: %v\nsource:\n%s", err, buf.String())
	}
	if strings.ContainsRune(field.Tags(), '`') {
		t.Fatalf("tag still contains a raw backtick: %s", field.Tags())
	}
}
