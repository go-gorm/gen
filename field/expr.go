package field

import (
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
func (e expr) IsNull() Expr {
	return e.setExpression(clause.Expr{SQL: "? IS NULL", Vars: []interface{}{e.RawExpr()}})
}

func (e expr) IsNotNull() Expr {
	return e.setExpression(clause.Expr{SQL: "? IS NOT NULL", Vars: []interface{}{e.RawExpr()}})
}

func (e expr) Count() Int {
	return Int{e.setExpression(clause.Expr{SQL: "COUNT(?)", Vars: []interface{}{e.RawExpr()}})}
}

func (e expr) Length() Int {
	return Int{e.setExpression(clause.Expr{SQL: "LENGTH(?)", Vars: []interface{}{e.RawExpr()}})}
}

func (e expr) Max() Float64 {
	return Float64{e.setExpression(clause.Expr{SQL: "MAX(?)", Vars: []interface{}{e.RawExpr()}})}
}

func (e expr) Min() Float64 {
	return Float64{e.setExpression(clause.Expr{SQL: "MIN(?)", Vars: []interface{}{e.RawExpr()}})}
}

func (e expr) Avg() Float64 {
	return Float64{e.setExpression(clause.Expr{SQL: "AVG(?)", Vars: []interface{}{e.RawExpr()}})}
}

func (e expr) Sum() Float64 {
	return Float64{e.setExpression(clause.Expr{SQL: "SUM(?)", Vars: []interface{}{e.RawExpr()}})}
}

func (e expr) WithTable(table string) Expr {
	e.Col.Table = table
	return e
}

// ======================== comparison between columns ========================
func (e expr) EqCol(col Expr) Expr {
	return e.setExpression(clause.Expr{SQL: "? = ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}})
}

func (e expr) NeqCol(col Expr) Expr {
	return e.setExpression(clause.Expr{SQL: "? <> ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}})
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
	return e.setExpression(clause.Expr{SQL: "? DESC", Vars: []interface{}{e.RawExpr()}})
}

// ======================== general experssion ========================
func (e expr) between(values []interface{}) expr {
	return e.setExpression(clause.Expr{SQL: "? BETWEEN ? AND ?", Vars: append([]interface{}{e.RawExpr()}, values...)})
}

func (e expr) add(value interface{}) expr {
	switch v := value.(type) {
	case time.Duration:
		return e.setExpression(clause.Expr{SQL: "DATE_ADD(?, INTERVAL ? MICROSECOND)", Vars: []interface{}{e.RawExpr(), v.Microseconds()}})
	default:
		return e.setExpression(clause.Expr{SQL: "?+?", Vars: []interface{}{e.RawExpr(), value}})
	}
}

func (e expr) sub(value interface{}) expr {
	switch v := value.(type) {
	case time.Duration:
		return e.setExpression(clause.Expr{SQL: "DATE_SUB(?, INTERVAL ? MICROSECOND)", Vars: []interface{}{e.RawExpr(), v.Microseconds()}})
	default:
		return e.setExpression(clause.Expr{SQL: "?-?", Vars: []interface{}{e.RawExpr(), value}})
	}
}

func (e expr) mul(value interface{}) expr {
	return e.setExpression(clause.Expr{SQL: "?*?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) div(value interface{}) expr {
	return e.setExpression(clause.Expr{SQL: "?/?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) mod(value interface{}) expr {
	return e.setExpression(clause.Expr{SQL: "?%?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) floorDiv(value interface{}) expr {
	return e.setExpression(clause.Expr{SQL: "? DIV ?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) floor() expr {
	return e.setExpression(clause.Expr{SQL: "FLOOR(?)", Vars: []interface{}{e.RawExpr()}})
}

func (e expr) rightShift(value interface{}) expr {
	return e.setExpression(clause.Expr{SQL: "?>>?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) leftShift(value interface{}) expr {
	return e.setExpression(clause.Expr{SQL: "?<<?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) bitXor(value interface{}) expr {
	return e.setExpression(clause.Expr{SQL: "?^?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) bitAnd(value interface{}) expr {
	return e.setExpression(clause.Expr{SQL: "?&?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) bitOr(value interface{}) expr {
	return e.setExpression(clause.Expr{SQL: "?|?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) bitFlip() expr {
	return e.setExpression(clause.Expr{SQL: "~?", Vars: []interface{}{e.RawExpr()}})
}

func (e expr) regexp(value interface{}) expr {
	return e.setExpression(clause.Expr{SQL: "? REGEXP ?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) not() expr {
	return e.setExpression(clause.Expr{SQL: "NOT ?", Vars: []interface{}{e.RawExpr()}})
}

func (e expr) is(value interface{}) expr {
	return e.setExpression(clause.Expr{SQL: "? IS ?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) and(value interface{}) expr {
	return e.setExpression(clause.Expr{SQL: "? AND ?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) or(value interface{}) expr {
	return e.setExpression(clause.Expr{SQL: "? OR ?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) xor(value interface{}) expr {
	return e.setExpression(clause.Expr{SQL: "? XOR ?", Vars: []interface{}{e.RawExpr(), value}})
}
