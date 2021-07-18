package field

import (
	"gorm.io/gorm/clause"
)

type Int Field

func (field Int) Eq(value int) Expr {
	return expr{expression: clause.Eq{Column: field.Col, Value: value}}
}

func (field Int) Neq(value int) Expr {
	return expr{expression: clause.Neq{Column: field.Col, Value: value}}
}

func (field Int) Gt(value int) Expr {
	return expr{expression: clause.Gt{Column: field.Col, Value: value}}
}

func (field Int) Gte(value int) Expr {
	return expr{expression: clause.Gte{Column: field.Col, Value: value}}
}

func (field Int) Lt(value int) Expr {
	return expr{expression: clause.Lt{Column: field.Col, Value: value}}
}

func (field Int) Lte(value int) Expr {
	return expr{expression: clause.Lte{Column: field.Col, Value: value}}
}

func (field Int) In(values ...int) Expr {
	return expr{expression: clause.IN{Column: field.Col, Values: field.toSlice(values...)}}
}

func (field Int) NotIn(values ...int) Expr {
	return expr{expression: clause.Not(field.In(values...))}
}

func (field Int) Between(left int, right int) Expr {
	return field.between([]interface{}{left, right})
}

func (field Int) NotBetween(left int, right int) Expr {
	return Not(field.Between(left, right))
}

func (field Int) Like(value int) Expr {
	return expr{expression: clause.Like{Column: field.Col, Value: value}}
}

func (field Int) NotLike(value int) Expr {
	return expr{expression: clause.Not(field.Like(value))}
}

func (field Int) Add(value int) Expr {
	return field.add(value)
}

func (field Int) Sub(value int) Expr {
	return field.sub(value)
}

func (field Int) Mul(value int) Expr {
	return field.mul(value)
}

func (field Int) Div(value int) Expr {
	return field.div(value)
}

func (field Int) Mod(value int) Expr {
	return field.mod(value)
}

func (field Int) FloorDiv(value int) Expr {
	return field.floorDiv(value)
}

func (field Int) RightShift(value int) Expr {
	return field.rightShift(value)
}

func (field Int) LeftShift(value int) Expr {
	return field.leftShift(value)
}

func (field Int) BitXor(value int) Expr {
	return field.bitXor(value)
}

func (field Int) BitAnd(value int) Expr {
	return field.bitAnd(value)
}

func (field Int) BitOr(value int) Expr {
	return field.bitOr(value)
}

func (field Int) BitFlip() Expr {
	return field.bitFlip()
}

func (field Int) toSlice(values ...int) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

type Int8 Int

func (field Int8) Eq(value int8) Expr {
	return expr{expression: clause.Eq{Column: field.Col, Value: value}}
}

func (field Int8) Neq(value int8) Expr {
	return expr{expression: clause.Neq{Column: field.Col, Value: value}}
}

func (field Int8) Gt(value int8) Expr {
	return expr{expression: clause.Gt{Column: field.Col, Value: value}}
}

func (field Int8) Gte(value int8) Expr {
	return expr{expression: clause.Gte{Column: field.Col, Value: value}}
}

func (field Int8) Lt(value int8) Expr {
	return expr{expression: clause.Lt{Column: field.Col, Value: value}}
}

func (field Int8) Lte(value int8) Expr {
	return expr{expression: clause.Lte{Column: field.Col, Value: value}}
}

func (field Int8) In(values ...int8) Expr {
	return expr{expression: clause.IN{Column: field.Col, Values: field.toSlice(values...)}}
}

func (field Int8) NotIn(values ...int8) Expr {
	return expr{expression: clause.Not(field.In(values...))}
}

func (field Int8) Between(left int8, right int8) Expr {
	return field.between([]interface{}{left, right})
}

func (field Int8) NotBetween(left int8, right int8) Expr {
	return Not(field.Between(left, right))
}

func (field Int8) Like(value int8) Expr {
	return expr{expression: clause.Like{Column: field.Col, Value: value}}
}

func (field Int8) NotLike(value int8) Expr {
	return expr{expression: clause.Not(field.Like(value))}
}

func (field Int8) Add(value int8) Expr {
	return field.add(value)
}

func (field Int8) Sub(value int8) Expr {
	return field.sub(value)
}

func (field Int8) Mul(value int8) Expr {
	return field.mul(value)
}

func (field Int8) Div(value int8) Expr {
	return field.div(value)
}

