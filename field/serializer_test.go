package field

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	gormtests "gorm.io/gorm/utils/tests"
)

type serializedTestModel struct {
	ID         uint
	Photos     []*string      `gorm:"serializer:json"`
	Attributes map[string]int `gorm:"serializer:json"`
	Name       string
}

type testSerializerValuer struct {
	value interface{}
	err   error
}

func (v testSerializerValuer) Value(context.Context, *schema.Field, reflect.Value, interface{}) (interface{}, error) {
	return v.value, v.err
}

func newSerializedTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	testDB, err := gorm.Open(gormtests.DummyDialector{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := schema.Parse(&serializedTestModel{}, &sync.Map{}, testDB.NamingStrategy)
	if err != nil {
		t.Fatal(err)
	}
	testDB.Statement = &gorm.Statement{
		DB:           testDB,
		Schema:       parsed,
		ReflectValue: reflect.ValueOf(&serializedTestModel{}),
	}
	return testDB
}

func TestSerializedValueUsesSchemaSerializer(t *testing.T) {
	one, two := "1", "2"
	tests := []struct {
		name   string
		column string
		value  interface{}
		want   interface{}
	}{
		{name: "slice", column: "photos", value: []*string{&one, &two}, want: `["1","2"]`},
		{name: "empty slice", column: "photos", value: []*string{}, want: "[]"},
		{name: "nil slice", column: "photos", value: []*string(nil), want: nil},
		{name: "map", column: "attributes", value: map[string]int{"one": 1}, want: `{"one":1}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB := newSerializedTestDB(t)
			got := (serializedValue{Column: tt.column, Value: tt.value}).GormValue(context.Background(), testDB)
			if testDB.Error != nil {
				t.Fatalf("GormValue() error = %v", testDB.Error)
			}
			if got.SQL != "?" || !reflect.DeepEqual(got.Vars, []interface{}{tt.want}) {
				t.Fatalf("GormValue() = %#v, want SQL ? and vars %#v", got, []interface{}{tt.want})
			}
		})
	}
}

func TestSerializedValuePrefersValueSerializerAndRecordsErrors(t *testing.T) {
	t.Run("value serializer", func(t *testing.T) {
		testDB := newSerializedTestDB(t)
		value := testSerializerValuer{value: "custom"}
		got := (serializedValue{Column: "photos", Value: value}).GormValue(context.Background(), testDB)
		if testDB.Error != nil {
			t.Fatalf("GormValue() error = %v", testDB.Error)
		}
		if !reflect.DeepEqual(got.Vars, []interface{}{"custom"}) {
			t.Fatalf("GormValue() vars = %#v, want custom value", got.Vars)
		}
	})

	t.Run("serializer error", func(t *testing.T) {
		testDB := newSerializedTestDB(t)
		wantErr := errors.New("serialize")
		value := testSerializerValuer{err: wantErr}
		_ = (serializedValue{Column: "photos", Value: value}).GormValue(context.Background(), testDB)
		if !errors.Is(testDB.Error, wantErr) {
			t.Fatalf("GormValue() error = %v, want %v", testDB.Error, wantErr)
		}
	})
}

func TestSerializedValueReportsInvalidSchemaState(t *testing.T) {
	tests := []struct {
		name   string
		column string
		setup  func(*gorm.DB)
	}{
		{name: "missing schema", column: "photos", setup: func(db *gorm.DB) { db.Statement.Schema = nil }},
		{name: "missing field", column: "missing", setup: func(*gorm.DB) {}},
		{name: "missing serializer", column: "name", setup: func(*gorm.DB) {}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB := newSerializedTestDB(t)
			tt.setup(testDB)
			got := (serializedValue{Column: tt.column, Value: "value"}).GormValue(context.Background(), testDB)
			if testDB.Error == nil {
				t.Fatal("GormValue() error = nil")
			}
			if got.SQL != "?" || !reflect.DeepEqual(got.Vars, []interface{}{"value"}) {
				t.Fatalf("GormValue() fallback = %#v", got)
			}
		})
	}
}
