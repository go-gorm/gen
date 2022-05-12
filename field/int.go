package field

import (
	"gorm.io/gorm/clause"
)

// Int int type field
type Int Field

// Eq equal to
func (field Int) Eq(value int) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq not equal to
func (field Int) Neq(value int) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// Gt greater than
func (field Int) Gt(value int) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte greater or equal to
func (field Int) Gte(value int) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt less than
func (field Int) Lt(value int) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte less or equal to
func (field Int) Lte(value int) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// In ...
func (field Int) In(values ...int) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// NotIn ...
func (field Int) NotIn(values ...int) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Between ...
func (field Int) Between(left int, right int) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field Int) NotBetween(left int, right int) Expr {
	return Not(field.Between(left, right))
}

// Like ...
func (field Int) Like(value int) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// NotLike ...
func (field Int) NotLike(value int) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

// Add ...
func (field Int) Add(value int) Int {
	return Int{field.add(value)}
}

// Sub ...
func (field Int) Sub(value int) Int {
	return Int{field.sub(value)}
}

// Mul ...
func (field Int) Mul(value int) Int {
	return Int{field.mul(value)}
}

// Div ...
func (field Int) Div(value int) Int {
	return Int{field.div(value)}
}

// Mod ...
func (field Int) Mod(value int) Int {
	return Int{field.mod(value)}
}

// FloorDiv ...
func (field Int) FloorDiv(value int) Int {
	return Int{field.floorDiv(value)}
}

// RightShift ...
func (field Int) RightShift(value int) Int {
	return Int{field.rightShift(value)}
}

// LeftShift ...
func (field Int) LeftShift(value int) Int {
	return Int{field.leftShift(value)}
}

// BitXor ...
func (field Int) BitXor(value int) Int {
	return Int{field.bitXor(value)}
}

// BitAnd ...
func (field Int) BitAnd(value int) Int {
	return Int{field.bitAnd(value)}
}

// BitOr ...
func (field Int) BitOr(value int) Int {
	return Int{field.bitOr(value)}
}

// BitFlip ...
func (field Int) BitFlip() Int {
	return Int{field.bitFlip()}
}

// Value set value
func (field Int) Value(value int) AssignExpr {
	return field.value(value)
}

// Zero set zero value
func (field Int) Zero() AssignExpr {
	return field.value(0)
}

// Sum ...
func (field Int) Sum() Int {
	return Int{field.sum()}
}

// IfNull ...
func (field Int) IfNull(value int) Expr {
	return field.ifNull(value)
}

func (field Int) toSlice(values ...int) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

// Int8 int8 type field
type Int8 Int

// Eq equal to
func (field Int8) Eq(value int8) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq not equal to
func (field Int8) Neq(value int8) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// Gt greater than
func (field Int8) Gt(value int8) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte greater or equal to
func (field Int8) Gte(value int8) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt less than
func (field Int8) Lt(value int8) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte less or equal to
func (field Int8) Lte(value int8) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// In ...
func (field Int8) In(values ...int8) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// NotIn ...
func (field Int8) NotIn(values ...int8) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Between ...
func (field Int8) Between(left int8, right int8) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field Int8) NotBetween(left int8, right int8) Expr {
	return Not(field.Between(left, right))
}

// Like ...
func (field Int8) Like(value int8) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// NotLike ...
func (field Int8) NotLike(value int8) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

// Add ...
func (field Int8) Add(value int8) Int8 {
	return Int8{field.add(value)}
}

// Sub ...
func (field Int8) Sub(value int8) Int8 {
	return Int8{field.sub(value)}
}

// Mul ...
func (field Int8) Mul(value int8) Int8 {
	return Int8{field.mul(value)}
}

// Div ...
func (field Int8) Div(value int8) Int8 {
	return Int8{field.div(value)}
}

// Mod ...
func (field Int8) Mod(value int8) Int8 {
	return Int8{field.mod(value)}
}