func (field Int8) Mod(value int8) Expr {
	return field.mod(value)
}

func (field Int8) FloorDiv(value int8) Expr {
	return field.floorDiv(value)
}

func (field Int8) RightShift(value int8) Expr {
	return field.rightShift(value)
}

func (field Int8) LeftShift(value int8) Expr {
	return field.leftShift(value)
}

func (field Int8) BitXor(value int8) Expr {
	return field.bitXor(value)
}

func (field Int8) BitAnd(value int8) Expr {
	return field.bitAnd(value)
}

func (field Int8) BitOr(value int8) Expr {
	return field.bitOr(value)
}

func (field Int8) BitFlip() Expr {
	return field.bitFlip()
}

func (field Int8) toSlice(values ...int8) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

type Int16 Int

func (field Int16) Eq(value int16) Expr {
	return expr{expression: clause.Eq{Column: field.Col, Value: value}}
}

func (field Int16) Neq(value int16) Expr {
	return expr{expression: clause.Neq{Column: field.Col, Value: value}}
}

func (field Int16) Gt(value int16) Expr {
	return expr{expression: clause.Gt{Column: field.Col, Value: value}}
}

func (field Int16) Gte(value int16) Expr {
	return expr{expression: clause.Gte{Column: field.Col, Value: value}}
}

func (field Int16) Lt(value int16) Expr {
	return expr{expression: clause.Lt{Column: field.Col, Value: value}}
}

func (field Int16) Lte(value int16) Expr {
	return expr{expression: clause.Lte{Column: field.Col, Value: value}}
}

func (field Int16) In(values ...int16) Expr {
	return expr{expression: clause.IN{Column: field.Col, Values: field.toSlice(values...)}}
}

func (field Int16) NotIn(values ...int16) Expr {
	return expr{expression: clause.Not(field.In(values...))}
}

func (field Int16) Between(left int16, right int16) Expr {
	return field.between([]interface{}{left, right})
}

func (field Int16) NotBetween(left int16, right int16) Expr {
	return Not(field.Between(left, right))
}

func (field Int16) Like(value int16) Expr {
	return expr{expression: clause.Like{Column: field.Col, Value: value}}
}

func (field Int16) NotLike(value int16) Expr {
	return expr{expression: clause.Not(field.Like(value))}
}

func (field Int16) Add(value int16) Expr {
	return field.add(value)
}

func (field Int16) Sub(value int16) Expr {
	return field.sub(value)
}

func (field Int16) Mul(value int16) Expr {
	return field.mul(value)
}

func (field Int16) Div(value int16) Expr {
	return field.div(value)
}

func (field Int16) Mod(value int16) Expr {
	return field.mod(value)
}

func (field Int16) FloorDiv(value int16) Expr {
	return field.floorDiv(value)
}

func (field Int16) RightShift(value int16) Expr {
	return field.rightShift(value)
}

func (field Int16) LeftShift(value int16) Expr {
	return field.leftShift(value)
}

func (field Int16) BitXor(value int16) Expr {
	return field.bitXor(value)
}

func (field Int16) BitAnd(value int16) Expr {
	return field.bitAnd(value)
}

func (field Int16) BitOr(value int16) Expr {
	return field.bitOr(value)
}

func (field Int16) BitFlip() Expr {
	return field.bitFlip()
}

func (field Int16) toSlice(values ...int16) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

type Int32 Int

func (field Int32) Eq(value int32) Expr {
	return expr{expression: clause.Eq{Column: field.Col, Value: value}}
}

func (field Int32) Neq(value int32) Expr {
	return expr{expression: clause.Neq{Column: field.Col, Value: value}}
}

func (field Int32) Gt(value int32) Expr {
	return expr{expression: clause.Gt{Column: field.Col, Value: value}}
}

func (field Int32) Gte(value int32) Expr {
	return expr{expression: clause.Gte{Column: field.Col, Value: value}}
}

func (field Int32) Lt(value int32) Expr {
	return expr{expression: clause.Lt{Column: field.Col, Value: value}}
}

func (field Int32) Lte(value int32) Expr {
	return expr{expression: clause.Lte{Column: field.Col, Value: value}}
}

func (field Int32) In(values ...int32) Expr {
	return expr{expression: clause.IN{Column: field.Col, Values: field.toSlice(values...)}}
}

func (field Int32) NotIn(values ...int32) Expr {
	return expr{expression: clause.Not(field.In(values...))}
}

