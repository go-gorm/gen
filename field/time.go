package field

import (
	"fmt"
	"time"

	"gorm.io/gorm/clause"
)

type Time Field

func (field Time) Eq(value time.Time) Expr {
	return expr{expression: clause.Eq{Column: field.RawExpr(), Value: value}}
}

func (field Time) Neq(value time.Time) Expr {
	return expr{expression: clause.Neq{Column: field.RawExpr(), Value: value}}
}

func (field Time) Gt(value time.Time) Expr {
	return expr{expression: clause.Gt{Column: field.RawExpr(), Value: value}}
}

func (field Time) Gte(value time.Time) Expr {
	return expr{expression: clause.Gte{Column: field.RawExpr(), Value: value}}
}

func (field Time) Lt(value time.Time) Expr {
	return expr{expression: clause.Lt{Column: field.RawExpr(), Value: value}}
}

func (field Time) Lte(value time.Time) Expr {
	return expr{expression: clause.Lte{Column: field.RawExpr(), Value: value}}
}

func (field Time) Between(left time.Time, right time.Time) Expr {
	return field.between([]interface{}{left, right})
}

func (field Time) NotBetween(left time.Time, right time.Time) Expr {
	return Not(field.Between(left, right))
}

func (field Time) In(values ...time.Time) Expr {
	return expr{expression: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

func (field Time) NotIn(values ...time.Time) Expr {
	return expr{expression: clause.Not(field.In(values...))}
}

func (field Time) Add(value time.Duration) Time {
	return Time{field.add(value)}
}

func (field Time) Sub(value time.Duration) Time {
	return Time{field.sub(value)}
}

func (field Time) Date() Time {
	return Time{expr{expression: clause.Expr{SQL: "DATE(?)", Vars: []interface{}{field.RawExpr()}}}}
}

func (field Time) DateDiff(value time.Time) Int {
	return Int{expr{expression: clause.Expr{SQL: "DATEDIFF(?,?)", Vars: []interface{}{field.RawExpr(), value}}}}
}

func (field Time) DateFormat(value string) String {
	return String{expr{expression: clause.Expr{SQL: "DATE_FORMAT(?,?)", Vars: []interface{}{field.RawExpr(), value}}}}
}

func (field Time) Now() Time {
	return Time{expr{expression: clause.Expr{SQL: "NOW()"}}}
}

func (field Time) CurDate() Time {
	return Time{expr{expression: clause.Expr{SQL: "CURDATE()"}}}
}

func (field Time) CurTime() Time {
	return Time{expr{expression: clause.Expr{SQL: "CURTIME()"}}}
}

func (field Time) DayName() String {
	return String{expr{expression: clause.Expr{SQL: "DAYNAME(?)", Vars: []interface{}{field.RawExpr()}}}}
}

func (field Time) MonthName() String {
	return String{expr{expression: clause.Expr{SQL: "MONTHNAME(?)", Vars: []interface{}{field.RawExpr()}}}}
}

func (field Time) Month() Int {
	return Int{expr{expression: clause.Expr{SQL: "MONTH(?)", Vars: []interface{}{field.RawExpr()}}}}
}

func (field Time) Day() Int {
	return Int{expr{expression: clause.Expr{SQL: "DAY(?)", Vars: []interface{}{field.RawExpr()}}}}
}

func (field Time) Hour() Int {
	return Int{expr{expression: clause.Expr{SQL: "HOUR(?)", Vars: []interface{}{field.RawExpr()}}}}
}

func (field Time) Minute() Int {
	return Int{expr{expression: clause.Expr{SQL: "MINUTE(?)", Vars: []interface{}{field.RawExpr()}}}}
}

func (field Time) Second() Int {
	return Int{expr{expression: clause.Expr{SQL: "SECOND(?)", Vars: []interface{}{field.RawExpr()}}}}
}

func (field Time) MicroSecond() Int {
	return Int{expr{expression: clause.Expr{SQL: "MICROSECOND(?)", Vars: []interface{}{field.RawExpr()}}}}
}

func (field Time) DayOfWeek() Int {
	return Int{expr{expression: clause.Expr{SQL: "DAYOFWEEK(?)", Vars: []interface{}{field.RawExpr()}}}}
}

func (field Time) DayOfMonth() Int {
	return Int{expr{expression: clause.Expr{SQL: "DAYOFMONTH(?)", Vars: []interface{}{field.RawExpr()}}}}
}

func (field Time) DayOfYear() Int {
	return Int{expr{expression: clause.Expr{SQL: "DAYOFYEAR(?)", Vars: []interface{}{field.RawExpr()}}}}
}

func (field Time) FromDays(value int) Time {
	return Time{expr{expression: clause.Expr{SQL: fmt.Sprintf("FROM_DAYS(%d)", value)}}}
}

func (field Time) FromUnixtime(value int) Time {
	return Time{expr{expression: clause.Expr{SQL: fmt.Sprintf("FROM_UNIXTIME(%d)", value)}}}
}

func (field Time) toSlice(values ...time.Time) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}
