package check

import (
	"fmt"
	"strings"

<<<<<<< HEAD
	"gorm.io/gen/internal/model"
=======
	"gorm.io/gen/internal/models"
>>>>>>> feat: index tag
	"gorm.io/gen/internal/parser"
)

// Clause a symbol of clause, it can be sql condition clause, if clause, where clause, set clause and else cluase
type Clause interface {
	String() string
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
<<<<<<< HEAD
	Type    model.Status
=======
	Type    models.Status
>>>>>>> feat: index tag
}

// SQLClause sql condition clause
type SQLClause struct {
	clause
	Value []string
}

func (s SQLClause) String() string {
	return strings.ReplaceAll(strings.Join(s.Value, "+"), `"+"`, "")
}

// IfClause if clause
type IfClause struct {
	clause
	Cond  string
	Value []Clause
	Else  []Clause
}

func (i IfClause) String() string {
	return fmt.Sprintf("helper.IfClause(%s)", i.VarName)
}

// ElseClause else clause
type ElseClause struct {
	clause
	Cond  string
	Value []Clause
}

func (e ElseClause) String() (res string) {
	condList := make([]string, len(e.Value))
	for i, v := range e.Value {
		condList[i] = v.String()
	}
	return strings.ReplaceAll(strings.Join(condList, "+"), `"+"`, "")
}

// WhereClause where clause
type WhereClause struct {
	clause
	Value []Clause
}

func (w WhereClause) String() string {
	return fmt.Sprintf("helper.WhereClause(%s)", w.VarName)
}

// SetClause set clause
type SetClause struct {
	clause
	Value []Clause
}

func (w SetClause) String() string {
	return fmt.Sprintf("helper.SetClause(%s)", w.VarName)
}

// Slices split sql into chunks
type Slices struct {
	slices       []slice
	tmpl         []string
	currentIndex int
<<<<<<< HEAD
	Names        map[model.Status]int
=======
	Names        map[models.Status]int
>>>>>>> feat: index tag
}

