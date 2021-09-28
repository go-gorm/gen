package field

import (
	"gorm.io/gorm/clause"
)

type String Field

func (field String) Eq(value string) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

func (field String) Neq(value string) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

func (field String) Gt(value string) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

func (field String) Gte(value string) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

func (field String) Lt(value string) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

func (field String) Lte(value string) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

func (field String) Between(left string, right string) Expr {
	return field.between([]interface{}{left, right})
}

func (field String) NotBetween(left string, right string) Expr {
	return Not(field.Between(left, right))
}

func (field String) In(values ...string) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values)}}
}

func (field String) NotIn(values ...string) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

func (field String) Like(value string) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

func (field String) NotLike(value string) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

func (field String) Regexp(value string) Expr {
	return field.regexp(value)
}

func (field String) NotRegxp(value string) Expr {
	return expr{e: clause.Not(field.Regexp(value).expression())}
}

func (field String) Value(value string) AssignExpr {
	return field.value(value)
}

func (field String) Zero() AssignExpr {
	return field.value("")
}

func (field String) toSlice(values []string) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

type Bytes String

func (field Bytes) Eq(value []byte) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

func (field Bytes) Neq(value []byte) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

func (field Bytes) Gt(value []byte) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

func (field Bytes) Gte(value []byte) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

func (field Bytes) Lt(value []byte) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

func (field Bytes) Lte(value []byte) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

func (field Bytes) Between(left []byte, right []byte) Expr {
	return field.between([]interface{}{left, right})
}

func (field Bytes) NotBetween(left []byte, right []byte) Expr {
	return Not(field.Between(left, right))
}

func (field Bytes) In(values ...[]byte) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values)}}
}

func (field Bytes) NotIn(values ...[]byte) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

func (field Bytes) Like(value string) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

func (field Bytes) NotLike(value string) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

func (field Bytes) Regexp(value string) Expr {
	return field.regexp(value)
}

func (field Bytes) NotRegxp(value string) Expr {
	return Not(field.Regexp(value))
}

func (field Bytes) Value(value []byte) AssignExpr {
	return field.value(value)
}

func (field Bytes) Zero() AssignExpr {
	return field.value([]byte{})
}

func (field Bytes) toSlice(values [][]byte) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}
