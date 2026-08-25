package model

import (
	"database/sql"
	"reflect"
	"strings"
	"testing"

	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
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
