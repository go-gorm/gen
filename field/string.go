package field

import (
	"gorm.io/gorm/clause"
)

// String string type field
type String Field

// Eq equal to
func (field String) Eq(value string) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq not equal to
func (field String) Neq(value string) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// Gt greater than
func (field String) Gt(value string) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte greater or equal to
func (field String) Gte(value string) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt less than
func (field String) Lt(value string) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte less or equal to
func (field String) Lte(value string) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// Between ...
func (field String) Between(left string, right string) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field String) NotBetween(left string, right string) Expr {
	return Not(field.Between(left, right))
}

// In ...
func (field String) In(values ...string) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values)}}
}

// NotIn ...
func (field String) NotIn(values ...string) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Like ...
func (field String) Like(value string) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// NotLike ...
func (field String) NotLike(value string) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

// Regexp ...
func (field String) Regexp(value string) Expr {
	return field.regexp(value)
}

// NotRegxp ...
func (field String) NotRegxp(value string) Expr {
	return expr{e: clause.Not(field.Regexp(value).expression())}
}

// Value ...
func (field String) Value(value string) AssignExpr {
	return field.value(value)
}

// Zero ...
func (field String) Zero() AssignExpr {
	return field.value("")
}

// IfNull ...
func (field String) IfNull(value string) Expr {
	return field.ifNull(value)
}

// FindInSet equal to FIND_IN_SET(field_name, input_string_list)
func (field String) FindInSet(targetList string) Expr {
	return expr{e: clause.Expr{SQL: "FIND_IN_SET(?,?)", Vars: []interface{}{field.RawExpr(), targetList}}}
}

// FindInSetWith equal to FIND_IN_SET(input_string, field_name)
func (field String) FindInSetWith(target string) Expr {
	return expr{e: clause.Expr{SQL: "FIND_IN_SET(?,?)", Vars: []interface{}{target, field.RawExpr()}}}
}

// Replace ...
func (field String) Replace(from, to string) String {
	return String{expr{e: clause.Expr{SQL: "REPLACE(?,?,?)", Vars: []interface{}{field.RawExpr(), from, to}}}}
}

// Concat ...
func (field String) Concat(before, after string) String {
	switch {
	case before != "" && after != "":
		return String{expr{e: clause.Expr{SQL: "CONCAT(?,?,?)", Vars: []interface{}{before, field.RawExpr(), after}}}}
	case before != "":
		return String{expr{e: clause.Expr{SQL: "CONCAT(?,?)", Vars: []interface{}{before, field.RawExpr()}}}}
	case after != "":
		return String{expr{e: clause.Expr{SQL: "CONCAT(?,?)", Vars: []interface{}{field.RawExpr(), after}}}}
	default:
		return field
	}
}

func (field String) toSlice(values []string) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

// Bytes []byte type field
type Bytes String

// Eq equal to
func (field Bytes) Eq(value []byte) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq not equal to
func (field Bytes) Neq(value []byte) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// Gt greater than
func (field Bytes) Gt(value []byte) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte greater or equal to
func (field Bytes) Gte(value []byte) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt less than
func (field Bytes) Lt(value []byte) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte less or equal to
func (field Bytes) Lte(value []byte) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// Between ...
func (field Bytes) Between(left []byte, right []byte) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field Bytes) NotBetween(left []byte, right []byte) Expr {
	return Not(field.Between(left, right))
}

// In ...
func (field Bytes) In(values ...[]byte) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values)}}
}

// NotIn ...
func (field Bytes) NotIn(values ...[]byte) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Like ...
func (field Bytes) Like(value string) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// NotLike ...
func (field Bytes) NotLike(value string) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

// Regexp ...
func (field Bytes) Regexp(value string) Expr {
	return field.regexp(value)
}

// NotRegxp ...
func (field Bytes) NotRegxp(value string) Expr {
	return Not(field.Regexp(value))
}

// Value ...
func (field Bytes) Value(value []byte) AssignExpr {
	return field.value(value)
}

// Zero ...
func (field Bytes) Zero() AssignExpr {
	return field.value([]byte{})
}

// IfNull ...
func (field Bytes) IfNull(value []byte) Expr {
	return field.ifNull(value)
}

// FindInSet FIND_IN_SET(field_name, input_string_list)
func (field Bytes) FindInSet(targetList string) Expr {
	return expr{e: clause.Expr{SQL: "FIND_IN_SET(?,?)", Vars: []interface{}{field.RawExpr(), targetList}}}
}

// FindInSetWith FIND_IN_SET(input_string, field_name)
func (field Bytes) FindInSetWith(target string) Expr {
	return expr{e: clause.Expr{SQL: "FIND_IN_SET(?,?)", Vars: []interface{}{target, field.RawExpr()}}}
}

func (field Bytes) toSlice(values [][]byte) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}
