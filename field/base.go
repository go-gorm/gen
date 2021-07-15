package field

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ Expr = new(Field)

// Expr a query expression about field
type Expr interface {
	clause.Expression

	As(alias string) Expr
	Column() clause.Column
	BuildColumn(*gorm.Statement, ...BuildOpt) string
	RawExpr() interface{}

	// pirvate do nothing, prevent users from implementing interfaces outside the package
	private()
}

func Or(exprs ...Expr) Expr {
	return &expr{expression: clause.Or(toExpression(exprs...)...)}
}

func And(exprs ...Expr) Expr {
	return &expr{expression: clause.And(toExpression(exprs...)...)}
}

func Not(exprs ...Expr) Expr {
	return &expr{expression: clause.Not(toExpression(exprs...)...)}
}

func toExpression(conds ...Expr) []clause.Expression {
	exprs := make([]clause.Expression, len(conds))
	for i, cond := range conds {
		exprs[i] = cond
	}
	return exprs
}

// Field a standard field struct
type Field struct {
	expr
}

type expr struct {
	Col clause.Column

	expression clause.Expression
}

func (expr) private() {}

func (e expr) Build(builder clause.Builder) {
	if e.expression == nil {
		e.expression = clause.NamedExpr{SQL: "?", Vars: []interface{}{e.Col}}
	}
	e.expression.Build(builder)
}

func (e expr) Column() clause.Column {
	return e.Col
}

type BuildOpt func(clause.Column) interface{}

var WithAlias BuildOpt = func(col clause.Column) interface{} {
	return col
}

func (e expr) BuildColumn(stmt *gorm.Statement, opts ...BuildOpt) string {
	var col interface{} = e.Col.Name
	for _, opt := range opts {
		col = opt(e.Col)
	}
	return stmt.Quote(col)
}

func (e expr) RawExpr() interface{} {
	if e.expression == nil {
		return e.Col
	}
	return e.expression
}

// ======================== basic function ========================
func (e expr) Count() Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "COUNT(?)", Vars: []interface{}{e.Col}}}
}

func (e expr) Length() Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "LENGTH(?)", Vars: []interface{}{e.Col}}}
}

func (e expr) Max() Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "MAX(?)", Vars: []interface{}{e.Col}}}
}

func (e expr) Min() Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "MIN(?)", Vars: []interface{}{e.Col}}}
}

func (e expr) Avg() Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "AVG(?)", Vars: []interface{}{e.Col}}}
}

func (e expr) Sum() Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "SUM(?)", Vars: []interface{}{e.Col}}}
}

func (e expr) WithTable(table string) Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "?.?", Vars: []interface{}{clause.Table{Name: table}, e.Col}}}
}

// ======================== comparison between columns ========================
func (e expr) EqCol(col Expr) Expr {
	return &expr{expression: clause.Expr{SQL: "? = ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}}}
}

func (e expr) GtCol(col Expr) Expr {
	return &expr{expression: clause.Expr{SQL: "? > ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}}}
}

func (e expr) GteCol(col Expr) Expr {
	return &expr{expression: clause.Expr{SQL: "? >= ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}}}
}

func (e expr) LtCol(col Expr) Expr {
	return &expr{expression: clause.Expr{SQL: "? < ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}}}
}

func (e expr) LteCol(col Expr) Expr {
	return &expr{expression: clause.Expr{SQL: "? <= ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}}}
}

// ======================== keyword ========================
func (e expr) As(alias string) Expr {
	if e.expression != nil {
		return &expr{Col: e.Col, expression: clause.Expr{SQL: "? AS ?", Vars: []interface{}{e.expression, clause.Column{Name: alias}}}}
	}
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "?", Vars: []interface{}{clause.Column{Name: e.Col.Name, Alias: alias}}}}
}

func (e expr) Desc() Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "? DESC", Vars: []interface{}{e.Col}}}
}

// ======================== general experssion ========================
func (e expr) between(values []interface{}) Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "? BETWEEN ? AND ?", Vars: append([]interface{}{e.Col}, values...)}}
}

func (e expr) add(value interface{}) Expr {
	switch v := value.(type) {
	case time.Duration:
		return &expr{Col: e.Col, expression: clause.Expr{SQL: "DATE_ADD(?, INTERVAL ? MICROSECOND)", Vars: []interface{}{e.Col, v.Microseconds()}}}
	default:
		return &expr{Col: e.Col, expression: clause.Expr{SQL: "?+?", Vars: []interface{}{e.Col, value}}}
	}
}

func (e expr) sub(value interface{}) Expr {
	switch v := value.(type) {
	case time.Duration:
		return &expr{Col: e.Col, expression: clause.Expr{SQL: "DATE_SUB(?, INTERVAL ? MICROSECOND)", Vars: []interface{}{e.Col, v.Microseconds()}}}
	default:
		return &expr{Col: e.Col, expression: clause.Expr{SQL: "?-?", Vars: []interface{}{e.Col, value}}}
	}
}

func (e expr) mul(value interface{}) Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "?*?", Vars: []interface{}{e.Col, value}}}
}

