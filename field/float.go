package field

import "gorm.io/gorm/clause"

type Float64 Field

func (field Float64) Eq(value float64) Expr {
	return expr{expression: clause.Eq{Column: field.Col, Value: value}}
}

func (field Float64) Neq(value float64) Expr {
	return expr{expression: clause.Neq{Column: field.Col, Value: value}}
}

func (field Float64) Gt(value float64) Expr {
	return expr{expression: clause.Gt{Column: field.Col, Value: value}}
}

func (field Float64) Gte(value float64) Expr {
	return expr{expression: clause.Gte{Column: field.Col, Value: value}}
}

func (field Float64) Lt(value float64) Expr {
	return expr{expression: clause.Lt{Column: field.Col, Value: value}}
}

func (field Float64) Lte(value float64) Expr {
	return expr{expression: clause.Lte{Column: field.Col, Value: value}}
}

func (field Float64) In(values ...float64) Expr {
	return expr{expression: clause.IN{Column: field.Col, Values: field.toSlice(values...)}}
}

func (field Float64) NotIn(values ...float64) Expr {
	return expr{expression: clause.Not(field.In(values...))}
}

func (field Float64) Between(left float64, right float64) Expr {
	return field.between([]interface{}{left, right})
}

func (field Float64) NotBetween(left float64, right float64) Expr {
	return Not(field.Between(left, right))
}

func (field Float64) Like(value float64) Expr {
	return expr{expression: clause.Like{Column: field.Col, Value: value}}
}

func (field Float64) NotLike(value float64) Expr {
	return expr{expression: clause.Not(field.Like(value))}
}

func (field Float64) Add(value float64) Expr {
	return field.add(value)
}

func (field Float64) Sub(value float64) Expr {
	return field.sub(value)
}

func (field Float64) Mul(value float64) Expr {
	return field.mul(value)
}

func (field Float64) Div(value float64) Expr {
	return field.div(value)
}

func (field Float64) FloorDiv(value float64) Expr {
	return field.floorDiv(value)
}

func (field Float64) toSlice(values ...float64) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

type Float32 Float64

func (field Float32) Eq(value float32) Expr {
	return expr{expression: clause.Eq{Column: field.Col, Value: value}}
}

func (field Float32) Neq(value float32) Expr {
	return expr{expression: clause.Neq{Column: field.Col, Value: value}}
}

func (field Float32) Gt(value float32) Expr {
	return expr{expression: clause.Gt{Column: field.Col, Value: value}}
}

func (field Float32) Gte(value float32) Expr {
	return expr{expression: clause.Gte{Column: field.Col, Value: value}}
}

func (field Float32) Lt(value float32) Expr {
	return expr{expression: clause.Lt{Column: field.Col, Value: value}}
}

func (field Float32) Lte(value float32) Expr {
	return expr{expression: clause.Lte{Column: field.Col, Value: value}}
}

func (field Float32) In(values ...float32) Expr {
	return expr{expression: clause.IN{Column: field.Col, Values: field.toSlice(values...)}}
}

func (field Float32) NotIn(values ...float32) Expr {
	return expr{expression: clause.Not(field.In(values...))}
}

func (field Float32) Between(left float32, right float32) Expr {
	return field.between([]interface{}{left, right})
}

func (field Float32) NotBetween(left float32, right float32) Expr {
	return Not(field.Between(left, right))
}

func (field Float32) Like(value float32) Expr {
	return expr{expression: clause.Like{Column: field.Col, Value: value}}
}

func (field Float32) NotLike(value float32) Expr {
	return expr{expression: clause.Not(field.Like(value))}
}

func (field Float32) Add(value float32) Expr {
	return field.add(value)
}

func (field Float32) Sub(value float32) Expr {
	return field.sub(value)
}

func (field Float32) Mul(value float32) Expr {
	return field.mul(value)
}

func (field Float32) Div(value float32) Expr {
	return field.div(value)
}

func (field Float32) FloorDiv(value float32) Expr {
	return field.floorDiv(value)
}

func (field Float32) toSlice(values ...float32) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}
