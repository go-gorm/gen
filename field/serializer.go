package field

import (
	"context"
	"fmt"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// ValuerType delays conversion of a serializer value until GORM builds the statement.
type ValuerType struct {
	// Column identifies the schema field whose serializer should perform the conversion.
	Column string
	// Value is the application value passed to the schema serializer.
	Value schema.SerializerValuerInterface
}

// GormValue converts Value with the serializer registered for Column.
// Conversion errors are recorded on db and the returned expression still uses
// a bind variable so statement construction remains parameterized.
func (v ValuerType) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) {
	stmt := db.Statement.Schema
	field := stmt.LookUpField(v.Column)
	newValue, err := v.Value.Value(ctx, field, reflect.ValueOf(v.Value), v.Value)
	_ = db.AddError(err)
	return clause.Expr{SQL: "?", Vars: []interface{}{newValue}}
}

// Serializer a standard field struct
type Serializer struct{ expr }

// Eq judge equal
func (field Serializer) Eq(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: ValuerType{Column: field.ColumnName().String(), Value: value}}}
}

// Neq judge not equal
func (field Serializer) Neq(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: ValuerType{Column: field.ColumnName().String(), Value: value}}}
}

// In ...
func (field Serializer) In(values ...schema.SerializerValuerInterface) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// Gt ...
func (field Serializer) Gt(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: ValuerType{Column: field.ColumnName().String(), Value: value}}}
}

// Gte ...
func (field Serializer) Gte(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: ValuerType{Column: field.ColumnName().String(), Value: value}}}
}

// Lt ...
func (field Serializer) Lt(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: ValuerType{Column: field.ColumnName().String(), Value: value}}}
}

// Lte ...
func (field Serializer) Lte(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: ValuerType{Column: field.ColumnName().String(), Value: value}}}
}

// Like ...
func (field Serializer) Like(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: ValuerType{Column: field.ColumnName().String(), Value: value}}}
}

// Value ...
func (field Serializer) Value(value schema.SerializerValuerInterface) AssignExpr {
	return field.value(ValuerType{Column: field.ColumnName().String(), Value: value})
}

// Sum ...
func (field Serializer) Sum() Number[float64] {
	return newNumber[float64](field.sum())
}

// IfNull ...
func (field Serializer) IfNull(value schema.SerializerValuerInterface) Expr {
	return field.ifNull(ValuerType{Column: field.ColumnName().String(), Value: value})
}

func (field Serializer) toSlice(values ...schema.SerializerValuerInterface) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = ValuerType{Column: field.ColumnName().String(), Value: v}
	}
	return slice
}

// Serialized is a field whose values are converted by the serializer configured
// on the corresponding GORM schema field.
type Serialized struct{ expr }

// Eq compares the field with a value after schema serialization.
func (field Serialized) Eq(value interface{}) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: field.wrap(value)}}
}

// Neq compares the field with a value for inequality after schema serialization.
func (field Serialized) Neq(value interface{}) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: field.wrap(value)}}
}

// In compares the field with serialized candidate values.
func (field Serialized) In(values ...interface{}) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSerializedSlice(values...)}}
}

// Gt compares the field with a serialized value using greater-than semantics.
func (field Serialized) Gt(value interface{}) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: field.wrap(value)}}
}

// Gte compares the field with a serialized value using greater-than-or-equal semantics.
func (field Serialized) Gte(value interface{}) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: field.wrap(value)}}
}

// Lt compares the field with a serialized value using less-than semantics.
func (field Serialized) Lt(value interface{}) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: field.wrap(value)}}
}

// Lte compares the field with a serialized value using less-than-or-equal semantics.
func (field Serialized) Lte(value interface{}) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: field.wrap(value)}}
}

// Like compares the field with a serialized pattern value.
func (field Serialized) Like(value interface{}) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: field.wrap(value)}}
}

// Value creates an assignment whose value is converted by the schema serializer.
func (field Serialized) Value(value interface{}) AssignExpr {
	return field.value(field.wrap(value))
}

// Sum returns a numeric sum expression for the field.
func (field Serialized) Sum() Number[float64] {
	return newNumber[float64](field.sum())
}

// IfNull returns the serialized fallback value when the field is NULL.
func (field Serialized) IfNull(value interface{}) Expr {
	return field.ifNull(field.wrap(value))
}

func (field Serialized) wrap(value interface{}) serializedValue {
	return serializedValue{Column: field.ColumnName().String(), Value: value}
}

func (field Serialized) toSerializedSlice(values ...interface{}) []interface{} {
	slice := make([]interface{}, len(values))
	for i, value := range values {
		slice[i] = field.wrap(value)
	}
	return slice
}

type serializedValue struct {
	Column string
	Value  interface{}
}

func (v serializedValue) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	expr := clause.Expr{SQL: "?", Vars: []interface{}{v.Value}}
	if db == nil || db.Statement == nil || db.Statement.Schema == nil {
		if db != nil {
			_ = db.AddError(fmt.Errorf("field: schema is unavailable for serialized column %q", v.Column))
		}
		return expr
	}

	schemaField := db.Statement.Schema.LookUpField(v.Column)
	if schemaField == nil {
		_ = db.AddError(fmt.Errorf("field: serialized column %q was not found in schema", v.Column))
		return expr
	}

	var valuer schema.SerializerValuerInterface = schemaField.Serializer
	if valueValuer, ok := v.Value.(schema.SerializerValuerInterface); ok {
		valuer = valueValuer
	}
	if valuer == nil {
		_ = db.AddError(fmt.Errorf("field: column %q has no schema serializer", v.Column))
		return expr
	}

	destination := db.Statement.ReflectValue
	if !destination.IsValid() {
		destination = reflect.ValueOf(v.Value)
	}
	newValue, err := valuer.Value(ctx, schemaField, destination, v.Value)
	_ = db.AddError(err)
	expr.Vars[0] = newValue
	return expr
}