// NewSlices create and initialize Slices
func NewSlices() *Slices {
	return &Slices{
<<<<<<< HEAD
		Names: map[model.Status]int{
			model.IF:    0,
			model.WHERE: 0,
			model.SET:   0,
=======
		Names: map[models.Status]int{
			models.IF:    0,
			models.WHERE: 0,
			models.SET:   0,
>>>>>>> feat: index tag
		},
	}
}

// Len return length of s.slices
func (s *Slices) Len() int {
	return len(s.slices)
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
<<<<<<< HEAD
func (s *Slices) GetName(status model.Status) string {
	switch status {
	case model.IF:
		defer func() { s.Names[model.IF]++ }()
		return fmt.Sprintf("ifCond%d", s.Names[model.IF])
	case model.WHERE:
		defer func() { s.Names[model.WHERE]++ }()
		return fmt.Sprintf("whereCond%d", s.Names[model.WHERE])
	case model.SET:
		defer func() { s.Names[model.SET]++ }()
		return fmt.Sprintf("setCond%d", s.Names[model.SET])
=======
func (s *Slices) GetName(status models.Status) string {
	switch status {
	case models.IF:
		defer func() { s.Names[models.IF]++ }()
		return fmt.Sprintf("ifCond%d", s.Names[models.IF])
	case models.WHERE:
		defer func() { s.Names[models.WHERE]++ }()
		return fmt.Sprintf("whereCond%d", s.Names[models.WHERE])
	case models.SET:
		defer func() { s.Names[models.SET]++ }()
		return fmt.Sprintf("setCond%d", s.Names[models.SET])
>>>>>>> feat: index tag
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
func (s *Slices) parse() ([]Clause, error) {
	if s.IsNull() {
		return nil, fmt.Errorf("sql is null")
	}

	name := "generateSQL"
	res := make([]Clause, 0, s.Len())
	for slice := s.Current(); ; slice = s.Next() {
		s.tmpl = append(s.tmpl, "")
		switch slice.Type {
<<<<<<< HEAD
		case model.SQL, model.DATA, model.VARIABLE:
			sqlClause := s.parseSQL(name)
			res = append(res, sqlClause)
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=%s", name, sqlClause.String()))
		case model.IF:
=======
		case models.SQL, models.DATA, models.VARIABLE:
			sqlClause := s.parseSQL(name)
			res = append(res, sqlClause)
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=%s", name, sqlClause.String()))
		case models.IF:
>>>>>>> feat: index tag
			ifClause, err := s.parseIF()
			if err != nil {
				return nil, err
			}
			res = append(res, ifClause)
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=helper.IfClause(%s)", name, ifClause.VarName))
<<<<<<< HEAD
		case model.WHERE:
=======
		case models.WHERE:
>>>>>>> feat: index tag
			whereClause, err := s.parseWhere()
			if err != nil {
				return nil, err
			}
			res = append(res, whereClause)
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=helper.WhereClause(%s)", name, whereClause.VarName))
<<<<<<< HEAD
		case model.SET:
=======
		case models.SET:
>>>>>>> feat: index tag
			setClause, err := s.parseSet()
			if err != nil {
				return nil, err
			}
			res = append(res, setClause)
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=helper.SetClause(%s)", name, setClause.VarName))
<<<<<<< HEAD
		case model.END:
=======
		case models.END:
>>>>>>> feat: index tag
		default:
			return nil, fmt.Errorf("unknow clause:%s", slice.Origin)
		}

		if !s.HasMore() {
			break
		}
	}
	return res, nil
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
<<<<<<< HEAD
		case model.SQL, model.DATA, model.VARIABLE:
			str := s.parseSQL(name)
			res.Value = append(res.Value, str)
			s.appendIfCond(name, res.Cond, str.String())
		case model.IF:
=======
		case models.SQL, models.DATA, models.VARIABLE:
			str := s.parseSQL(name)
			res.Value = append(res.Value, str)
			s.appendIfCond(name, res.Cond, str.String())
		case models.IF:
>>>>>>> feat: index tag
			var ifClause IfClause
			ifClause, err = s.parseIF()
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendIfCond(name, res.Cond, ifClause.String())
<<<<<<< HEAD
		case model.WHERE:
=======
		case models.WHERE:
>>>>>>> feat: index tag
			var whereClause WhereClause
			whereClause, err = s.parseWhere()
			if err != nil {
				return
			}
			res.Value = append(res.Value, whereClause)
			s.appendIfCond(name, res.Cond, whereClause.String())
<<<<<<< HEAD
		case model.SET:
=======
		case models.SET:
>>>>>>> feat: index tag
			var setClause SetClause
			setClause, err = s.parseSet()
			if err != nil {
				return
			}
			res.Value = append(res.Value, setClause)
			s.appendIfCond(name, res.Cond, setClause.String())
<<<<<<< HEAD
		case model.ELSEIF:
=======
		case models.ELSEIF:
>>>>>>> feat: index tag
			elseClause := s.parseElSE(name)
			elseCond := elseClause.Cond
			elseClause.Cond = fmt.Sprintf("!(%s) && %s", strings.Join(cond, " || "), elseCond)
			res.Else = append(res.Else, elseClause)
			s.appendIfCond(name, elseClause.Cond, elseClause.String())
			cond = append(cond, elseCond)
<<<<<<< HEAD
		case model.ELSE:
=======
		case models.ELSE:
>>>>>>> feat: index tag
			elseClause := s.parseElSE(name)
			elseClause.Cond = fmt.Sprintf("!(%s)", strings.Join(cond, " || "))
			res.Else = append(res.Else, elseClause)
			s.appendIfCond(name, elseClause.Cond, elseClause.String())
<<<<<<< HEAD
		case model.END:
=======
		case models.END:
>>>>>>> feat: index tag
			return
		default:
			err = fmt.Errorf("unknow clause : %s", n.Origin)
			return
		}
	}
<<<<<<< HEAD
	if s.Current().Type == model.END {
=======
	if s.Current().Type == models.END {
>>>>>>> feat: index tag
		return
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

	if !s.HasMore() {
		return
	}
	for n := s.Next(); s.HasMore(); n = s.Next() {
		switch n.Type {
<<<<<<< HEAD
		case model.SQL, model.DATA, model.VARIABLE:
			res.Value = append(res.Value, s.parseSQL(name))
		case model.IF:
=======
		case models.SQL, models.DATA, models.VARIABLE:
			res.Value = append(res.Value, s.parseSQL(name))
		case models.IF:
>>>>>>> feat: index tag
			ifClause, err := s.parseIF()
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
<<<<<<< HEAD
		case model.WHERE:
=======
		case models.WHERE:
>>>>>>> feat: index tag
			whereClause, err := s.parseWhere()
			if err != nil {
				return
			}
			res.Value = append(res.Value, whereClause)
<<<<<<< HEAD
		case model.SET:
=======
		case models.SET:
>>>>>>> feat: index tag
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
<<<<<<< HEAD
		case model.SQL, model.DATA, model.VARIABLE:
			strClause := s.parseSQL(name)
			res.Value = append(res.Value, strClause)
			s.appendSetValue(name, strClause.String())
		case model.IF:
=======
		case models.SQL, models.DATA, models.VARIABLE:
			strClause := s.parseSQL(name)
			res.Value = append(res.Value, strClause)
			s.appendSetValue(name, strClause.String())
		case models.IF:
>>>>>>> feat: index tag
			var ifClause IfClause
			ifClause, err = s.parseIF()
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendSetValue(name, ifClause.String())
<<<<<<< HEAD
		case model.END:
=======
		case models.END:
>>>>>>> feat: index tag
			return
		default:
			err = fmt.Errorf("unknow clause : %s", n.Origin)
			return
		}
	}
<<<<<<< HEAD
	if s.Current().Type == model.END {
=======
	if s.Current().Type == models.END {
>>>>>>> feat: index tag
		return
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
<<<<<<< HEAD
		case model.SQL, model.DATA, model.VARIABLE:
			strClause := s.parseSQL(name)
			res.Value = append(res.Value, strClause)
			s.appendSetValue(name, strClause.String())
		case model.IF:
=======
		case models.SQL, models.DATA, models.VARIABLE:
			strClause := s.parseSQL(name)
			res.Value = append(res.Value, strClause)
			s.appendSetValue(name, strClause.String())
		case models.IF:
>>>>>>> feat: index tag
			var ifClause IfClause
			ifClause, err = s.parseIF()
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendSetValue(name, ifClause.String())
<<<<<<< HEAD
		case model.END:
=======
		case models.END:
>>>>>>> feat: index tag
			return
		default:
			err = fmt.Errorf("unknow clause : %s", n.Origin)
			return
		}
	}
<<<<<<< HEAD
	if s.Current().Type == model.END {
=======
	if s.Current().Type == models.END {
>>>>>>> feat: index tag
		return
	}
	err = fmt.Errorf("incomplete SQL,set not end")
	return
}

// parseSQL parse sql condition, the clause' type must be one of SQL condition, VARIABLE, Data
func (s *Slices) parseSQL(name string) (res SQLClause) {
	res.VarName = name
<<<<<<< HEAD
	res.Type = model.SQL
	for slice := s.Current(); ; slice = s.Next() {
		switch slice.Type {
		case model.SQL, model.VARIABLE, model.DATA:
=======
	res.Type = models.SQL
	for slice := s.Current(); ; slice = s.Next() {
		switch slice.Type {
		case models.SQL, models.VARIABLE, models.DATA:
>>>>>>> feat: index tag
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
<<<<<<< HEAD
	Type    model.Status
=======
	Type    models.Status
>>>>>>> feat: index tag
	Value   string
	IsArray bool
}

func checkFragment(s string, params []parser.Param) (f fragment, err error) {
	digital := func(str string) string {
		if isDigit(str) {
			return "<integer>"
		}
		return str
	}

<<<<<<< HEAD
	f = fragment{Type: model.UNKNOWN, Value: strings.Trim(s, " ")}
	str := strings.ToLower(strings.Trim(s, " "))
	switch digital(str) {
	case "<integer>":
		f.Type = model.INT
	case "&&", "||":
		f.Type = model.LOGICAL
	case ">", "<", ">=", "<=", "==", "!=":
		f.Type = model.EXPRESSION
	case "end":
		f.Type = model.END
	case "if":
		f.Type = model.IF
	case "set":
		f.Type = model.SET
	case "else":
		f.Type = model.ELSE
	case "where":
		f.Type = model.WHERE
	case "true", "false":
		f.Type = model.BOOL
	case "nil":
		f.Type = model.NIL
	default:
		f.fragmentByParams(params)
		if f.Type == model.UNKNOWN {
=======
	f = fragment{Type: models.UNKNOWN, Value: strings.Trim(s, " ")}
	str := strings.ToLower(strings.Trim(s, " "))
	switch digital(str) {
	case "<integer>":
		f.Type = models.INT
	case "&&", "||":
		f.Type = models.LOGICAL
	case ">", "<", ">=", "<=", "==", "!=":
		f.Type = models.EXPRESSION
	case "end":
		f.Type = models.END
	case "if":
		f.Type = models.IF
	case "set":
		f.Type = models.SET
	case "else":
		f.Type = models.ELSE
	case "where":
		f.Type = models.WHERE
	case "true", "false":
		f.Type = models.BOOL
	case "nil":
		f.Type = models.NIL
	default:
		f.fragmentByParams(params)
		if f.Type == models.UNKNOWN {
>>>>>>> feat: index tag
			err = fmt.Errorf("unknow parameter: %s", s)
		}
	}
	return
}

func (f *fragment) fragmentByParams(params []parser.Param) {
	for _, param := range params {
		if param.Name == f.Value {
			f.IsArray = param.IsArray
			switch param.Type {
			case "bool":
<<<<<<< HEAD
				f.Type = model.BOOL
				return
			case "int":
				f.Type = model.INT
				return
			case "string":
				f.Type = model.STRING
				return
			case "Time":
				f.Type = model.TIME
			default:
				f.Type = model.OTHER
=======
				f.Type = models.BOOL
				return
			case "int":
				f.Type = models.INT
				return
			case "string":
				f.Type = models.STRING
				return
			case "Time":
				f.Type = models.TIME
			default:
				f.Type = models.OTHER
>>>>>>> feat: index tag
			}
		}
	}
}

func splitTemplate(tmpl string, params []parser.Param) (fragList []fragment, err error) {
<<<<<<< HEAD
	var buf model.SQLBuffer
=======
	var buf models.Sql
>>>>>>> feat: index tag
	var f fragment
	for i := 0; !strOutrange(i, tmpl); i++ {
		switch tmpl[i] {
		case '"':
			_ = buf.WriteByte(tmpl[i])
			for i++; ; i++ {
				if strOutrange(i, tmpl) {
					return nil, fmt.Errorf("incomplete code:%s", tmpl)
				}
				_ = buf.WriteByte(tmpl[i])

				if tmpl[i] == '"' && tmpl[i-1] != '\\' {
<<<<<<< HEAD
					fragList = append(fragList, fragment{Type: model.STRING, Value: buf.Dump()})
=======
					fragList = append(fragList, fragment{Type: models.STRING, Value: buf.Dump()})
>>>>>>> feat: index tag
					break
				}
			}
		case ' ':
			if sqlClause := buf.Dump(); sqlClause != "" {
				f, err = checkFragment(sqlClause, params)
				if err != nil {
					return nil, err
				}
				fragList = append(fragList, f)
			}
		case '>', '<', '=', '!':
			if sqlClause := buf.Dump(); sqlClause != "" {
				f, err = checkFragment(sqlClause, params)
				if err != nil {
					return nil, err
				}
				fragList = append(fragList, f)
			}

			_ = buf.WriteByte(tmpl[i])

			if strOutrange(i+1, tmpl) {
				return nil, fmt.Errorf("incomplete code:%s", tmpl)
			}
			if tmpl[i+1] == '=' {
				_ = buf.WriteByte(tmpl[i+1])
				i++
			}

			f, err = checkFragment(buf.Dump(), params)
			if err != nil {
				return nil, err
			}
			fragList = append(fragList, f)
		case '&', '|':
			if strOutrange(i+1, tmpl) {
				return nil, fmt.Errorf("incomplete code:%s", tmpl)
			}

			if tmpl[i+1] == tmpl[i] {
				i++

				if sqlClause := buf.Dump(); sqlClause != "" {
					f, err = checkFragment(sqlClause, params)
					if err != nil {
						return nil, err
					}
					fragList = append(fragList, f)
				}

				// write && or ||
				fragList = append(fragList, fragment{
<<<<<<< HEAD
					Type:  model.LOGICAL,
=======
					Type:  models.LOGICAL,
>>>>>>> feat: index tag
					Value: tmpl[i-1 : i+1],
				})
			}
		default:
			_ = buf.WriteByte(tmpl[i])
		}
	}

	if sqlClause := buf.Dump(); sqlClause != "" {
		f, err = checkFragment(sqlClause, params)
		if err != nil {
			return nil, err
		}
		fragList = append(fragList, f)
	}
	return fragList, nil
}

// check validition of clause's value
func checkTempleFragmentValid(list []fragment) error {
	for i := 1; i < len(list); i++ {
		switch list[i].Type {
<<<<<<< HEAD
		case model.IF, model.ELSE, model.END, model.BOOL, model.LOGICAL, model.WHERE, model.SET:
			continue
		case model.INT, model.STRING, model.OTHER, model.TIME, model.NIL:
=======
		case models.IF, models.ELSE, models.END, models.BOOL, models.LOGICAL, models.WHERE, models.SET:
			continue
		case models.INT, models.STRING, models.OTHER, models.TIME, models.NIL:
>>>>>>> feat: index tag
			if i+2 < len(list) {
				if isExpressionValid(list[i : i+3]) {
					i += 2
				} else {
					return fmt.Errorf("condition type not match：%s", fragmentToString(list[i:i+3]))
				}
			}
		default:
			return fmt.Errorf("unknow fragment ： %s ", list[i].Value)
		}
	}
	return nil
}

// isExpressionValid  check express valid
func isExpressionValid(expr []fragment) bool {
	if len(expr) != 3 {
		return false
	}
<<<<<<< HEAD
	if expr[1].Type != model.EXPRESSION {
		return false
	}
	//Only arrays can be compared with nil
	if expr[0].Type == model.NIL || expr[2].Type == model.NIL {
=======
	if expr[1].Type != models.EXPRESSION {
		return false
	}
	//Only arrays can be compared with nil
	if expr[0].Type == models.NIL || expr[2].Type == models.NIL {
>>>>>>> feat: index tag
		return expr[0].IsArray || expr[2].IsArray
	}

	return expr[0].Type == expr[2].Type
}

func fragmentToString(list []fragment) string {
	var values []string

	if len(list) == 0 {
		return ""
	}
	for _, t := range list {
		values = append(values, t.Value)
	}
	return strings.Join(values, " ")
}

func fragmentToSLice(list []fragment) (part slice, err error) {
	var values []string

	if len(list) == 0 {
		return
	}
	for _, t := range list {
		values = append(values, t.Value)
	}
	part.Origin = strings.Join(values, " ")
	switch strings.ToLower(values[0]) {
	case "if":
		if len(values) > 1 {
<<<<<<< HEAD
			part.Type = model.IF
=======
			part.Type = models.IF
>>>>>>> feat: index tag
			part.Value = strings.Join(values[1:], " ")
			return
		}
	case "else":
		if len(values) == 1 {
<<<<<<< HEAD
			part.Type = model.ELSE
=======
			part.Type = models.ELSE
>>>>>>> feat: index tag
			return
		} else {
			if strings.ToLower(values[1]) == "if" && len(values) > 2 {
				part.Value = strings.Join(values[2:], " ")
<<<<<<< HEAD
				part.Type = model.ELSEIF
=======
				part.Type = models.ELSEIF
>>>>>>> feat: index tag
				return
			}
		}
	case "where":
<<<<<<< HEAD
		part.Type = model.WHERE
		return
	case "set":
		part.Type = model.SET
		return
	case "end":
		part.Type = model.END
=======
		part.Type = models.WHERE
		return
	case "set":
		part.Type = models.SET
		return
	case "end":
		part.Type = models.END
>>>>>>> feat: index tag
		return
	}

	err = fmt.Errorf("syntax error:%s", strings.Join(values, " "))
	return
}
