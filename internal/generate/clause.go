package generate

import (
	"fmt"
	"strings"

	"gorm.io/gen/internal/model"
)

// Clause a symbol of clause, it can be sql condition clause, if clause, where clause, set clause and else clause
type Clause interface {
	// String returns the Go expression used to build this clause.
	String() string
	// Create returns the Go statement that initializes or emits this clause.
	Create() string
}

var (
	_ Clause = new(SQLClause)
	_ Clause = new(IfClause)
	_ Clause = new(ElseClause)
	_ Clause = new(WhereClause)
	_ Clause = new(SetClause)
)

type clause struct {
	VarName string
	Type    model.Status
}

// SQLClause sql condition clause
type SQLClause struct {
	clause
	// Value contains Go string expressions joined into one SQL fragment.
	Value []string
}

// String normalizes Value into a Go expression that ends with one SQL space.
func (s SQLClause) String() string {
	sqlString := strings.Join(s.Value, "+")
	// trim left space
	if strings.HasPrefix(sqlString, "\"") {
		sqlString = `"` + strings.TrimLeft(sqlString, `" `)
	}
	// make sure right has only one space
	if !strings.HasSuffix(sqlString, ` "`) {
		sqlString += `+" "`
	}
	// Remove redundant connection symbols
	return strings.ReplaceAll(sqlString, `"+"`, "")
}

// Create create clause
func (s SQLClause) Create() string {
	return fmt.Sprintf("%s.WriteString(%s)", s.VarName, s.String())
}

// Finish finish clause
func (s SQLClause) Finish() string {
	return fmt.Sprintf("%s.WriteString(%s)", s.VarName, s.String())
}

// IfClause if clause
type IfClause struct {
	clause
	// Value contains the clauses nested inside the conditional block.
	Value []Clause
	slice section
}

// String returns the generated Go condition.
func (i IfClause) String() string {
	return i.slice.Value
}

// Create create clause
func (i IfClause) Create() string {
	return fmt.Sprintf("%s {", i.String())
}

// Finish finish clause
func (i IfClause) Finish() string {
	return "}"
}

// ElseClause else clause
type ElseClause struct {
	IfClause
}

// String returns the generated else or else-if expression.
func (e ElseClause) String() (res string) {
	return e.slice.Value
}

// Create create clause
func (e ElseClause) Create() string {
	return fmt.Sprintf("} %s {", e.String())
}

// Finish finish clause
func (e ElseClause) Finish() string {
	return ""
}

// WhereClause where clause
type WhereClause struct {
	clause
	// Value contains the clauses collected into the WHERE builder.
	Value []Clause
}

// String returns the expression that trims the accumulated WHERE fragment.
func (w WhereClause) String() string {
	return fmt.Sprintf("helper.WhereTrim(%s.String())", w.VarName)
}

// Create create clause
func (w WhereClause) Create() string {
	return fmt.Sprintf("var %s strings.Builder", w.VarName)
}

// Finish finish clause
func (w WhereClause) Finish(name string) string {
	return fmt.Sprintf("helper.JoinWhereBuilder(&%s,%s)", name, w.VarName)
}

// SetClause set clause
type SetClause struct {
	clause
	// Value contains the clauses collected into the SET builder.
	Value []Clause
}

// String returns the expression that trims the accumulated SET fragment.
func (s SetClause) String() string {
	return fmt.Sprintf("helper.SetTrim(%s.String())", s.VarName)
}

// Create create clause
func (s SetClause) Create() string {
	return fmt.Sprintf("var %s strings.Builder", s.VarName)
}

// Finish finish clause
func (s SetClause) Finish(name string) string {
	return fmt.Sprintf("helper.JoinSetBuilder(&%s,%s)", name, s.VarName)
}

// TrimClause set clause
type TrimClause struct {
	clause
	// Value contains the clauses whose edge separators are trimmed.
	Value []Clause
}

// String returns the expression that trims logical operators and separators.
func (s TrimClause) String() string {
	return fmt.Sprintf("helper.TrimALL(%s.String())", s.VarName)
}

// Create create trim clause
func (s TrimClause) Create() string {
	return fmt.Sprintf("var %s strings.Builder", s.VarName)
}

// Finish finish trim clause
func (s TrimClause) Finish(name string) string {
	return fmt.Sprintf("helper.JoinTrimAllBuilder(&%s,%s)", name, s.VarName)
}

// ForClause set clause
type ForClause struct {
	clause
	// Value contains the clauses emitted for each iteration.
	Value []Clause
	// ForRange describes the generated Go range statement.
	ForRange ForRange
	forSlice section
}

// String returns the generated Go range statement and opening brace.
func (f ForClause) String() string {
	return f.forSlice.Value + "{"
}

// Create create clause
func (f ForClause) Create() string {
	return f.String()
}

// Finish finish clause
func (f ForClause) Finish() string {
	return "}"
}
