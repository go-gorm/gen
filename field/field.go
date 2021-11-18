package field

import (
	"database/sql/driver"

	"gorm.io/gorm/clause"
)

type ScanValuer interface {
	Scan(src interface{}) error   // sql.Scanner
	Value() (driver.Value, error) // driver.Valuer
}

// Field a standard field struct
type Field struct{ expr }

func (field Field) Eq(value ScanValuer) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

func (field Field) Neq(value ScanValuer) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

func (field Field) In(values ...ScanValuer) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

func (field Field) Gt(value ScanValuer) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

func (field Field) Gte(value ScanValuer) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

func (field Field) Lt(value ScanValuer) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

func (field Field) Lte(value ScanValuer) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

func (field Field) Like(value ScanValuer) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

func (field Field) Value(value ScanValuer) AssignExpr {
	return field.value(value)
}

func (field Field) Sum() Field {
	return Field{field.sum()}
}

func (field Field) IfNull(value ScanValuer) Expr {
	return field.ifNull(value)
}

func (field Field) toSlice(values ...ScanValuer) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}
