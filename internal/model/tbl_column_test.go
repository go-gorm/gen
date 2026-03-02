package model

import (
	"database/sql"
	"reflect"
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
	comment       string
}

func (t testColumnType) Name() string                      { return t.name }
func (t testColumnType) DatabaseTypeName() string          { return t.databaseType }
func (t testColumnType) ColumnType() (string, bool)        { return t.columnType, t.columnType != "" }
func (t testColumnType) PrimaryKey() (bool, bool)          { return t.primaryKey, true }
func (t testColumnType) AutoIncrement() (bool, bool)       { return t.autoIncrement, true }
func (t testColumnType) Nullable() (bool, bool)            { return t.nullable, true }
func (t testColumnType) Unique() (bool, bool)              { return t.unique, true }
func (t testColumnType) DefaultValue() (string, bool)      { return t.defaultValue, t.defaultValue != "" }
func (t testColumnType) Comment() (string, bool)           { return t.comment, t.comment != "" }
func (t testColumnType) ScanType() reflect.Type            { return nil }
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
