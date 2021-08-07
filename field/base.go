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
	BuildExpr(stmt *gorm.Statement) string

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

var (
	// WithTable build column with table
	WithTable BuildOpt = func(col clause.Column) interface{} { return clause.Column{Table: col.Table, Name: col.Name} }

	// WithAll build column with table and alias
	WithAll BuildOpt = func(col clause.Column) interface{} { return col }
)

func (e expr) BuildColumn(stmt *gorm.Statement, opts ...BuildOpt) string {
	var col interface{} = e.Col.Name
	for _, opt := range opts {
		col = opt(e.Col)
	}
	return stmt.Quote(col)
}

func (e expr) BuildExpr(stmt *gorm.Statement) string {
	if e.expression == nil {
		return e.BuildColumn(stmt, WithAll)
	}
	newStmt := &gorm.Statement{DB: stmt.DB, Table: stmt.Table, Schema: stmt.Schema}
	e.expression.Build(newStmt)
	return newStmt.SQL.String()
}

func (e expr) RawExpr() interface{} {
	if e.expression == nil {
		return e.Col
	}
	return e.expression
}

func (e expr) setExpression(expression clause.Expression) expr {
	e.expression = expression
	return e
}

// ======================== basic function ========================
func (e expr) Count() Expr {
	return e.setExpression(clause.Expr{SQL: "COUNT(?)", Vars: []interface{}{e.Col}})
}

func (e expr) Length() Expr {
	return e.setExpression(clause.Expr{SQL: "LENGTH(?)", Vars: []interface{}{e.Col}})
}

func (e expr) Max() Expr {
	return e.setExpression(clause.Expr{SQL: "MAX(?)", Vars: []interface{}{e.Col}})
}

func (e expr) Min() Expr {
	return e.setExpression(clause.Expr{SQL: "MIN(?)", Vars: []interface{}{e.Col}})
}

func (e expr) Avg() Expr {
	return e.setExpression(clause.Expr{SQL: "AVG(?)", Vars: []interface{}{e.Col}})
}

func (e expr) Sum() Expr {
	return e.setExpression(clause.Expr{SQL: "SUM(?)", Vars: []interface{}{e.Col}})
}

func (e expr) WithTable(table string) Expr {
	e.Col.Table = table
	return e
}