// FloorDiv ...
func (field Int8) FloorDiv(value int8) Int8 {
	return Int8{field.floorDiv(value)}
}

// RightShift ...
func (field Int8) RightShift(value int8) Int8 {
	return Int8{field.rightShift(value)}
}

// LeftShift ...
func (field Int8) LeftShift(value int8) Int8 {
	return Int8{field.leftShift(value)}
}

// BitXor ...
func (field Int8) BitXor(value int8) Int8 {
	return Int8{field.bitXor(value)}
}

// BitAnd ...
func (field Int8) BitAnd(value int8) Int8 {
	return Int8{field.bitAnd(value)}
}

// BitOr ...
func (field Int8) BitOr(value int8) Int8 {
	return Int8{field.bitOr(value)}
}

// BitFlip ...
func (field Int8) BitFlip() Int8 {
	return Int8{field.bitFlip()}
}

// Value set value
func (field Int8) Value(value int8) AssignExpr {
	return field.value(value)
}

// Zero set zero value
func (field Int8) Zero() AssignExpr {
	return field.value(0)
}

// Sum ...
func (field Int8) Sum() Int8 {
	return Int8{field.sum()}
}

// IfNull ...
func (field Int8) IfNull(value int8) Expr {
	return field.ifNull(value)
}

func (field Int8) toSlice(values ...int8) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

// Int16 int16 type field
type Int16 Int

// Eq equal to
func (field Int16) Eq(value int16) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq not equal to
func (field Int16) Neq(value int16) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// Gt greater than
func (field Int16) Gt(value int16) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte greater or equal to
func (field Int16) Gte(value int16) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt less than
func (field Int16) Lt(value int16) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte less or equal to
func (field Int16) Lte(value int16) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// In ...
func (field Int16) In(values ...int16) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// NotIn ...
func (field Int16) NotIn(values ...int16) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Between ...
func (field Int16) Between(left int16, right int16) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field Int16) NotBetween(left int16, right int16) Expr {
	return Not(field.Between(left, right))
}

// Like ...
func (field Int16) Like(value int16) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// NotLike ...
func (field Int16) NotLike(value int16) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

// Add ...
func (field Int16) Add(value int16) Int16 {
	return Int16{field.add(value)}
}

// Sub ...
func (field Int16) Sub(value int16) Int16 {
	return Int16{field.sub(value)}
}

// Mul ...
func (field Int16) Mul(value int16) Int16 {
	return Int16{field.mul(value)}
}

// Div ...
func (field Int16) Div(value int16) Int16 {
	return Int16{field.div(value)}
}

// Mod ...
func (field Int16) Mod(value int16) Int16 {
	return Int16{field.mod(value)}
}

// FloorDiv ...
func (field Int16) FloorDiv(value int16) Int16 {
	return Int16{field.floorDiv(value)}
}

// RightShift ...
func (field Int16) RightShift(value int16) Int16 {
	return Int16{field.rightShift(value)}
}

// LeftShift ...
func (field Int16) LeftShift(value int16) Int16 {
	return Int16{field.leftShift(value)}
}

// BitXor ...
func (field Int16) BitXor(value int16) Int16 {
	return Int16{field.bitXor(value)}
}

// BitAnd ...
func (field Int16) BitAnd(value int16) Int16 {
	return Int16{field.bitAnd(value)}
}

// BitOr ...
func (field Int16) BitOr(value int16) Int16 {
	return Int16{field.bitOr(value)}
}

// BitFlip ...
func (field Int16) BitFlip() Int16 {
	return Int16{field.bitFlip()}
}

// Value set value
func (field Int16) Value(value int16) AssignExpr {
	return field.value(value)
}

// Zero set zero value
func (field Int16) Zero() AssignExpr {
	return field.value(0)
}

// Sum ...
func (field Int16) Sum() Int16 {
	return Int16{field.sum()}
}

// IfNull ...
func (field Int16) IfNull(value int16) Expr {
	return field.ifNull(value)
}