func (field Int32) Between(left int32, right int32) Expr {
	return field.between([]interface{}{left, right})
}

func (field Int32) NotBetween(left int32, right int32) Expr {
	return Not(field.Between(left, right))
}

func (field Int32) Like(value int32) Expr {
	return expr{expression: clause.Like{Column: field.Col, Value: value}}
}

func (field Int32) NotLike(value int32) Expr {
	return expr{expression: clause.Not(field.Like(value))}
}

func (field Int32) Add(value int32) Expr {
	return field.add(value)
}

func (field Int32) Sub(value int32) Expr {
	return field.sub(value)
}

func (field Int32) Mul(value int32) Expr {
	return field.mul(value)
}

func (field Int32) Div(value int32) Expr {
	return field.div(value)
}

func (field Int32) Mod(value int32) Expr {
	return field.mod(value)
}

func (field Int32) FloorDiv(value int32) Expr {
	return field.floorDiv(value)
}

func (field Int32) RightShift(value int32) Expr {
	return field.rightShift(value)
}

func (field Int32) LeftShift(value int32) Expr {
	return field.leftShift(value)
}

func (field Int32) BitXor(value int32) Expr {
	return field.bitXor(value)
}

func (field Int32) BitAnd(value int32) Expr {
	return field.bitAnd(value)
}

func (field Int32) BitOr(value int32) Expr {
	return field.bitOr(value)
}

func (field Int32) BitFlip() Expr {
	return field.bitFlip()
}

func (field Int32) toSlice(values ...int32) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

type Int64 Int

func (field Int64) Eq(value int64) Expr {
	return expr{expression: clause.Eq{Column: field.Col, Value: value}}
}

func (field Int64) Neq(value int64) Expr {
	return expr{expression: clause.Neq{Column: field.Col, Value: value}}
}

func (field Int64) Gt(value int64) Expr {
	return expr{expression: clause.Gt{Column: field.Col, Value: value}}
}

func (field Int64) Gte(value int64) Expr {
	return expr{expression: clause.Gte{Column: field.Col, Value: value}}
}

func (field Int64) Lt(value int64) Expr {
	return expr{expression: clause.Lt{Column: field.Col, Value: value}}
}

func (field Int64) Lte(value int64) Expr {
	return expr{expression: clause.Lte{Column: field.Col, Value: value}}
}

func (field Int64) In(values ...int64) Expr {
	return expr{expression: clause.IN{Column: field.Col, Values: field.toSlice(values...)}}
}

func (field Int64) NotIn(values ...int64) Expr {
	return expr{expression: clause.Not(field.In(values...))}
}

func (field Int64) Between(left int64, right int64) Expr {
	return field.between([]interface{}{left, right})
}

func (field Int64) NotBetween(left int64, right int64) Expr {
	return Not(field.Between(left, right))
}

func (field Int64) Like(value int64) Expr {
	return expr{expression: clause.Like{Column: field.Col, Value: value}}
}

func (field Int64) NotLike(value int64) Expr {
	return expr{expression: clause.Not(field.Like(value))}
}

func (field Int64) Add(value int64) Expr {
	return field.add(value)
}

func (field Int64) Sub(value int64) Expr {
	return field.sub(value)
}

func (field Int64) Mul(value int64) Expr {
	return field.mul(value)
}

func (field Int64) Div(value int64) Expr {
	return field.div(value)
}

func (field Int64) Mod(value int64) Expr {
	return field.mod(value)
}

func (field Int64) FloorDiv(value int64) Expr {
	return field.floorDiv(value)
}

func (field Int64) RightShift(value int64) Expr {
	return field.rightShift(value)
}

func (field Int64) LeftShift(value int64) Expr {
	return field.leftShift(value)
}

func (field Int64) BitXor(value int64) Expr {
	return field.bitXor(value)
}

func (field Int64) BitAnd(value int64) Expr {
	return field.bitAnd(value)
}

func (field Int64) BitOr(value int64) Expr {
	return field.bitOr(value)
}

func (field Int64) BitFlip() Expr {
	return field.bitFlip()
}

func (field Int64) toSlice(values ...int64) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

type Uint Int

func (field Uint) Eq(value uint) Expr {
	return expr{expression: clause.Eq{Column: field.Col, Value: value}}
}

func (field Uint) Neq(value uint) Expr {
	return expr{expression: clause.Neq{Column: field.Col, Value: value}}
}