// ======================== comparison between columns ========================
func (e expr) EqCol(col Expr) Expr {
	return e.setExpression(clause.Expr{SQL: "? = ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}})
}

func (e expr) GtCol(col Expr) Expr {
	return e.setExpression(clause.Expr{SQL: "? > ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}})
}

func (e expr) GteCol(col Expr) Expr {
	return e.setExpression(clause.Expr{SQL: "? >= ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}})
}

func (e expr) LtCol(col Expr) Expr {
	return e.setExpression(clause.Expr{SQL: "? < ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}})
}

func (e expr) LteCol(col Expr) Expr {
	return e.setExpression(clause.Expr{SQL: "? <= ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}})
}

// ======================== keyword ========================
func (e expr) As(alias string) Expr {
	if e.expression != nil {
		return e.setExpression(clause.Expr{SQL: "? AS ?", Vars: []interface{}{e.expression, clause.Column{Name: alias}}})
	}
	e.Col.Alias = alias
	return e
}

func (e expr) Desc() Expr {
	return e.setExpression(clause.Expr{SQL: "? DESC", Vars: []interface{}{e.Col}})
}

// ======================== general experssion ========================
func (e expr) between(values []interface{}) Expr {
	return e.setExpression(clause.Expr{SQL: "? BETWEEN ? AND ?", Vars: append([]interface{}{e.Col}, values...)})
}

func (e expr) add(value interface{}) Expr {
	switch v := value.(type) {
	case time.Duration:
		return e.setExpression(clause.Expr{SQL: "DATE_ADD(?, INTERVAL ? MICROSECOND)", Vars: []interface{}{e.Col, v.Microseconds()}})
	default:
		return e.setExpression(clause.Expr{SQL: "?+?", Vars: []interface{}{e.Col, value}})
	}
}

func (e expr) sub(value interface{}) Expr {
	switch v := value.(type) {
	case time.Duration:
		return e.setExpression(clause.Expr{SQL: "DATE_SUB(?, INTERVAL ? MICROSECOND)", Vars: []interface{}{e.Col, v.Microseconds()}})
	default:
		return e.setExpression(clause.Expr{SQL: "?-?", Vars: []interface{}{e.Col, value}})
	}
}

func (e expr) mul(value interface{}) Expr {
	return e.setExpression(clause.Expr{SQL: "?*?", Vars: []interface{}{e.Col, value}})
}

func (e expr) div(value interface{}) Expr {
	return e.setExpression(clause.Expr{SQL: "?/?", Vars: []interface{}{e.Col, value}})
}

func (e expr) mod(value interface{}) Expr {
	return e.setExpression(clause.Expr{SQL: "?%?", Vars: []interface{}{e.Col, value}})
}

func (e expr) floorDiv(value interface{}) Expr {
	return e.setExpression(clause.Expr{SQL: "? DIV ?", Vars: []interface{}{e.Col, value}})
}

func (e expr) rightShift(value interface{}) Expr {
	return e.setExpression(clause.Expr{SQL: "?>>?", Vars: []interface{}{e.Col, value}})
}

func (e expr) leftShift(value interface{}) Expr {
	return e.setExpression(clause.Expr{SQL: "?<<?", Vars: []interface{}{e.Col, value}})
}

func (e expr) bitXor(value interface{}) Expr {
	return e.setExpression(clause.Expr{SQL: "?^?", Vars: []interface{}{e.Col, value}})
}

func (e expr) bitAnd(value interface{}) Expr {
	return e.setExpression(clause.Expr{SQL: "?&?", Vars: []interface{}{e.Col, value}})
}

func (e expr) bitOr(value interface{}) Expr {
	return e.setExpression(clause.Expr{SQL: "?|?", Vars: []interface{}{e.Col, value}})
}

func (e expr) bitFlip() Expr {
	return e.setExpression(clause.Expr{SQL: "~?", Vars: []interface{}{e.Col}})
}

func (e expr) regexp(value interface{}) Expr {
	return e.setExpression(clause.Expr{SQL: "? REGEXP ?", Vars: []interface{}{e.Col, value}})
}

func (e expr) not() Expr {
	return e.setExpression(clause.Expr{SQL: "NOT ?", Vars: []interface{}{e.Col}})
}

func (e expr) is(value interface{}) Expr {
	return e.setExpression(clause.Expr{SQL: "? IS ?", Vars: []interface{}{e.Col, value}})
}

func (e expr) and(value interface{}) Expr {
	return e.setExpression(clause.Expr{SQL: "? AND ?", Vars: []interface{}{e.Col, value}})
}

func (e expr) or(value interface{}) Expr {
	return e.setExpression(clause.Expr{SQL: "? OR ?", Vars: []interface{}{e.Col, value}})
}

func (e expr) xor(value interface{}) Expr {
	return e.setExpression(clause.Expr{SQL: "? XOR ?", Vars: []interface{}{e.Col, value}})
}

// ======================== subquery method ========================
func ContainsSubQuery(columns []Expr, subQuery *gorm.DB) Expr {
	switch len(columns) {
	case 0:
		return expr{expression: clause.Expr{}}
	case 1:
		return expr{expression: clause.Expr{
			SQL:  "? IN (?)",
			Vars: append([]interface{}{columns[0].RawExpr()}, subQuery),
		}}
	default: // len(columns) > 0
		vars := make([]string, len(columns))
		queryCols := make([]interface{}, len(columns))
		for i, c := range columns {
			vars[i], queryCols[i] = "?", c.RawExpr()
		}
		return expr{expression: clause.Expr{
			SQL:  fmt.Sprintf("(%s) IN (?)", strings.Join(vars, ", ")),
			Vars: append(queryCols, subQuery),
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

func CompareSubQuery(op CompareOperate, column Expr, subQuery *gorm.DB) Expr {
	return expr{expression: clause.Expr{
		SQL:  fmt.Sprint("?", op, "(?)"),
		Vars: append([]interface{}{column.RawExpr()}, subQuery),
	}}
}