func (field Int16) toSlice(values ...int16) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

// Int32 int32 type field
type Int32 Int

// Eq equal to
func (field Int32) Eq(value int32) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq not equal to
func (field Int32) Neq(value int32) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// Gt greater than
func (field Int32) Gt(value int32) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte greater or equal to
func (field Int32) Gte(value int32) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt less than
func (field Int32) Lt(value int32) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte less or equal to
func (field Int32) Lte(value int32) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// In ...
func (field Int32) In(values ...int32) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// NotIn ...
func (field Int32) NotIn(values ...int32) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Between ...
func (field Int32) Between(left int32, right int32) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field Int32) NotBetween(left int32, right int32) Expr {
	return Not(field.Between(left, right))
}

// Like ...
func (field Int32) Like(value int32) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// NotLike ...
func (field Int32) NotLike(value int32) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

// Add ...
func (field Int32) Add(value int32) Int32 {
	return Int32{field.add(value)}
}

// Sub ...
func (field Int32) Sub(value int32) Int32 {
	return Int32{field.sub(value)}
}

// Mul ...
func (field Int32) Mul(value int32) Int32 {
	return Int32{field.mul(value)}
}

// Div ...
func (field Int32) Div(value int32) Int32 {
	return Int32{field.div(value)}
}

// Mod ...
func (field Int32) Mod(value int32) Int32 {
	return Int32{field.mod(value)}
}

// FloorDiv ...
func (field Int32) FloorDiv(value int32) Int32 {
	return Int32{field.floorDiv(value)}
}

// RightShift ...
func (field Int32) RightShift(value int32) Int32 {
	return Int32{field.rightShift(value)}
}

// LeftShift ...
func (field Int32) LeftShift(value int32) Int32 {
	return Int32{field.leftShift(value)}
}

// BitXor ...
func (field Int32) BitXor(value int32) Int32 {
	return Int32{field.bitXor(value)}
}

// BitAnd ...
func (field Int32) BitAnd(value int32) Int32 {
	return Int32{field.bitAnd(value)}
}

// BitOr ...
func (field Int32) BitOr(value int32) Int32 {
	return Int32{field.bitOr(value)}
}

// BitFlip ...
func (field Int32) BitFlip() Int32 {
	return Int32{field.bitFlip()}
}

// Value set value
func (field Int32) Value(value int32) AssignExpr {
	return field.value(value)
}

// Zero set zero value
func (field Int32) Zero() AssignExpr {
	return field.value(0)
}

// Sum ...
func (field Int32) Sum() Int32 {
	return Int32{field.sum()}
}

// IfNull ...
func (field Int32) IfNull(value int32) Expr {
	return field.ifNull(value)
}

func (field Int32) toSlice(values ...int32) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

// Int64 int64 type field
type Int64 Int

// Eq equal to
func (field Int64) Eq(value int64) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq not equal to
func (field Int64) Neq(value int64) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// Gt greater than
func (field Int64) Gt(value int64) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte greater or equal to
func (field Int64) Gte(value int64) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt less than
func (field Int64) Lt(value int64) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte less or equal to
func (field Int64) Lte(value int64) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// In ...
func (field Int64) In(values ...int64) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// NotIn ...
func (field Int64) NotIn(values ...int64) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Between ...
func (field Int64) Between(left int64, right int64) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field Int64) NotBetween(left int64, right int64) Expr {
	return Not(field.Between(left, right))
}

// Like ...
func (field Int64) Like(value int64) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// NotLike ...
func (field Int64) NotLike(value int64) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

// Add ...
func (field Int64) Add(value int64) Int64 {
	return Int64{field.add(value)}
}

// Sub ...
func (field Int64) Sub(value int64) Int64 {
	return Int64{field.sub(value)}
}

// Mul ...
func (field Int64) Mul(value int64) Int64 {
	return Int64{field.mul(value)}
}

// Div ...
func (field Int64) Div(value int64) Int64 {
	return Int64{field.div(value)}
}

