package helper

import "strings"

type Cond struct {
	Cond   bool
	Result string
}

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

func WhereClause(conds []string) string {
	return joinClause(conds, "WHERE", whereValue, " ")
}

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
	case strings.HasPrefix(lowercase, "or"):
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

func JoinBuilder(src *strings.Builder, conjunction string, joinValue strings.Builder) {
	value := trimAll(joinValue.String())
	if value != "" {
		src.WriteString(conjunction)
		src.WriteString(value)
	}
}

func JoinBuilder2(src *strings.Builder, conjunction string, joinValues ...strings.Builder) {
	for _, value := range joinValues {
		trimValue := trimAll(value.String())
		if trimValue != "" {
			src.WriteString(conjunction)
			src.WriteString(trimValue)
		}
	}
}
