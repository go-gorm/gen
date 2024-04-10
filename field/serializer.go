package field

import (
	"context"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type valuerType struct {
	Column  string
	Value   schema.SerializerValuerInterface
	FucName string
}

func (v valuerType) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) {
	stmt := db.Statement.Schema
	field := stmt.LookUpField(v.Column)
	newValue, err := v.Value.Value(context.WithValue(ctx, "func_name", v.FucName), field, reflect.ValueOf(v.Value), v.Value)
	db.AddError(err)
	return clause.Expr{SQL: "?", Vars: []interface{}{newValue}}
}

// Serializer a standard field struct
type Serializer struct{ expr }

// Eq judge equal
func (s Serializer) Eq(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Eq{Column: s.RawExpr(), Value: valuerType{Column: s.ColumnName().String(), Value: value, FucName: "Eq"}}}
}

// Neq judge not equal
func (s Serializer) Neq(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Neq{Column: s.RawExpr(), Value: valuerType{Column: s.ColumnName().String(), Value: value, FucName: "Neq"}}}
}

// In ...
func (s Serializer) In(values ...schema.SerializerValuerInterface) Expr {
	return expr{e: clause.IN{Column: s.RawExpr(), Values: s.toSlice(values...)}}
}

// Gt ...
func (s Serializer) Gt(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Gt{Column: s.RawExpr(), Value: valuerType{Column: s.ColumnName().String(), Value: value, FucName: "Gt"}}}
}

// Gte ...
func (s Serializer) Gte(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Gte{Column: s.RawExpr(), Value: valuerType{Column: s.ColumnName().String(), Value: value, FucName: "Gte"}}}
}

// Lt ...
func (s Serializer) Lt(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Lt{Column: s.RawExpr(), Value: valuerType{Column: s.ColumnName().String(), Value: value, FucName: "Lt"}}}
}

// Lte ...
func (s Serializer) Lte(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Lte{Column: s.RawExpr(), Value: valuerType{Column: s.ColumnName().String(), Value: value, FucName: "Lte"}}}
}

// Like ...
func (s Serializer) Like(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Like{Column: s.RawExpr(), Value: valuerType{Column: s.ColumnName().String(), Value: value, FucName: "Like"}}}
}

// Value ...
func (s Serializer) Value(value schema.SerializerValuerInterface) AssignExpr {
	return s.value(valuerType{Column: s.ColumnName().String(), Value: value, FucName: "Value"})
}

// Sum ...
func (s Serializer) Sum() Number[float64] {
	return newNumber[float64](s.sum())
}

// IfNull ...
func (s Serializer) IfNull(value schema.SerializerValuerInterface) Expr {
	return s.ifNull(valuerType{Column: s.ColumnName().String(), Value: value, FucName: "IfNull"})
}

func (s Serializer) toSlice(values ...schema.SerializerValuerInterface) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = valuerType{Column: s.ColumnName().String(), Value: v, FucName: "In"}
	}
	return slice
}
