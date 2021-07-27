package check

import (
	"bytes"
	"fmt"
	"strings"

	"gorm.io/gen/internal/parser"
)

// Clause a symbol of clause, it can be sql condition clause, if clause, where clause, set clause and else cluase
type Clause interface {
	String() string
}

// SQLClause sql condition clause
type SQLClause struct {
	VarName string
	Type    Status
	Value   []string
}

func (s SQLClause) String() string {
	return strings.ReplaceAll(strings.Join(s.Value, "+"), `"+"`, "")
}

// IfClause if clause
type IfClause struct {
	VarName string
	Cond    string
	Type    Status
	Value   []Clause
	Else    []Clause
}

func (i IfClause) String() string {
	return fmt.Sprintf("helper.IfClause(%s)", i.VarName)
}

// ElseClause else clause
type ElseClause struct {
	VarName string
	Cond    string
	Type    Status
	Value   []Clause
}

func (e ElseClause) String() (res string) {
	strList := make([]string, len(e.Value))
	for i, v := range e.Value {
		strList[i] = v.String()
	}
	return strings.ReplaceAll(strings.Join(strList, "+"), `"+"`, "")
}

// WhereClause where clause
type WhereClause struct {
	VarName string
	Type    Status
	Value   []Clause
}

func (w WhereClause) String() string {
	return fmt.Sprintf("helper.WhereClause(%s)", w.VarName)
}

// SetClause set clause
type SetClause struct {
	VarName string
	Type    Status
	Value   []Clause
}

func (w SetClause) String() string {
	return fmt.Sprintf("helper.SetClause(%s)", w.VarName)
}

// Slices split sql into chunks
type Slices struct {
	slices       []slice
	tmpl         []string
	currentIndex int
	Names        map[Status]int
}

// NewSlices create and initialize Slices
func NewSlices() *Slices {
	return &Slices{
		Names: map[Status]int{
			IF:    0,
			WHERE: 0,
			SET:   0,
		},
	}
}

// Next return next slice and increase index by 1
func (s *Slices) Next() slice {
	s.currentIndex++
	return s.slices[s.currentIndex]
}

// SubIndex take index one step back
func (s *Slices) SubIndex() {
	s.currentIndex--
}

// HasMore whether has more slice
func (s *Slices) HasMore() bool {
	return s.currentIndex < len(s.slices)-1
}

// IsNull whether slice is empty
func (s *Slices) IsNull() bool {
	return len(s.slices) == 0
}

// Current return current slice
func (s *Slices) Current() slice {
	return s.slices[s.currentIndex]
}

// GetName ...
func (s *Slices) GetName(status Status) string {
	switch status {
	case IF:
		defer func() { s.Names[IF]++ }()
		return fmt.Sprintf("ifCond%d", s.Names[IF])
	case WHERE:
		defer func() { s.Names[WHERE]++ }()
		return fmt.Sprintf("whereCond%d", s.Names[WHERE])
	case SET:
		defer func() { s.Names[SET]++ }()
		return fmt.Sprintf("setCond%d", s.Names[SET])
	default:
		return fmt.Sprintf("Cond%d", s.currentIndex)
	}
}

func (s *Slices) appendIfCond(name, cond, result string) {
	s.tmpl = append(s.tmpl, fmt.Sprintf("%s = append(%s, helper.Cond{%s, %s})", name, name, cond, result))
}

func (s *Slices) appendSetValue(name, result string) {
	s.tmpl = append(s.tmpl, fmt.Sprintf("%s = append(%s,  %s)", name, name, strings.TrimSpace(result)))
}

// CreateIf create if clause code
func (s *Slices) CreateIf(name string) {
	s.tmpl = append(s.tmpl, fmt.Sprintf("%s := make([]helper.Cond, 0, 100)", name))
}

// CreateStringSet create string set
func (s *Slices) CreateStringSet(name string) {
	s.tmpl = append(s.tmpl, fmt.Sprintf("%s := make([]string, 0, 100)", name))
}