func (e expr) div(value interface{}) Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "?/?", Vars: []interface{}{e.Col, value}}}
}

func (e expr) mod(value interface{}) Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "?%?", Vars: []interface{}{e.Col, value}}}
}

func (e expr) floorDiv(value interface{}) Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "? DIV ?", Vars: []interface{}{e.Col, value}}}
}

func (e expr) rightShift(value interface{}) Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "?>>?", Vars: []interface{}{e.Col, value}}}
}

func (e expr) leftShift(value interface{}) Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "?<<?", Vars: []interface{}{e.Col, value}}}
}

func (e expr) bitXor(value interface{}) Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "?^?", Vars: []interface{}{e.Col, value}}}
}

func (e expr) bitAnd(value interface{}) Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "?&?", Vars: []interface{}{e.Col, value}}}
}

func (e expr) bitOr(value interface{}) Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "?|?", Vars: []interface{}{e.Col, value}}}
}

func (e expr) bitFlip() Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "~?", Vars: []interface{}{e.Col}}}
}

func (e expr) regexp(value interface{}) Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "? REGEXP ?", Vars: []interface{}{e.Col, value}}}
}

func (e expr) not() Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "NOT ?", Vars: []interface{}{e.Col}}}
}

func (e expr) is(value interface{}) Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "? IS ?", Vars: []interface{}{e.Col, value}}}
}

func (e expr) and(value interface{}) Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "? AND ?", Vars: []interface{}{e.Col, value}}}
}

func (e expr) or(value interface{}) Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "? OR ?", Vars: []interface{}{e.Col, value}}}
}

func (e expr) xor(value interface{}) Expr {
	return &expr{Col: e.Col, expression: clause.Expr{SQL: "? XOR ?", Vars: []interface{}{e.Col, value}}}
}

// ======================== subquery method ========================
func ContainsSubQuery(columns []Expr, subQuery *gorm.Statement) Expr {
	switch len(columns) {
	case 0:
		return &expr{expression: clause.Expr{}}
	case 1:
		return &expr{expression: clause.Expr{
			SQL:  fmt.Sprint("?", " IN ", wrap(subQuery.SQL.String())),
			Vars: append([]interface{}{columns[0].RawExpr()}, subQuery.Statement.Vars...),
		}}
	default: // len(columns) > 0
		vars := make([]string, len(columns))
		queryCols := make([]interface{}, len(columns))
		for i, c := range columns {
			vars[i], queryCols[i] = "?", c.RawExpr()
		}
		return &expr{expression: clause.Expr{
			SQL:  fmt.Sprint(wrap(strings.Join(vars, ", ")), " IN ", wrap(subQuery.SQL.String())),
			Vars: append(queryCols, subQuery.Statement.Vars...),
		}}
	}
}

type CompareOperate string

const (
	EqOp  CompareOperate = " = "
	GtOp  CompareOperate = " > "
	GteOp CompareOperate = " >= "
	LtOp  CompareOperate = " < "
	LteOp CompareOperate = " <= "
)

func CompareSubQuery(op CompareOperate, column Expr, subQuery *gorm.Statement) Expr {
	return &expr{expression: clause.Expr{
		SQL:  fmt.Sprint("?", op, wrap(subQuery.SQL.String())),
		Vars: append([]interface{}{column.RawExpr()}, subQuery.Statement.Vars...),
	}}
}

func wrap(subQuery string) string {
	return "(" + subQuery + ")"
}
