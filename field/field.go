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
	return expr{expression: clause.Eq{Column: field.RawExpr(), Value: value}}
}

func (field Field) In(values ...ScanValuer) Expr {
	return expr{expression: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

func (field Field) Gt(value ScanValuer) Expr {
	return expr{expression: clause.Gt{Column: field.RawExpr(), Value: value}}
}

func (field Field) Gte(value ScanValuer) Expr {
	return expr{expression: clause.Gte{Column: field.RawExpr(), Value: value}}
}

func (field Field) Lt(value ScanValuer) Expr {
	return expr{expression: clause.Lt{Column: field.RawExpr(), Value: value}}
}

func (field Field) Lte(value ScanValuer) Expr {
	return expr{expression: clause.Lte{Column: field.RawExpr(), Value: value}}
}

func (field Field) Like(value ScanValuer) Expr {
	return expr{expression: clause.Like{Column: field.RawExpr(), Value: value}}
}

func (field Field) toSlice(values ...ScanValuer) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}