// parse slice and append result to tmpl, return a Clause array
func (s *Slices) parse() (res []Clause, err error) {
	if s.IsNull() {
		err = fmt.Errorf("sql is null")
		return
	}
	name := "generateSQL"
	for slice := s.Current(); ; slice = s.Next() {
		s.tmpl = append(s.tmpl, "")
		switch slice.Type {
		case SQL, DATA, VARIABLE:
			sqlClause := s.parseSQL(name)
			res = append(res, sqlClause)
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=%s", name, sqlClause.String()))
		case IF:
			var ifClause IfClause
			ifClause, err = s.parseIF()
			if err != nil {
				return
			}
			res = append(res, ifClause)
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=helper.IfClause(%s)", name, ifClause.VarName))
		case WHERE:
			var whereClause WhereClause
			whereClause, err = s.parseWhere()
			if err != nil {
				return
			}
			res = append(res, whereClause)
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=helper.WhereClause(%s)", name, whereClause.VarName))
		case SET:
			var setClause SetClause
			setClause, err = s.parseSet()
			if err != nil {
				return
			}
			res = append(res, setClause)
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=helper.SetClause(%s)", name, setClause.VarName))
		case END:
		default:
			err = fmt.Errorf("unknow clause:%s", slice.Origin)
			return
		}

		if !s.HasMore() {
			return
		}
	}
}

// parseIF parse if clause
func (s *Slices) parseIF() (res IfClause, err error) {
	slice := s.Current()
	name := s.GetName(slice.Type)
	s.CreateIf(name)

	res.Type = slice.Type
	res.Cond = slice.Value
	res.VarName = name
	cond := []string{res.Cond}
	for s.HasMore() {
		n := s.Next()
		switch n.Type {
		case SQL, DATA, VARIABLE:
			str := s.parseSQL(name)
			res.Value = append(res.Value, str)
			s.appendIfCond(name, res.Cond, str.String())
		case IF:
			var ifClause IfClause
			ifClause, err = s.parseIF()
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendIfCond(name, res.Cond, ifClause.String())
		case WHERE:
			var whereClause WhereClause
			whereClause, err = s.parseWhere()
			if err != nil {
				return
			}
			res.Value = append(res.Value, whereClause)
			s.appendIfCond(name, res.Cond, whereClause.String())
		case SET:
			var setClause SetClause
			setClause, err = s.parseSet()
			if err != nil {
				return
			}
			res.Value = append(res.Value, setClause)
			s.appendIfCond(name, res.Cond, setClause.String())
		case ELSEIF:
			elseClause := s.parseElSE(name)
			elseCond := elseClause.Cond
			elseClause.Cond = fmt.Sprintf("!(%s) && %s", strings.Join(cond, " || "), elseCond)
			res.Else = append(res.Else, elseClause)
			s.appendIfCond(name, elseClause.Cond, elseClause.String())
			cond = append(cond, elseCond)
		case ELSE:
			elseClause := s.parseElSE(name)
			elseClause.Cond = fmt.Sprintf("!(%s)", strings.Join(cond, " || "))
			res.Else = append(res.Else, elseClause)
			s.appendIfCond(name, elseClause.Cond, elseClause.String())
		case END:
			return
		default:
			err = fmt.Errorf("unknow clause : %s", n.Origin)
			return
		}
	}
	err = fmt.Errorf("incomplete SQL,if not end")
	return
}

// parseElSE parse else clause, the clause' type must be one of if, where, set, SQL condition
func (s *Slices) parseElSE(name string) (res ElseClause) {
	slice := s.Current()
	res.Cond = slice.Value
	res.VarName = name
	res.Type = slice.Type
	for n := s.Next(); s.HasMore(); n = s.Next() {
		switch n.Type {
		case SQL, DATA, VARIABLE:
			res.Value = append(res.Value, s.parseSQL(name))
		case IF:
			ifClause, err := s.parseIF()
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
		case WHERE:
			whereClause, err := s.parseWhere()
			if err != nil {
				return
			}
			res.Value = append(res.Value, whereClause)
		case SET:
			setClause, err := s.parseSet()
			if err != nil {
				return
			}
			res.Value = append(res.Value, setClause)
		default:
			s.SubIndex()
			return
		}
	}
	return
}