func (field Uint) Gt(value uint) Expr {
	return expr{expression: clause.Gt{Column: field.Col, Value: value}}
}

func (field Uint) Gte(value uint) Expr {
	return expr{expression: clause.Gte{Column: field.Col, Value: value}}
}

func (field Uint) Lt(value uint) Expr {
	return expr{expression: clause.Lt{Column: field.Col, Value: value}}
}

func (field Uint) Lte(value uint) Expr {
	return expr{expression: clause.Lte{Column: field.Col, Value: value}}
}

func (field Uint) In(values ...uint) Expr {
	return expr{expression: clause.IN{Column: field.Col, Values: field.toSlice(values...)}}
}

func (field Uint) NotIn(values ...uint) Expr {
	return expr{expression: clause.Not(field.In(values...))}
}

func (field Uint) Between(left uint, right uint) Expr {
	return field.between([]interface{}{left, right})
}

func (field Uint) NotBetween(left uint, right uint) Expr {
	return Not(field.Between(left, right))
}

func (field Uint) Like(value uint) Expr {
	return expr{expression: clause.Like{Column: field.Col, Value: value}}
}

func (field Uint) NotLike(value uint) Expr {
	return expr{expression: clause.Not(field.Like(value))}
}

func (field Uint) Add(value uint) Expr {
	return field.add(value)
}

func (field Uint) Sub(value uint) Expr {
	return field.sub(value)
}

func (field Uint) Mul(value uint) Expr {
	return field.mul(value)
}

func (field Uint) Div(value uint) Expr {
	return field.div(value)
}

func (field Uint) Mod(value uint) Expr {
	return field.mod(value)
}

func (field Uint) FloorDiv(value uint) Expr {
	return field.floorDiv(value)
}

func (field Uint) RightShift(value uint) Expr {
	return field.rightShift(value)
}

func (field Uint) LeftShift(value uint) Expr {
	return field.leftShift(value)
}

func (field Uint) BitXor(value uint) Expr {
	return field.bitXor(value)
}

func (field Uint) BitAnd(value uint) Expr {
	return field.bitAnd(value)
}

func (field Uint) BitOr(value uint) Expr {
	return field.bitOr(value)
}

func (field Uint) BitFlip() Expr {
	return field.bitFlip()
}

func (field Uint) toSlice(values ...uint) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

type Uint8 Int

func (field Uint8) Eq(value uint8) Expr {
	return expr{expression: clause.Eq{Column: field.Col, Value: value}}
}

func (field Uint8) Neq(value uint8) Expr {
	return expr{expression: clause.Neq{Column: field.Col, Value: value}}
}

func (field Uint8) Gt(value uint8) Expr {
	return expr{expression: clause.Gt{Column: field.Col, Value: value}}
}

func (field Uint8) Gte(value uint8) Expr {
	return expr{expression: clause.Gte{Column: field.Col, Value: value}}
}

func (field Uint8) Lt(value uint8) Expr {
	return expr{expression: clause.Lt{Column: field.Col, Value: value}}
}

func (field Uint8) Lte(value uint8) Expr {
	return expr{expression: clause.Lte{Column: field.Col, Value: value}}
}

func (field Uint8) In(values ...uint8) Expr {
	return expr{expression: clause.IN{Column: field.Col, Values: field.toSlice(values...)}}
}

func (field Uint8) NotIn(values ...uint8) Expr {
	return expr{expression: clause.Not(field.In(values...))}
}

func (field Uint8) Between(left uint8, right uint8) Expr {
	return field.between([]interface{}{left, right})
}

func (field Uint8) NotBetween(left uint8, right uint8) Expr {
	return Not(field.Between(left, right))
}

func (field Uint8) Like(value uint8) Expr {
	return expr{expression: clause.Like{Column: field.Col, Value: value}}
}

func (field Uint8) NotLike(value uint8) Expr {
	return expr{expression: clause.Not(field.Like(value))}
}

func (field Uint8) Add(value uint8) Expr {
	return field.add(value)
}

func (field Uint8) Sub(value uint8) Expr {
	return field.sub(value)
}

func (field Uint8) Mul(value uint8) Expr {
	return field.mul(value)
}

func (field Uint8) Div(value uint8) Expr {
	return field.div(value)
}

func (field Uint8) Mod(value uint8) Expr {
	return field.mod(value)
}

func (field Uint8) FloorDiv(value uint8) Expr {
	return field.floorDiv(value)
}