// Mod ...
func (field Int64) Mod(value int64) Int64 {
	return Int64{field.mod(value)}
}

// FloorDiv ...
func (field Int64) FloorDiv(value int64) Int64 {
	return Int64{field.floorDiv(value)}
}

// RightShift ...
func (field Int64) RightShift(value int64) Int64 {
	return Int64{field.rightShift(value)}
}

// LeftShift ...
func (field Int64) LeftShift(value int64) Int64 {
	return Int64{field.leftShift(value)}
}

// BitXor ...
func (field Int64) BitXor(value int64) Int64 {
	return Int64{field.bitXor(value)}
}

// BitAnd ...
func (field Int64) BitAnd(value int64) Int64 {
	return Int64{field.bitAnd(value)}
}

// BitOr ...
func (field Int64) BitOr(value int64) Int64 {
	return Int64{field.bitOr(value)}
}

// BitFlip ...
func (field Int64) BitFlip() Int64 {
	return Int64{field.bitFlip()}
}

// Value set value
func (field Int64) Value(value int64) AssignExpr {
	return field.value(value)
}

// Zero set zero value
func (field Int64) Zero() AssignExpr {
	return field.value(0)
}

// Sum ...
func (field Int64) Sum() Int64 {
	return Int64{field.sum()}
}

// IfNull ...
func (field Int64) IfNull(value int64) Expr {
	return field.ifNull(value)
}

func (field Int64) toSlice(values ...int64) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

// Uint uint type field
type Uint Int

// Eq equal to
func (field Uint) Eq(value uint) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq not equal to
func (field Uint) Neq(value uint) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// Gt greater than
func (field Uint) Gt(value uint) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte greater or equal to
func (field Uint) Gte(value uint) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt less than
func (field Uint) Lt(value uint) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte less or equal to
func (field Uint) Lte(value uint) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// In ...
func (field Uint) In(values ...uint) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// NotIn ...
func (field Uint) NotIn(values ...uint) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Between ...
func (field Uint) Between(left uint, right uint) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field Uint) NotBetween(left uint, right uint) Expr {
	return Not(field.Between(left, right))
}

// Like ...
func (field Uint) Like(value uint) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// NotLike ...
func (field Uint) NotLike(value uint) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

// Add ...
func (field Uint) Add(value uint) Uint {
	return Uint{field.add(value)}
}

// Sub ...
func (field Uint) Sub(value uint) Uint {
	return Uint{field.sub(value)}
}

// Mul ...
func (field Uint) Mul(value uint) Uint {
	return Uint{field.mul(value)}
}

// Div ...
func (field Uint) Div(value uint) Uint {
	return Uint{field.mul(value)}
}

// Mod ...
func (field Uint) Mod(value uint) Uint {
	return Uint{field.mod(value)}
}

// FloorDiv ...
func (field Uint) FloorDiv(value uint) Uint {
	return Uint{field.floorDiv(value)}
}

// RightShift ...
func (field Uint) RightShift(value uint) Uint {
	return Uint{field.rightShift(value)}
}

// LeftShift ...
func (field Uint) LeftShift(value uint) Uint {
	return Uint{field.leftShift(value)}
}

// BitXor ...
func (field Uint) BitXor(value uint) Uint {
	return Uint{field.bitXor(value)}
}

// BitAnd ...
func (field Uint) BitAnd(value uint) Uint {
	return Uint{field.bitAnd(value)}
}

// BitOr ...
func (field Uint) BitOr(value uint) Uint {
	return Uint{field.bitOr(value)}
}

// BitFlip ...
func (field Uint) BitFlip() Uint {
	return Uint{field.bitFlip()}
}

// Value set value
func (field Uint) Value(value uint) AssignExpr {
	return field.value(value)
}

// Zero set zero value
func (field Uint) Zero() AssignExpr {
	return field.value(0)
}

// Sum ...
func (field Uint) Sum() Uint {
	return Uint{field.sum()}
}