// parseWhere parse where clause, the clause' type must be one of if, SQL condition
func (s *Slices) parseWhere() (res WhereClause, err error) {
	slice := s.Current()
	name := s.GetName(slice.Type)
	s.CreateStringSet(name)

	res.VarName = name
	res.Type = slice.Type
	for s.HasMore() {
		n := s.Next()
		switch n.Type {
		case SQL, DATA, VARIABLE:
			strClause := s.parseSQL(name)
			res.Value = append(res.Value, strClause)
			s.appendSetValue(name, strClause.String())
		case IF:
			var ifClause IfClause
			ifClause, err = s.parseIF()
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendSetValue(name, ifClause.String())
		case END:
			return
		default:
			err = fmt.Errorf("unknow clause : %s", n.Origin)
			return
		}
	}
	err = fmt.Errorf("incomplete SQL,where not end")
	return
}

// parseSet parse set clause, the clause' type must be one of if, SQL condition
func (s *Slices) parseSet() (res SetClause, err error) {
	slice := s.Current()
	name := s.GetName(slice.Type)
	s.CreateStringSet(name)

	res.VarName = name
	res.Type = slice.Type
	for s.HasMore() {
		n := s.Next()
		switch n.Type {
		case SQL, DATA, VARIABLE:
			strClause := s.parseSQL(name)
			res.Value = append(res.Value, strClause)
			s.appendSetValue(name, strClause.String())
		case IF:
			var ifClause IfClause
			ifClause, err = s.parseIF()
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendSetValue(name, ifClause.String())
		case END:
			return
		default:
			err = fmt.Errorf("unknow clause : %s", n.Origin)
			return
		}
	}
	err = fmt.Errorf("incomplete SQL,set not end")
	return
}

// parseSQL parse sql condition, the clause' type must be one of SQL condition, VARIABLE, Data
func (s *Slices) parseSQL(name string) (res SQLClause) {
	res.VarName = name
	res.Type = SQL
	for slice := s.Current(); ; slice = s.Next() {
		switch slice.Type {
		case SQL, VARIABLE, DATA:
			res.Value = append(res.Value, slice.Value)
		default:
			s.SubIndex()
			return
		}

		if !s.HasMore() {
			return
		}
	}
}

// sql fragment
type fragment struct {
	Type  Status
	value string
}

func checkFragment(s string, params []parser.Param) (f fragment, err error) {
	f = fragment{Type: UNKNOWN, value: strings.Trim(s, " ")}
	str := strings.ToLower(strings.Trim(s, " "))
	switch str {
	case "&&", "||":
		f.Type = LOGICAL
		return
	case ">", "<", ">=", "<=", "==", "!=":
		f.Type = EXPRESSION
		return
	case "end":
		f.Type = END
		return
	case "if":
		f.Type = IF
		return
	case "set":
		f.Type = SET
		return
	case "else":
		f.Type = ELSE
		return
	case "where":
		f.Type = WHERE
		return
	case "true", "false":
		f.Type = BOOL
		return
	case "nil":
		f.Type = OTHER
		return
	default:
		if isDigit(str) {
			f.Type = INT
			return
		}

		f.fragmentByParams(params)
		if f.Type != UNKNOWN {
			return
		}
	}
	// TODO double check
	return f, fmt.Errorf("unknow parameter: %s", s)
}

func (f *fragment) fragmentByParams(params []parser.Param) {
	for _, param := range params {
		if param.Name == f.value {
			switch param.Type {
			case "bool":
				f.Type = BOOL
				return
			case "int":
				f.Type = INT
				return
			case "string":
				f.Type = STRING
				return
			case "Time":
				f.Type = TIME
			default:
				f.Type = OTHER
			}
		}
	}
}