func (field Uint8) RightShift(value uint8) Expr {
	return field.rightShift(value)
}

func (field Uint8) LeftShift(value uint8) Expr {
	return field.leftShift(value)
}

func (field Uint8) BitXor(value uint8) Expr {
	return field.bitXor(value)
}

func (field Uint8) BitAnd(value uint8) Expr {
	return field.bitAnd(value)
}

func (field Uint8) BitOr(value uint8) Expr {
	return field.bitOr(value)
}

func (field Uint8) BitFlip() Expr {
	return field.bitFlip()
}

func (field Uint8) toSlice(values ...uint8) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

type Uint16 Int

func (field Uint16) Eq(value uint16) Expr {
	return expr{expression: clause.Eq{Column: field.Col, Value: value}}
}

func (field Uint16) Neq(value uint16) Expr {
	return expr{expression: clause.Neq{Column: field.Col, Value: value}}
}

func (field Uint16) Gt(value uint16) Expr {
	return expr{expression: clause.Gt{Column: field.Col, Value: value}}
}

func (field Uint16) Gte(value uint16) Expr {
	return expr{expression: clause.Gte{Column: field.Col, Value: value}}
}

func (field Uint16) Lt(value uint16) Expr {
	return expr{expression: clause.Lt{Column: field.Col, Value: value}}
}

func (field Uint16) Lte(value uint16) Expr {
	return expr{expression: clause.Lte{Column: field.Col, Value: value}}
}

func (field Uint16) In(values ...uint16) Expr {
	return expr{expression: clause.IN{Column: field.Col, Values: field.toSlice(values...)}}
}

func (field Uint16) NotIn(values ...uint16) Expr {
	return expr{expression: clause.Not(field.In(values...))}
}

func (field Uint16) Between(left uint16, right uint16) Expr {
	return field.between([]interface{}{left, right})
}

func (field Uint16) NotBetween(left uint16, right uint16) Expr {
	return Not(field.Between(left, right))
}

func (field Uint16) Like(value uint16) Expr {
	return expr{expression: clause.Like{Column: field.Col, Value: value}}
}

func (field Uint16) NotLike(value uint16) Expr {
	return expr{expression: clause.Not(field.Like(value))}
}

func (field Uint16) Add(value uint16) Expr {
	return field.add(value)
}

func (field Uint16) Sub(value uint16) Expr {
	return field.sub(value)
}

func (field Uint16) Mul(value uint16) Expr {
	return field.mul(value)
}

func (field Uint16) Div(value uint16) Expr {
	return field.div(value)
}

func (field Uint16) Mod(value uint16) Expr {
	return field.mod(value)
}

func (field Uint16) FloorDiv(value uint16) Expr {
	return field.floorDiv(value)
}

func (field Uint16) RightShift(value uint16) Expr {
	return field.rightShift(value)
}

func (field Uint16) LeftShift(value uint16) Expr {
	return field.leftShift(value)
}

func (field Uint16) BitXor(value uint16) Expr {
	return field.bitXor(value)
}

func (field Uint16) BitAnd(value uint16) Expr {
	return field.bitAnd(value)
}

func (field Uint16) BitOr(value uint16) Expr {
	return field.bitOr(value)
}

func (field Uint16) BitFlip() Expr {
	return field.bitFlip()
}

func (field Uint16) toSlice(values ...uint16) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

type Uint32 Int

func (field Uint32) Eq(value uint32) Expr {
	return expr{expression: clause.Eq{Column: field.Col, Value: value}}
}

func (field Uint32) Neq(value uint32) Expr {
	return expr{expression: clause.Neq{Column: field.Col, Value: value}}
}

func (field Uint32) Gt(value uint32) Expr {
	return expr{expression: clause.Gt{Column: field.Col, Value: value}}
}

func (field Uint32) Gte(value uint32) Expr {
	return expr{expression: clause.Gte{Column: field.Col, Value: value}}
}

func (field Uint32) Lt(value uint32) Expr {
	return expr{expression: clause.Lt{Column: field.Col, Value: value}}
}

func (field Uint32) Lte(value uint32) Expr {
	return expr{expression: clause.Lte{Column: field.Col, Value: value}}
}

func (field Uint32) In(values ...uint32) Expr {
	return expr{expression: clause.IN{Column: field.Col, Values: field.toSlice(values...)}}
}

