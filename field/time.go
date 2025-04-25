package field

import (
	"fmt"
	"time"

	"gorm.io/gorm/clause"
)

// Time time type field
type Time struct {
	genericsField[time.Time]
}

// Between ...
func (field Time) Between(left time.Time, right time.Time) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field Time) NotBetween(left time.Time, right time.Time) Expr {
	return Not(field.Between(left, right))
}

// NotIn ...
func (field Time) NotIn(values ...time.Time) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Add ...
func (field Time) Add(value time.Duration) Time {
	return Time{genericsField[time.Time]{field.add(value)}}
}

// Sub ...
func (field Time) Sub(value time.Duration) Time {
	return Time{genericsField[time.Time]{field.sub(value)}}
}

// Date convert to data, equal to "DATE(time_expr)"
func (field Time) Date() Time {
	return Time{genericsField[time.Time]{expr{e: clause.Expr{SQL: "DATE(?)", Vars: []interface{}{field.RawExpr()}}}}}
}

// DateDiff equal to DATADIFF(self, value)
func (field Time) DateDiff(value time.Time) Number[int] {
	return newNumber[int](expr{e: clause.Expr{SQL: "DATEDIFF(?,?)", Vars: []interface{}{field.RawExpr(), value}}})
}

// DateFormat equal to DATE_FORMAT(self, value)
func (field Time) DateFormat(value string) String {
	return newChars[string](expr{e: clause.Expr{SQL: "DATE_FORMAT(?,?)", Vars: []interface{}{field.RawExpr(), value}}})
}

// Now return result of NOW()
func (field Time) Now() Time {
	return Time{genericsField[time.Time]{expr{e: clause.Expr{SQL: "NOW()"}}}}
}

// CurDate return result of CURDATE()
func (field Time) CurDate() Time {
	return Time{genericsField[time.Time]{expr{e: clause.Expr{SQL: "CURDATE()"}}}}
}

// CurTime return result of CURTIME()
func (field Time) CurTime() Time {
	return Time{genericsField[time.Time]{expr{e: clause.Expr{SQL: "CURTIME()"}}}}
}

// DayName equal to DAYNAME(self)
func (field Time) DayName() String {
	return newChars[string](expr{e: clause.Expr{SQL: "DAYNAME(?)", Vars: []interface{}{field.RawExpr()}}})

}

// MonthName equal to MONTHNAME(self)
func (field Time) MonthName() String {
	return newChars[string](expr{e: clause.Expr{SQL: "MONTHNAME(?)", Vars: []interface{}{field.RawExpr()}}})
}

// Year equal to YEAR(self)
func (field Time) Year() Number[int] {
	return newNumber[int](expr{e: clause.Expr{SQL: "YEAR(?)", Vars: []interface{}{field.RawExpr()}}})
}

// Month equal to MONTH(self)
func (field Time) Month() Number[int] {
	return newNumber[int](expr{e: clause.Expr{SQL: "MONTH(?)", Vars: []interface{}{field.RawExpr()}}})
}

// Day equal to DAY(self)
func (field Time) Day() Number[int] {
	return newNumber[int](expr{e: clause.Expr{SQL: "DAY(?)", Vars: []interface{}{field.RawExpr()}}})
}

// Hour equal to HOUR(self)
func (field Time) Hour() Number[int] {
	return newNumber[int](expr{e: clause.Expr{SQL: "HOUR(?)", Vars: []interface{}{field.RawExpr()}}})
}

// Minute equal to MINUTE(self)
func (field Time) Minute() Number[int] {
	return newNumber[int](expr{e: clause.Expr{SQL: "MINUTE(?)", Vars: []interface{}{field.RawExpr()}}})
}

// Second equal to SECOND(self)
func (field Time) Second() Number[int] {
	return newNumber[int](expr{e: clause.Expr{SQL: "SECOND(?)", Vars: []interface{}{field.RawExpr()}}})
}

// MicroSecond equal to MICROSECOND(self)
func (field Time) MicroSecond() Number[int] {
	return newNumber[int](expr{e: clause.Expr{SQL: "MICROSECOND(?)", Vars: []interface{}{field.RawExpr()}}})
}

// DayOfWeek equal to DAYOFWEEK(self)
func (field Time) DayOfWeek() Number[int] {
	return newNumber[int](expr{e: clause.Expr{SQL: "DAYOFWEEK(?)", Vars: []interface{}{field.RawExpr()}}})
}

// DayOfMonth equal to DAYOFMONTH(self)
func (field Time) DayOfMonth() Number[int] {
	return newNumber[int](expr{e: clause.Expr{SQL: "DAYOFMONTH(?)", Vars: []interface{}{field.RawExpr()}}})
}

// DayOfYear equal to DAYOFYEAR(self)
func (field Time) DayOfYear() Number[int] {
	return newNumber[int](expr{e: clause.Expr{SQL: "DAYOFYEAR(?)", Vars: []interface{}{field.RawExpr()}}})
}

// FromDays equal to FROM_DAYS(self)
func (field Time) FromDays(value int) Time {
	return Time{genericsField[time.Time]{expr{e: clause.Expr{SQL: fmt.Sprintf("FROM_DAYS(%d)", value)}}}}
}

// FromUnixTime equal to FROM_UNIXTIME(self)
func (field Time) FromUnixTime(value int) Time {
	return Time{genericsField[time.Time]{expr{e: clause.Expr{SQL: fmt.Sprintf("FROM_UNIXTIME(%d)", value)}}}}
}

// Zero set zero value
func (field Time) Zero() AssignExpr {
	return field.value(time.Time{})
}
