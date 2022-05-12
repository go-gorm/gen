package field

import (
	"fmt"
	"time"

	"gorm.io/gorm/clause"
)

// Time time type field
type Time Field

// Eq equal to
func (field Time) Eq(value time.Time) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq not equal to
func (field Time) Neq(value time.Time) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// Gt greater than
func (field Time) Gt(value time.Time) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte greater or equal to
func (field Time) Gte(value time.Time) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt less than
func (field Time) Lt(value time.Time) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte less or equal to
func (field Time) Lte(value time.Time) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// Between ...
func (field Time) Between(left time.Time, right time.Time) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field Time) NotBetween(left time.Time, right time.Time) Expr {
	return Not(field.Between(left, right))
}

// In ...
func (field Time) In(values ...time.Time) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// NotIn ...
func (field Time) NotIn(values ...time.Time) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Add ...
func (field Time) Add(value time.Duration) Time {
	return Time{field.add(value)}
}

// Sub ...
func (field Time) Sub(value time.Duration) Time {
	return Time{field.sub(value)}
}

// Date convert to data, equal to "DATE(time_expr)"
func (field Time) Date() Time {
	return Time{expr{e: clause.Expr{SQL: "DATE(?)", Vars: []interface{}{field.RawExpr()}}}}
}

// DateDiff equal to DATADIFF(self, value)
func (field Time) DateDiff(value time.Time) Int {
	return Int{expr{e: clause.Expr{SQL: "DATEDIFF(?,?)", Vars: []interface{}{field.RawExpr(), value}}}}
}

// DateFormat equal to DATE_FORMAT(self, value)
func (field Time) DateFormat(value string) String {
	return String{expr{e: clause.Expr{SQL: "DATE_FORMAT(?,?)", Vars: []interface{}{field.RawExpr(), value}}}}
}

// Now return result of NOW()
func (field Time) Now() Time {
	return Time{expr{e: clause.Expr{SQL: "NOW()"}}}
}

// CurDate return result of CURDATE()
func (field Time) CurDate() Time {
	return Time{expr{e: clause.Expr{SQL: "CURDATE()"}}}
}

// CurTime return result of CURTIME()
func (field Time) CurTime() Time {
	return Time{expr{e: clause.Expr{SQL: "CURTIME()"}}}
}

// DayName equal to DAYNAME(self)
func (field Time) DayName() String {
	return String{expr{e: clause.Expr{SQL: "DAYNAME(?)", Vars: []interface{}{field.RawExpr()}}}}
}

// MonthName equal to MONTHNAME(self)
func (field Time) MonthName() String {
	return String{expr{e: clause.Expr{SQL: "MONTHNAME(?)", Vars: []interface{}{field.RawExpr()}}}}
}

// Month equal to MONTH(self)
func (field Time) Month() Int {
	return Int{expr{e: clause.Expr{SQL: "MONTH(?)", Vars: []interface{}{field.RawExpr()}}}}
}

// Day equal to DAY(self)
func (field Time) Day() Int {
	return Int{expr{e: clause.Expr{SQL: "DAY(?)", Vars: []interface{}{field.RawExpr()}}}}
}

// Hour equal to HOUR(self)
func (field Time) Hour() Int {
	return Int{expr{e: clause.Expr{SQL: "HOUR(?)", Vars: []interface{}{field.RawExpr()}}}}
}

// Minute equal to MINUTE(self)
func (field Time) Minute() Int {
	return Int{expr{e: clause.Expr{SQL: "MINUTE(?)", Vars: []interface{}{field.RawExpr()}}}}
}

// Second equal to SECOND(self)
func (field Time) Second() Int {
	return Int{expr{e: clause.Expr{SQL: "SECOND(?)", Vars: []interface{}{field.RawExpr()}}}}
}

// MicroSecond equal to MICROSECOND(self)
func (field Time) MicroSecond() Int {
	return Int{expr{e: clause.Expr{SQL: "MICROSECOND(?)", Vars: []interface{}{field.RawExpr()}}}}
}

// DayOfWeek equal to DAYOFWEEK(self)
func (field Time) DayOfWeek() Int {
	return Int{expr{e: clause.Expr{SQL: "DAYOFWEEK(?)", Vars: []interface{}{field.RawExpr()}}}}
}

// DayOfMonth equal to DAYOFMONTH(self)
func (field Time) DayOfMonth() Int {
	return Int{expr{e: clause.Expr{SQL: "DAYOFMONTH(?)", Vars: []interface{}{field.RawExpr()}}}}
}

// DayOfYear equal to DAYOFYEAR(self)
func (field Time) DayOfYear() Int {
	return Int{expr{e: clause.Expr{SQL: "DAYOFYEAR(?)", Vars: []interface{}{field.RawExpr()}}}}
}

// FromDays equal to FROM_DAYS(self)
func (field Time) FromDays(value int) Time {
	return Time{expr{e: clause.Expr{SQL: fmt.Sprintf("FROM_DAYS(%d)", value)}}}
}

// FromUnixtime equal to FROM_UNIXTIME(self)
func (field Time) FromUnixtime(value int) Time {
	return Time{expr{e: clause.Expr{SQL: fmt.Sprintf("FROM_UNIXTIME(%d)", value)}}}
}

// Value set value
func (field Time) Value(value time.Time) AssignExpr {
	return field.value(value)
}

// Zero set zero value
func (field Time) Zero() AssignExpr {
	return field.value(time.Time{})
}

// Sum calc sum
func (field Time) Sum() Time {
	return Time{field.sum()}
}

// IfNull ...
func (field Time) IfNull(value Time) Expr {
	return field.ifNull(value)
}

func (field Time) toSlice(values ...time.Time) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}