// IfNull ...
func (field Uint) IfNull(value uint) Expr {
	return field.ifNull(value)
}

func (field Uint) toSlice(values ...uint) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

// Uint8 uint8 type field
type Uint8 Int

// Eq equal to
func (field Uint8) Eq(value uint8) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq not equal to
func (field Uint8) Neq(value uint8) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// Gt greater than
func (field Uint8) Gt(value uint8) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte greater or equal to
func (field Uint8) Gte(value uint8) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt less than
func (field Uint8) Lt(value uint8) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte less or equal to
func (field Uint8) Lte(value uint8) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// In ...
func (field Uint8) In(values ...uint8) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// NotIn ...
func (field Uint8) NotIn(values ...uint8) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Between ...
func (field Uint8) Between(left uint8, right uint8) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field Uint8) NotBetween(left uint8, right uint8) Expr {
	return Not(field.Between(left, right))
}

// Like ...
func (field Uint8) Like(value uint8) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// NotLike ...
func (field Uint8) NotLike(value uint8) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

// Add ...
func (field Uint8) Add(value uint8) Uint8 {
	return Uint8{field.add(value)}
}

// Sub ...
func (field Uint8) Sub(value uint8) Uint8 {
	return Uint8{field.sub(value)}
}

// Mul ...
func (field Uint8) Mul(value uint8) Uint8 {
	return Uint8{field.mul(value)}
}

// Div ...
func (field Uint8) Div(value uint8) Uint8 {
	return Uint8{field.mul(value)}
}

// Mod ...
func (field Uint8) Mod(value uint8) Uint8 {
	return Uint8{field.mod(value)}
}

// FloorDiv ...
func (field Uint8) FloorDiv(value uint8) Uint8 {
	return Uint8{field.floorDiv(value)}
}

// RightShift ...
func (field Uint8) RightShift(value uint8) Uint8 {
	return Uint8{field.rightShift(value)}
}

// LeftShift ...
func (field Uint8) LeftShift(value uint8) Uint8 {
	return Uint8{field.leftShift(value)}
}

// BitXor ...
func (field Uint8) BitXor(value uint8) Uint8 {
	return Uint8{field.bitXor(value)}
}

// BitAnd ...
func (field Uint8) BitAnd(value uint8) Uint8 {
	return Uint8{field.bitAnd(value)}
}

// BitOr ...
func (field Uint8) BitOr(value uint8) Uint8 {
	return Uint8{field.bitOr(value)}
}

// BitFlip ...
func (field Uint8) BitFlip() Uint8 {
	return Uint8{field.bitFlip()}
}

// Value set value
func (field Uint8) Value(value uint8) AssignExpr {
	return field.value(value)
}

// Zero set zero value
func (field Uint8) Zero() AssignExpr {
	return field.value(0)
}

// Sum ...
func (field Uint8) Sum() Uint8 {
	return Uint8{field.sum()}
}

// IfNull ...
func (field Uint8) IfNull(value uint8) Expr {
	return field.ifNull(value)
}

func (field Uint8) toSlice(values ...uint8) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

// Uint16 uint16 type field
type Uint16 Int

// Eq equal to
func (field Uint16) Eq(value uint16) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq not equal to
func (field Uint16) Neq(value uint16) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// Gt greater than
func (field Uint16) Gt(value uint16) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte greater or equal to
func (field Uint16) Gte(value uint16) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt less than
func (field Uint16) Lt(value uint16) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte less or equal to
func (field Uint16) Lte(value uint16) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// In ...
func (field Uint16) In(values ...uint16) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// NotIn ...
func (field Uint16) NotIn(values ...uint16) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Between ...
func (field Uint16) Between(left uint16, right uint16) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field Uint16) NotBetween(left uint16, right uint16) Expr {
	return Not(field.Between(left, right))
}

// Like ...
func (field Uint16) Like(value uint16) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// NotLike ...
func (field Uint16) NotLike(value uint16) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