func (field Uint32) NotIn(values ...uint32) Expr {
	return expr{expression: clause.Not(field.In(values...))}
}

func (field Uint32) Between(left uint32, right uint32) Expr {
	return field.between([]interface{}{left, right})
}

func (field Uint32) NotBetween(left uint32, right uint32) Expr {
	return Not(field.Between(left, right))
}

func (field Uint32) Like(value uint32) Expr {
	return expr{expression: clause.Like{Column: field.Col, Value: value}}
}

func (field Uint32) NotLike(value uint32) Expr {
	return expr{expression: clause.Not(field.Like(value))}
}

func (field Uint32) Add(value uint32) Expr {
	return field.add(value)
}

func (field Uint32) Sub(value uint32) Expr {
	return field.sub(value)
}

func (field Uint32) Mul(value uint32) Expr {
	return field.mul(value)
}

func (field Uint32) Div(value uint32) Expr {
	return field.div(value)
}

func (field Uint32) Mod(value uint32) Expr {
	return field.mod(value)
}

func (field Uint32) FloorDiv(value uint32) Expr {
	return field.floorDiv(value)
}

func (field Uint32) RightShift(value uint32) Expr {
	return field.rightShift(value)
}

func (field Uint32) LeftShift(value uint32) Expr {
	return field.leftShift(value)
}

func (field Uint32) BitXor(value uint32) Expr {
	return field.bitXor(value)
}

func (field Uint32) BitAnd(value uint32) Expr {
	return field.bitAnd(value)
}

func (field Uint32) BitOr(value uint32) Expr {
	return field.bitOr(value)
}

func (field Uint32) BitFlip() Expr {
	return field.bitFlip()
}

func (field Uint32) toSlice(values ...uint32) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

type Uint64 Int

func (field Uint64) Eq(value uint64) Expr {
	return expr{expression: clause.Eq{Column: field.Col, Value: value}}
}

func (field Uint64) Neq(value uint64) Expr {
	return expr{expression: clause.Neq{Column: field.Col, Value: value}}
}

func (field Uint64) Gt(value uint64) Expr {
	return expr{expression: clause.Gt{Column: field.Col, Value: value}}
}

func (field Uint64) Gte(value uint64) Expr {
	return expr{expression: clause.Gte{Column: field.Col, Value: value}}
}

func (field Uint64) Lt(value uint64) Expr {
	return expr{expression: clause.Lt{Column: field.Col, Value: value}}
}

func (field Uint64) Lte(value uint64) Expr {
	return expr{expression: clause.Lte{Column: field.Col, Value: value}}
}

func (field Uint64) In(values ...uint64) Expr {
	return expr{expression: clause.IN{Column: field.Col, Values: field.toSlice(values...)}}
}

func (field Uint64) NotIn(values ...uint64) Expr {
	return expr{expression: clause.Not(field.In(values...))}
}

func (field Uint64) Between(left uint64, right uint64) Expr {
	return field.between([]interface{}{left, right})
}

func (field Uint64) NotBetween(left uint64, right uint64) Expr {
	return Not(field.Between(left, right))
}

func (field Uint64) Like(value uint64) Expr {
	return expr{expression: clause.Like{Column: field.Col, Value: value}}
}

func (field Uint64) NotLike(value uint64) Expr {
	return expr{expression: clause.Not(field.Like(value))}
}

func (field Uint64) Add(value uint64) Expr {
	return field.add(value)
}

func (field Uint64) Sub(value uint64) Expr {
	return field.sub(value)
}

func (field Uint64) Mul(value uint64) Expr {
	return field.mul(value)
}

func (field Uint64) Div(value uint64) Expr {
	return field.div(value)
}

func (field Uint64) Mod(value uint64) Expr {
	return field.mod(value)
}

func (field Uint64) FloorDiv(value uint64) Expr {
	return field.floorDiv(value)
}

func (field Uint64) RightShift(value uint64) Expr {
	return field.rightShift(value)
}

func (field Uint64) LeftShift(value uint64) Expr {
	return field.leftShift(value)
}

func (field Uint64) BitXor(value uint64) Expr {
	return field.bitXor(value)
}

func (field Uint64) BitAnd(value uint64) Expr {
	return field.bitAnd(value)
}

func (field Uint64) BitOr(value uint64) Expr {
	return field.bitOr(value)
}

func (field Uint64) BitFlip() Expr {
	return field.bitFlip()
}

func (field Uint64) toSlice(values ...uint64) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}
