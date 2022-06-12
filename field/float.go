package field

import "gorm.io/gorm/clause"

// Float64 float64 type field
type Float64 Field

// Eq equal to
func (field Float64) Eq(value float64) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq not equal to
func (field Float64) Neq(value float64) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// Gt greater than
func (field Float64) Gt(value float64) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte greater or equal to
func (field Float64) Gte(value float64) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt less than
func (field Float64) Lt(value float64) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte less or equal to
func (field Float64) Lte(value float64) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// In ...
func (field Float64) In(values ...float64) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// NotIn ...
func (field Float64) NotIn(values ...float64) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Between ...
func (field Float64) Between(left float64, right float64) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field Float64) NotBetween(left float64, right float64) Expr {
	return Not(field.Between(left, right))
}

// Like ...
func (field Float64) Like(value float64) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// NotLike ...
func (field Float64) NotLike(value float64) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

// Add ...
func (field Float64) Add(value float64) Float64 {
	return Float64{field.add(value)}
}

// Sub ...
func (field Float64) Sub(value float64) Float64 {
	return Float64{field.sub(value)}
}

// Mul ...
func (field Float64) Mul(value float64) Float64 {
	return Float64{field.mul(value)}
}

// Div ...
func (field Float64) Div(value float64) Float64 {
	return Float64{field.div(value)}
}

// FloorDiv ...
func (field Float64) FloorDiv(value float64) Int {
	return Int{field.floorDiv(value)}
}

// Floor ...
func (field Float64) Floor() Int {
	return Int{field.floor()}
}

// Value set value
func (field Float64) Value(value float64) AssignExpr {
	return field.value(value)
}

// Zero set zero value
func (field Float64) Zero() AssignExpr {
	return field.value(0)
}

// Sum ...
func (field Float64) Sum() Float64 {
	return Float64{field.sum()}
}

// IfNull ...
func (field Float64) IfNull(value float64) Expr {
	return field.ifNull(value)
}

func (field Float64) toSlice(values ...float64) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

// Float32 float32 type field
type Float32 Float64

// Eq equal to
func (field Float32) Eq(value float32) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq not equal to
func (field Float32) Neq(value float32) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// Gt greater than
func (field Float32) Gt(value float32) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte greater or equal to
func (field Float32) Gte(value float32) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt less than
func (field Float32) Lt(value float32) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte less or equal to
func (field Float32) Lte(value float32) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// In ...
func (field Float32) In(values ...float32) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// NotIn ...
func (field Float32) NotIn(values ...float32) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Between ...
func (field Float32) Between(left float32, right float32) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field Float32) NotBetween(left float32, right float32) Expr {
	return Not(field.Between(left, right))
}

// Like ...
func (field Float32) Like(value float32) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// NotLike ...
func (field Float32) NotLike(value float32) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

// Add ...
func (field Float32) Add(value float32) Float32 {
	return Float32{field.add(value)}
}

// Sub ...
func (field Float32) Sub(value float32) Float32 {
	return Float32{field.sub(value)}
}

// Mul ...
func (field Float32) Mul(value float32) Float32 {
	return Float32{field.mul(value)}
}

// Div ...
func (field Float32) Div(value float32) Float32 {
	return Float32{field.div(value)}
}

// FloorDiv ...
func (field Float32) FloorDiv(value float32) Int {
	return Int{field.floorDiv(value)}
}

// Floor ...
func (field Float32) Floor() Int {
	return Int{field.floor()}
}

// Value set value
func (field Float32) Value(value float32) AssignExpr {
	return field.value(value)
}

// Zero set zero value
func (field Float32) Zero() AssignExpr {
	return field.value(0)
}

// Sum ...
func (field Float32) Sum() Float32 {
	return Float32{field.sum()}
}

// IfNull ...
func (field Float32) IfNull(value float32) Expr {
	return field.ifNull(value)
}

func (field Float32) toSlice(values ...float32) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}