// Add ...
func (field Uint16) Add(value uint16) Uint16 {
	return Uint16{field.add(value)}
}

// Sub ...
func (field Uint16) Sub(value uint16) Uint16 {
	return Uint16{field.sub(value)}
}

// Mul ...
func (field Uint16) Mul(value uint16) Uint16 {
	return Uint16{field.mul(value)}
}

// Div ...
func (field Uint16) Div(value uint16) Uint16 {
	return Uint16{field.mul(value)}
}

// Mod ...
func (field Uint16) Mod(value uint16) Uint16 {
	return Uint16{field.mod(value)}
}

// FloorDiv ...
func (field Uint16) FloorDiv(value uint16) Uint16 {
	return Uint16{field.floorDiv(value)}
}

// RightShift ...
func (field Uint16) RightShift(value uint16) Uint16 {
	return Uint16{field.rightShift(value)}
}

// LeftShift ...
func (field Uint16) LeftShift(value uint16) Uint16 {
	return Uint16{field.leftShift(value)}
}

// BitXor ...
func (field Uint16) BitXor(value uint16) Uint16 {
	return Uint16{field.bitXor(value)}
}

// BitAnd ...
func (field Uint16) BitAnd(value uint16) Uint16 {
	return Uint16{field.bitAnd(value)}
}

// BitOr ...
func (field Uint16) BitOr(value uint16) Uint16 {
	return Uint16{field.bitOr(value)}
}

// BitFlip ...
func (field Uint16) BitFlip() Uint16 {
	return Uint16{field.bitFlip()}
}

// Value set value
func (field Uint16) Value(value uint16) AssignExpr {
	return field.value(value)
}

// Zero set zero value
func (field Uint16) Zero() AssignExpr {
	return field.value(0)
}

// Sum ...
func (field Uint16) Sum() Uint16 {
	return Uint16{field.sum()}
}

// IfNull ...
func (field Uint16) IfNull(value uint16) Expr {
	return field.ifNull(value)
}

func (field Uint16) toSlice(values ...uint16) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

// Uint32 uint32 type field
type Uint32 Int

// Eq equal to
func (field Uint32) Eq(value uint32) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq not equal to
func (field Uint32) Neq(value uint32) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// Gt greater than
func (field Uint32) Gt(value uint32) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte greater or equal to
func (field Uint32) Gte(value uint32) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt less than
func (field Uint32) Lt(value uint32) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte less or equal to
func (field Uint32) Lte(value uint32) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// In ...
func (field Uint32) In(values ...uint32) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// NotIn ...
func (field Uint32) NotIn(values ...uint32) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Between ...
func (field Uint32) Between(left uint32, right uint32) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field Uint32) NotBetween(left uint32, right uint32) Expr {
	return Not(field.Between(left, right))
}

// Like ...
func (field Uint32) Like(value uint32) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// NotLike ...
func (field Uint32) NotLike(value uint32) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

// Add ...
func (field Uint32) Add(value uint32) Uint32 {
	return Uint32{field.add(value)}
}

// Sub ...
func (field Uint32) Sub(value uint32) Uint32 {
	return Uint32{field.sub(value)}
}

// Mul ...
func (field Uint32) Mul(value uint32) Uint32 {
	return Uint32{field.mul(value)}
}

// Div ...
func (field Uint32) Div(value uint32) Uint32 {
	return Uint32{field.mul(value)}
}

// Mod ...
func (field Uint32) Mod(value uint32) Uint32 {
	return Uint32{field.mod(value)}
}

// FloorDiv ...
func (field Uint32) FloorDiv(value uint32) Uint32 {
	return Uint32{field.floorDiv(value)}
}

// RightShift ...
func (field Uint32) RightShift(value uint32) Uint32 {
	return Uint32{field.rightShift(value)}
}

// LeftShift ...
func (field Uint32) LeftShift(value uint32) Uint32 {
	return Uint32{field.leftShift(value)}
}

