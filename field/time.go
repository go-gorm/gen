package field

import (
	"time"

	"gorm.io/gorm/clause"
)

type Time Field

func (field Time) Eq(value time.Time) Expr {
	return expr{expression: clause.Eq{Column: field.Col, Value: value}}
}

func (field Time) Neq(value time.Time) Expr {
	return expr{expression: clause.Neq{Column: field.Col, Value: value}}
}

func (field Time) Gt(value time.Time) Expr {
	return expr{expression: clause.Gt{Column: field.Col, Value: value}}
}

func (field Time) Gte(value time.Time) Expr {
	return expr{expression: clause.Gte{Column: field.Col, Value: value}}
}

func (field Time) Lt(value time.Time) Expr {
	return expr{expression: clause.Lt{Column: field.Col, Value: value}}
}

func (field Time) Lte(value time.Time) Expr {
	return expr{expression: clause.Lte{Column: field.Col, Value: value}}
}

func (field Time) Between(left time.Time, right time.Time) Expr {
	return field.between([]interface{}{left, right})
}

func (field Time) NotBetween(left time.Time, right time.Time) Expr {
	return Not(field.Between(left, right))
}

func (field Time) In(values ...time.Time) Expr {
	return expr{expression: clause.IN{Column: field.Col, Values: field.toSlice(values...)}}
}

func (field Time) NotIn(values ...time.Time) Expr {
	return expr{expression: clause.Not(field.In(values...))}
}

func (field Time) Add(value time.Duration) Expr {
	return field.add(value)
}

func (field Time) Sub(value time.Duration) Expr {
	return field.sub(value)
}

func (field Time) toSlice(values ...time.Time) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}
