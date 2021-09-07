package field

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FieldOption func(clause.Column) clause.Column

// TODO implement validator options

// ======================== generic field =======================

func NewField(table, column string, opts ...FieldOption) Field {
	return Field{expr: expr{Col: toColumn(table, column, opts...)}}
}

// ======================== integer =======================

func NewInt(table, column string, opts ...FieldOption) Int {
	return Int{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewInt8(table, column string, opts ...FieldOption) Int8 {
	return Int8{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewInt16(table, column string, opts ...FieldOption) Int16 {
	return Int16{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewInt32(table, column string, opts ...FieldOption) Int32 {
	return Int32{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewInt64(table, column string, opts ...FieldOption) Int64 {
	return Int64{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewUint(table, column string, opts ...FieldOption) Uint {
	return Uint{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewUint8(table, column string, opts ...FieldOption) Uint8 {
	return Uint8{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewUint16(table, column string, opts ...FieldOption) Uint16 {
	return Uint16{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewUint32(table, column string, opts ...FieldOption) Uint32 {
	return Uint32{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewUint64(table, column string, opts ...FieldOption) Uint64 {
	return Uint64{expr: expr{Col: toColumn(table, column, opts...)}}
}

// ======================== float =======================

func NewFloat32(table, column string, opts ...FieldOption) Float32 {
	return Float32{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewFloat64(table, column string, opts ...FieldOption) Float64 {
	return Float64{expr: expr{Col: toColumn(table, column, opts...)}}
}

// ======================== string =======================

func NewString(table, column string, opts ...FieldOption) String {
	return String{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewBytes(table, column string, opts ...FieldOption) Bytes {
	return Bytes{expr: expr{Col: toColumn(table, column, opts...)}}
}

// ======================== bool =======================

func NewBool(table, column string, opts ...FieldOption) Bool {
	return Bool{expr: expr{Col: toColumn(table, column, opts...)}}
}

// ======================== time =======================

func NewTime(table, column string, opts ...FieldOption) Time {
	return Time{expr: expr{Col: toColumn(table, column, opts...)}}
}

func toColumn(table, column string, opts ...FieldOption) clause.Column {
	col := clause.Column{Table: table, Name: column}
	for _, opt := range opts {
		col = opt(col)
	}
	return col
}

// ======================== boolean operate ========================
func Or(exprs ...Expr) Expr {
	return &expr{e: clause.Or(toExpression(exprs...)...)}
}

func And(exprs ...Expr) Expr {
	return &expr{e: clause.And(toExpression(exprs...)...)}
}

func Not(exprs ...Expr) Expr {
	return &expr{e: clause.Not(toExpression(exprs...)...)}
}

func toExpression(conds ...Expr) []clause.Expression {
	exprs := make([]clause.Expression, len(conds))
	for i, cond := range conds {
		exprs[i] = cond.expression()
	}
	return exprs
}

// ======================== subquery method ========================
func ContainsSubQuery(columns []Expr, subQuery *gorm.DB) Expr {
	switch len(columns) {
	case 0:
		return expr{e: clause.Expr{}}
	case 1:
		return expr{e: clause.Expr{
			SQL:  "? IN (?)",
			Vars: []interface{}{columns[0].RawExpr(), subQuery},
		}}
	default: // len(columns) > 0
		vars := make([]string, len(columns))
		queryCols := make([]interface{}, len(columns))
		for i, c := range columns {
			vars[i], queryCols[i] = "?", c.RawExpr()
		}
		return expr{e: clause.Expr{
			SQL:  fmt.Sprintf("(%s) IN (?)", strings.Join(vars, ", ")),
			Vars: append(queryCols, subQuery),
		}}
	}
}

type CompareOperate string

const (
	EqOp  CompareOperate = " = "
	NeqOp CompareOperate = " <> "
	GtOp  CompareOperate = " > "
	GteOp CompareOperate = " >= "
	LtOp  CompareOperate = " < "
	LteOp CompareOperate = " <= "
)

func CompareSubQuery(op CompareOperate, column Expr, subQuery *gorm.DB) Expr {
	return expr{e: clause.Expr{
		SQL:  fmt.Sprint("?", op, "(?)"),
		Vars: []interface{}{column.RawExpr(), subQuery},
	}}
}

type Value interface{ expr() clause.Expr }

type val struct{ clause.Expr }

func (v *val) expr() clause.Expr { return v.Expr }

func Values(value interface{}) Value {
	return &val{clause.Expr{
		SQL:                "?",
		Vars:               []interface{}{value},
		WithoutParentheses: true,
	}}
}

func ContainsValue(columns []Expr, value Value) Expr {
	switch len(columns) {
	case 0:
		return expr{e: clause.Expr{}}
	case 1:
		return expr{e: clause.Expr{
			SQL:  "? IN (?)",
			Vars: []interface{}{columns[0].RawExpr(), value.expr()},
		}}
	default: // len(columns) > 0
		vars := make([]string, len(columns))
		queryCols := make([]interface{}, len(columns))
		for i, c := range columns {
			vars[i], queryCols[i] = "?", c.RawExpr()
		}
		return expr{e: clause.Expr{
			SQL:  fmt.Sprintf("(%s) IN (?)", strings.Join(vars, ", ")),
			Vars: append(queryCols, value.expr()),
		}}
	}
}

func EmptyExpr() Expr { return expr{e: clause.Expr{}} }
