package field

import (
	"database/sql/driver"

	"gorm.io/gorm/clause"
)

// ScanValuer interface for Field
type ScanValuer interface {
	Scan(src interface{}) error   // sql.Scanner
	Value() (driver.Value, error) // driver.Valuer
}

// Field a standard field struct
type Field struct{ expr }

// Add ...
func (field Field) Add(value driver.Valuer) Expr {
	return field.add(value)
}

// Sub ...
func (field Field) Sub(value driver.Valuer) Expr {
	return field.sub(value)
}

// Mul ...
func (field Field) Mul(value driver.Valuer) Expr {
	return field.mul(value)
}

// Div ...
func (field Field) Div(value driver.Valuer) Expr {
	return field.div(value)
}

// Eq judge equal
func (field Field) Eq(value driver.Valuer) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq judge not equal
func (field Field) Neq(value driver.Valuer) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// In ...
func (field Field) In(values ...driver.Valuer) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// Gt ...
func (field Field) Gt(value driver.Valuer) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte ...
func (field Field) Gte(value driver.Valuer) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt ...
func (field Field) Lt(value driver.Valuer) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte ...
func (field Field) Lte(value driver.Valuer) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// Like ...
func (field Field) Like(value driver.Valuer) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// Value ...
func (field Field) Value(value driver.Valuer) AssignExpr {
	return field.value(value)
}

// Sum ...
func (field Field) Sum() Field {
	return Field{field.sum()}
}

// IfNull ...
func (field Field) IfNull(value driver.Valuer) Expr {
	return field.ifNull(value)
}

func (field Field) toSlice(values ...driver.Valuer) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}
