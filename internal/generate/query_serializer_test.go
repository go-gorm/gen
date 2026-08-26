package generate

import (
	"context"
	"reflect"
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	gormtests "gorm.io/gorm/utils/tests"
)

type legacySerializerValue struct{}

func (*legacySerializerValue) Scan(context.Context, *schema.Field, reflect.Value, interface{}) error {
	return nil
}

func (*legacySerializerValue) Value(context.Context, *schema.Field, reflect.Value, interface{}) (interface{}, error) {
	return "legacy", nil
}

type serializerQueryModel struct {
	Legacy legacySerializerValue
	Photos []*string `gorm:"serializer:json"`
}

func TestParseStructDistinguishesSerializerSources(t *testing.T) {
	db, err := gorm.Open(gormtests.DummyDialector{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	meta := &QueryStructMeta{db: db}
	if err := meta.parseStruct(serializerQueryModel{}); err != nil {
		t.Fatal(err)
	}

	genTypes := make(map[string]string, len(meta.Fields))
	for _, modelField := range meta.Fields {
		genTypes[modelField.Name] = modelField.GenType()
	}
	if got := genTypes["Legacy"]; got != "Serializer" {
		t.Fatalf("legacy serializer GenType() = %q, want Serializer", got)
	}
	if got := genTypes["Photos"]; got != "SerializerField[[]*string]" {
		t.Fatalf("tagged field GenType() = %q, want SerializerField[[]*string]", got)
	}
}
