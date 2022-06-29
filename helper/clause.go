package helper

import (
	"strings"

	"gorm.io/gorm/clause"
)

// Cond ...
type Cond struct {
	Cond   bool
	Result string
}

// IfClause if clause
func IfClause(conds []Cond) string {
	judge := func(c Cond) string {
		if c.Cond {
			return c.Result
		}
		return ""
	}

	clauses := make([]string, len(conds))
	for i, cond := range conds {
		clauses[i] = strings.Trim(judge(cond), " ")
	}
	return " " + strings.Join(clauses, " ")
}

// WhereClause where clause
func WhereClause(conds []string) string {
	return joinClause(conds, "WHERE", whereValue, " ")
}

// SetClause set clause
func SetClause(conds []string) string {
	return joinClause(conds, "SET", setValue, ",")
}

func joinClause(conds []string, keyword string, deal func(string) string, sep string) string {
	clauses := make([]string, len(conds))
	for i, clause := range conds {
		clauses[i] = deal(clause)
	}

	sql := trimAll(strings.Join(clauses, sep))
	if sql != "" {
		sql = " " + keyword + " " + sql
	}
	return sql
}

func trimAll(input string) string {
	return trimRight(trimLeft(input))
}

func trimLeft(input string) string {
	input = strings.TrimSpace(input)
	lowercase := strings.ToLower(input)
	switch {
	case strings.HasPrefix(lowercase, "and "):
		return input[4:]
	case strings.HasPrefix(lowercase, "or "):
		return input[3:]
	case strings.HasPrefix(lowercase, "xor "):
		return input[4:]
	case strings.HasPrefix(lowercase, ","):
		return input[1:]
	default:
		return input
	}
}
func trimRight(input string) string {
	input = strings.TrimSpace(input)
	lowercase := strings.ToLower(input)
	switch {
	case strings.HasSuffix(lowercase, " and"):
		return input[:len(input)-3]
	case strings.HasSuffix(lowercase, " or"):
		return input[:len(input)-2]
	case strings.HasSuffix(lowercase, " xor"):
		return input[:len(input)-3]
	case strings.HasSuffix(lowercase, ","):
		return input[:len(input)-1]
	default:
		return input
	}
}

// whereValue append a new condition with prefix "AND"
func whereValue(value string) string {
	value = strings.Trim(value, " ")
	lowercase := strings.ToLower(value)
	switch {
	case lowercase == "":
		return ""
	case strings.HasPrefix(lowercase, "and "):
		return value
	case strings.HasPrefix(lowercase, "or "):
		return value
	case strings.HasPrefix(lowercase, "xor "):
		return value
	default:
		return "AND " + value
	}
}

func setValue(value string) string {
	return strings.Trim(value, ", ")
}

// JoinWhereBuilder join where builder
func JoinWhereBuilder(src *strings.Builder, whereValue strings.Builder) {
	value := trimAll(whereValue.String())
	if value != "" {
		src.WriteString("WHERE ")
		src.WriteString(value)
		src.WriteString(" ")
	}
}

// JoinSetBuilder join set builder
func JoinSetBuilder(src *strings.Builder, setValue strings.Builder) {
	value := trimAll(setValue.String())
	if value != "" {
		src.WriteString("SET ")
		src.WriteString(value)
		src.WriteString(" ")
	}
}

// JoinTblExpr join clause with table expression(sub query...)
type JoinTblExpr struct {
	clause.Join
	TableExpr clause.Expression
}

// NewJoinTblExpr create new join table expr
func NewJoinTblExpr(join clause.Join, tbExpr clause.Expression) JoinTblExpr {
	return JoinTblExpr{Join: join, TableExpr: tbExpr}
}

// Build ...
func (join JoinTblExpr) Build(builder clause.Builder) {
	if builder == nil {
		return
	}
	if join.Type != "" {
		_, _ = builder.WriteString(string(join.Type))
		_ = builder.WriteByte(' ')
	}
	_, _ = builder.WriteString("JOIN ")
	if join.TableExpr != nil {
		join.TableExpr.Build(builder)
	}

	if len(join.ON.Exprs) > 0 {
		_, _ = builder.WriteString(" ON ")
		join.ON.Build(builder)
	} else if len(join.Using) > 0 {
		_, _ = builder.WriteString(" USING (")
		for idx, c := range join.Using {
			if idx > 0 {
				_ = builder.WriteByte(',')
			}
			builder.WriteQuoted(c)
		}
		_ = builder.WriteByte(')')
	}
}
