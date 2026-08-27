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

type destinationSerializer struct {
	destination reflect.Value
}

func (*destinationSerializer) Scan(context.Context, *schema.Field, reflect.Value, interface{}) error {
	return nil
}

func (s *destinationSerializer) Value(_ context.Context, _ *schema.Field, destination reflect.Value, _ interface{}) (interface{}, error) {
	s.destination = destination
	return "serialized", nil
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

func TestSchemaSerializerValueUsesSchemaSerializer(t *testing.T) {
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
			got := (schemaSerializerValue{Column: tt.column, Value: tt.value}).GormValue(context.Background(), testDB)
			if testDB.Error != nil {
				t.Fatalf("GormValue() error = %v", testDB.Error)
			}
			if got.SQL != "?" || !reflect.DeepEqual(got.Vars, []interface{}{tt.want}) {
				t.Fatalf("GormValue() = %#v, want SQL ? and vars %#v", got, []interface{}{tt.want})
			}
		})
	}
}

func TestSerializerFieldPreservesValueSerializerBehavior(t *testing.T) {
	t.Run("value serializer", func(t *testing.T) {
		testDB := newSerializedTestDB(t)
		value := testSerializerValuer{value: "custom"}
		wrapped, ok := NewSerializer("", "photos").wrap(value).(ValuerType)
		if !ok {
			t.Fatalf("wrap() type = %T, want ValuerType", wrapped)
		}
		got := wrapped.GormValue(context.Background(), testDB)
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
		wrapped, ok := NewSerializer("", "photos").wrap(value).(ValuerType)
		if !ok {
			t.Fatalf("wrap() type = %T, want ValuerType", wrapped)
		}
		_ = wrapped.GormValue(context.Background(), testDB)
		if !errors.Is(testDB.Error, wantErr) {
			t.Fatalf("GormValue() error = %v, want %v", testDB.Error, wantErr)
		}
	})
}

func TestSchemaSerializerValueReportsInvalidSchemaState(t *testing.T) {
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
			got := (schemaSerializerValue{Column: tt.column, Value: "value"}).GormValue(context.Background(), testDB)
			if testDB.Error == nil {
				t.Fatal("GormValue() error = nil")
			}
			if got.SQL != "?" || !reflect.DeepEqual(got.Vars, []interface{}{"value"}) {
				t.Fatalf("GormValue() fallback = %#v", got)
			}
		})
	}
}

func TestSchemaSerializerValueSuppliesValidDestination(t *testing.T) {
	testDB := newSerializedTestDB(t)
	testDB.Statement.ReflectValue = reflect.Value{}
	serializer := new(destinationSerializer)
	testDB.Statement.Schema.LookUpField("photos").Serializer = serializer

	got := (schemaSerializerValue{Column: "photos", Value: nil}).GormValue(context.Background(), testDB)
	if testDB.Error != nil {
		t.Fatalf("GormValue() error = %v", testDB.Error)
	}
	if !serializer.destination.IsValid() {
		t.Fatal("serializer destination is invalid")
	}
	wantType := reflect.New(testDB.Statement.Schema.ModelType).Type()
	if serializer.destination.Type() != wantType {
		t.Fatalf("serializer destination type = %v, want %v", serializer.destination.Type(), wantType)
	}
	if !reflect.DeepEqual(got.Vars, []interface{}{"serialized"}) {
		t.Fatalf("GormValue() vars = %#v, want serialized", got.Vars)
	}
}