// BitXor ...
func (field Uint32) BitXor(value uint32) Uint32 {
	return Uint32{field.bitXor(value)}
}

// BitAnd ...
func (field Uint32) BitAnd(value uint32) Uint32 {
	return Uint32{field.bitAnd(value)}
}

// BitOr ...
func (field Uint32) BitOr(value uint32) Uint32 {
	return Uint32{field.bitOr(value)}
}

// BitFlip ...
func (field Uint32) BitFlip() Uint32 {
	return Uint32{field.bitFlip()}
}

// Value set value
func (field Uint32) Value(value uint32) AssignExpr {
	return field.value(value)
}

// Zero set zero value
func (field Uint32) Zero() AssignExpr {
	return field.value(0)
}

// Sum ...
func (field Uint32) Sum() Uint32 {
	return Uint32{field.sum()}
}

// IfNull ...
func (field Uint32) IfNull(value uint32) Expr {
	return field.ifNull(value)
}

func (field Uint32) toSlice(values ...uint32) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

// Uint64 uint64 type field
type Uint64 Int

// Eq equal to
func (field Uint64) Eq(value uint64) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq not equal to
func (field Uint64) Neq(value uint64) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// Gt greater than
func (field Uint64) Gt(value uint64) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte greater or equal to
func (field Uint64) Gte(value uint64) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt less than
func (field Uint64) Lt(value uint64) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte less or equal to
func (field Uint64) Lte(value uint64) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// In ...
func (field Uint64) In(values ...uint64) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// NotIn ...
func (field Uint64) NotIn(values ...uint64) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Between ...
func (field Uint64) Between(left uint64, right uint64) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field Uint64) NotBetween(left uint64, right uint64) Expr {
	return Not(field.Between(left, right))
}

// Like ...
func (field Uint64) Like(value uint64) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// NotLike ...
func (field Uint64) NotLike(value uint64) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

// Add ...
func (field Uint64) Add(value uint64) Uint64 {
	return Uint64{field.add(value)}
}

// Sub ...
func (field Uint64) Sub(value uint64) Uint64 {
	return Uint64{field.sub(value)}
}

// Mul ...
func (field Uint64) Mul(value uint64) Uint64 {
	return Uint64{field.mul(value)}
}

// Div ...
func (field Uint64) Div(value uint64) Uint64 {
	return Uint64{field.mul(value)}
}

// Mod ...
func (field Uint64) Mod(value uint64) Uint64 {
	return Uint64{field.mod(value)}
}

// FloorDiv ...
func (field Uint64) FloorDiv(value uint64) Uint64 {
	return Uint64{field.floorDiv(value)}
}

// RightShift ...
func (field Uint64) RightShift(value uint64) Uint64 {
	return Uint64{field.rightShift(value)}
}

// LeftShift ...
func (field Uint64) LeftShift(value uint64) Uint64 {
	return Uint64{field.leftShift(value)}
}

// BitXor ...
func (field Uint64) BitXor(value uint64) Uint64 {
	return Uint64{field.bitXor(value)}
}

// BitAnd ...
func (field Uint64) BitAnd(value uint64) Uint64 {
	return Uint64{field.bitAnd(value)}
}

// BitOr ...
func (field Uint64) BitOr(value uint64) Uint64 {
	return Uint64{field.bitOr(value)}
}

// BitFlip ...
func (field Uint64) BitFlip() Uint64 {
	return Uint64{field.bitFlip()}
}

// Value set value
func (field Uint64) Value(value uint64) AssignExpr {
	return field.value(value)
}

// Zero set zero value
func (field Uint64) Zero() AssignExpr {
	return field.value(0)
}

// Sum ...
func (field Uint64) Sum() Uint64 {
	return Uint64{field.sum()}
}

// IfNull ...
func (field Uint64) IfNull(value uint64) Expr {
	return field.ifNull(value)
}

func (field Uint64) toSlice(values ...uint64) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}
