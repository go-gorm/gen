package generate

import (
	"fmt"
	"strings"

	"gorm.io/gen/internal/model"
)

// Clause a symbol of clause, it can be sql condition clause, if clause, where clause, set clause and else clause
type Clause interface {
	String() string
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
	Value []string
}

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
	Value []Clause
	slice section
}

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
	Value []Clause
}

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
	Value []Clause
}

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

// ForClause set clause
type ForClause struct {
	clause
	Value    []Clause
	ForRange ForRange
	forSlice section
}

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