func splitTemplate(tmpl string, params []parser.Param) (newList []fragment, err error) {
	tmpl += " "
	var out bytes.Buffer
	var t fragment
	for i := 0; i < len(tmpl); i++ {
		switch tmpl[i] {
		case '"':
			for {
				out.WriteByte(tmpl[i])
				if !stringHasMore(i, tmpl) {
					err = fmt.Errorf("incomplete code:%s", tmpl)
					return
				}
				i++
				if tmpl[i] == '"' && tmpl[i-1] != '\\' {
					out.WriteByte(tmpl[i])
					newList = append(newList, fragment{Type: STRING, value: out.String()})
					out.Reset()
					break
				}
			}
			continue
		case ' ':
			t, err = checkFragment(out.String(), params)
			if err != nil {
				return
			}
			if t.value != "" {
				newList = append(newList, t)
			}
			out.Reset()
		case '>', '<', '=', '!':
			t, err = checkFragment(out.String(), params)
			if err != nil {
				return
			}
			if t.value != "" {
				newList = append(newList, t)
			}
			out.Reset()

			out.WriteByte(tmpl[i])
			if tmpl[i+1] == '=' {
				out.WriteByte(tmpl[i+1])
				i++
			}
			t, err = checkFragment(out.String(), params)
			if err != nil {
				return
			}
			if t.value != "" {
				newList = append(newList, t)
			}
			out.Reset()
			continue
		case '&', '|':
			if tmpl[i+1] == tmpl[i] {
				t, err = checkFragment(out.String(), params)
				if err != nil {
					return
				}
				if t.value != "" {
					newList = append(newList, t)
				}
				out.Reset()
				out.WriteString(tmpl[i : i+2])
				t, err = checkFragment(out.String(), params)
				if err != nil {
					return
				}
				if t.value != "" {
					newList = append(newList, t)
				}
				out.Reset()
				i++
				continue
			}

		}

		out.WriteByte(tmpl[i])
	}
	t, err = checkFragment(out.String(), params)
	if err != nil {
		return
	}
	if t.value != "" {
		newList = append(newList, t)
	}
	// TODO check if verbose?
	if len(newList) == 0 {
		return
	}
	return
}

// check validition of clause's value
func checkTempleFragmentValid(list []fragment) error {
	for i := 1; i < len(list); i++ {
		switch list[i].Type {
		case IF, ELSE, END, BOOL, LOGICAL, WHERE, SET:
			continue
		case INT, STRING, OTHER, TIME:
			if i+2 < len(list) {
				if list[i+1].Type == EXPRESSION && list[i+2].Type == list[i].Type {
					i += 2
				} else {
					return fmt.Errorf("condition type not match：%s", fragmentToString(list[i:i+2]))
				}
			}
		default:
			return fmt.Errorf("unknow fragment ： %s ", list[i].value)
		}
	}
	return nil
}

func fragmentToString(list []fragment) string {
	var values []string

	if len(list) == 0 {
		return ""
	}
	for _, t := range list {
		values = append(values, t.value)
	}
	return strings.Join(values, " ")
}

func fragmentToSLice(list []fragment) (part slice, err error) {
	var values []string

	if len(list) == 0 {
		return
	}
	for _, t := range list {
		values = append(values, t.value)
	}
	part.Origin = strings.Join(values, " ")
	switch strings.ToLower(values[0]) {
	case "if":
		if len(values) > 1 {
			part.Type = IF
			part.Value = strings.Join(values[1:], " ")
			return
		}
	case "else":
		if len(values) == 1 {
			part.Type = ELSE
			return
		} else {
			if strings.ToLower(values[1]) == "if" && len(values) > 2 {
				part.Value = strings.Join(values[2:], " ")
				part.Type = ELSEIF
				return
			}
		}
	case "where":
		part.Type = WHERE
		return
	case "set":
		part.Type = SET
		return
	case "end":
		part.Type = END
		return
	}

	err = fmt.Errorf("syntax error:%s", strings.Join(values, " "))
	return
}
