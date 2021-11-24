package check

import (
	"fmt"
	"strconv"

	//"strconv"
	"strings"

	"gorm.io/gen/internal/model"
	"gorm.io/gen/internal/parser"
)

// Clause a symbol of clause, it can be sql condition clause, if clause, where clause, set clause and else cluase
type Clause interface {
	String() string
	Empty() bool
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
	return strings.ReplaceAll(strings.Join(s.Value, "+"), `"+"`, " ")
}

func (s SQLClause) Empty() bool {
	return s.VarName == ""
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

func (i IfClause) Empty() bool {
	return i.VarName == ""
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

func (e ElseClause) Empty() bool {
	return e.VarName == ""
}

// WhereClause where clause
type WhereClause struct {
	clause
	Value []Clause
}

func (w WhereClause) String() string {
	return fmt.Sprintf("helper.WhereClause(%s)", w.VarName)
}

func (w WhereClause) Empty() bool {
	return w.VarName == ""
}

// SetClause set clause
type SetClause struct {
	clause
	Value []Clause
}

func (w SetClause) String() string {
	return fmt.Sprintf("helper.SetClause(%s)", w.VarName)
}

func (w SetClause) Empty() bool {
	return w.VarName == ""
}

type RangeClause struct {
	clause
	Value []Clause
}

func (r RangeClause) String() string {
	//return fmt.Sprintf("helper.RangeClause(%s)", r.VarName)
	return fmt.Sprintf("helper.RangeClause(%s)", r.VarName)
}

func (r RangeClause) Empty() bool {
	return r.VarName == ""
}

// Slices split sql into chunks
type Slices struct {
	slices       []slice
	tmpl         []string
	currentIndex int
	Names        map[model.Status]int
}

// NewSlices create and initialize Slices
func NewSlices() *Slices {
	return &Slices{
		Names: map[model.Status]int{
			model.IF:    0,
			model.WHERE: 0,
			model.SET:   0,
			model.RANGE: 0,
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
	case model.RANGE:
		defer func() { s.Names[model.RANGE]++ }()
		return fmt.Sprintf("rangeCond%d", s.Names[model.RANGE])
	default:
		return fmt.Sprintf("Cond%d", s.currentIndex)
	}
}

func (s *Slices) appendIfCond(name, cond, result string) {
	cond = strings.ReplaceAll(cond, "$", "")
	s.tmpl = append(s.tmpl, fmt.Sprintf("%s = append(%s, helper.Cond{Cond: %s, Result: %s})", name, name, cond, result))
}

func (s *Slices) appendSetValue(name, result string) {
	s.tmpl = append(s.tmpl, fmt.Sprintf("%s = append(%s,  %s)", name, name, strings.TrimSpace(result)))
}

func (s *Slices) appendRangeValue(name, clause string) {
	s.tmpl = append(s.tmpl, fmt.Sprintf("\t%s = append(%s, %s)", name, name, clause))
}

// CreateIf create if clause code
func (s *Slices) CreateIf(name string) {
	s.tmpl = append(s.tmpl, fmt.Sprintf("%s := make([]helper.Cond, 0, 100)", name))
}

// CreateRange create range clause code
func (s *Slices) CreateRange(name, param string) {
	//s.tmpl = append(s.tmpl, fmt.Sprintf("var(\n\t%s\n)", param))
	s.tmpl = append(s.tmpl, fmt.Sprintf("%s := make([]string, 0, 100)", name))
	s.tmpl = append(s.tmpl, param)
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
		switch {
		case slice.Type.In(model.SQL, model.DATA, model.VARIABLE):
			sqlClause := s.parseSQL(name)
			res = append(res, sqlClause)
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=%s", name, sqlClause.String()))
		case slice.Type.In(model.IF):
			ifClause, err := s.parseIF()
			if err != nil {
				return nil, err
			}
			res = append(res, ifClause)
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=helper.IfClause(%s)", name, ifClause.VarName))
		case slice.Type.In(model.WHERE):
			whereClause, err := s.parseWhere()
			if err != nil {
				return nil, err
			}
			res = append(res, whereClause)
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=helper.WhereClause(%s)", name, whereClause.VarName))
		case slice.Type.In(model.SET):
			setClause, err := s.parseSet()
			if err != nil {
				return nil, err
			}
			res = append(res, setClause)
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=helper.SetClause(%s)", name, setClause.VarName))
		case slice.Type.In(model.RANGE):
			rangeClause, err := s.parseRange()
			if err != nil {
				return nil, err
			}
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=helper.RangeClause(%s)", name, rangeClause.VarName))
		case slice.Type.In(model.END):
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
		switch {
		case n.Type.In(model.SQL, model.DATA, model.VARIABLE, model.RANGEBODY):
			str := s.parseSQL(name)
			res.Value = append(res.Value, str)
			s.appendIfCond(name, res.Cond, str.String())
		case n.Type.In(model.IF):
			var ifClause IfClause
			ifClause, err = s.parseIF()
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendIfCond(name, res.Cond, ifClause.String())
		case n.Type.In(model.WHERE):
			var whereClause WhereClause
			whereClause, err = s.parseWhere()
			if err != nil {
				return
			}
			res.Value = append(res.Value, whereClause)
			s.appendIfCond(name, res.Cond, whereClause.String())
		case n.Type.In(model.SET):
			var setClause SetClause
			setClause, err = s.parseSet()
			if err != nil {
				return
			}
			res.Value = append(res.Value, setClause)
			s.appendIfCond(name, res.Cond, setClause.String())
		case n.Type.In(model.ELSEIF):
			elseClause := s.parseElSE(name)
			elseCond := elseClause.Cond
			elseClause.Cond = fmt.Sprintf("!(%s) && %s", strings.Join(cond, " || "), elseCond)
			res.Else = append(res.Else, elseClause)
			s.appendIfCond(name, elseClause.Cond, elseClause.String())
			cond = append(cond, elseCond)
		case n.Type.In(model.ELSE):
			elseClause := s.parseElSE(name)
			elseClause.Cond = fmt.Sprintf("!(%s)", strings.Join(cond, " || "))
			res.Else = append(res.Else, elseClause)
			s.appendIfCond(name, elseClause.Cond, elseClause.String())
		case n.Type.In(model.RANGE):
			var rangeClause RangeClause
			rangeClause, err = s.parseRange()
			if err != nil {
				return
			}
			res.Value = append(res.Value, rangeClause)
			s.appendIfCond(name, res.Cond, rangeClause.String())
		case n.Type.In(model.END):
			return
		default:
			err = fmt.Errorf("unknow clause : %s", n.Origin)
			return
		}
	}
	if s.Current().Type.In(model.END) {
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
		switch {
		case n.Type.In(model.SQL, model.DATA, model.VARIABLE, model.RANGEBODY):
			res.Value = append(res.Value, s.parseSQL(name))
		case n.Type.In(model.IF):
			ifClause, err := s.parseIF()
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
		case n.Type.In(model.WHERE):
			whereClause, err := s.parseWhere()
			if err != nil {
				return
			}
			res.Value = append(res.Value, whereClause)
		case n.Type.In(model.SET):
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
		switch {
		case n.Type.In(model.SQL, model.DATA, model.VARIABLE, model.RANGEBODY):
			strClause := s.parseSQL(name)
			res.Value = append(res.Value, strClause)
			s.appendSetValue(name, strClause.String())
		case n.Type.In(model.IF):
			var ifClause IfClause
			ifClause, err = s.parseIF()
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendSetValue(name, ifClause.String())
		case n.Type.In(model.RANGE):
			var rangeClause RangeClause
			rangeClause, err = s.parseRange()
			if err != nil {
				return
			}
			res.Value = append(res.Value, rangeClause)
			s.appendSetValue(name, rangeClause.String())
		case n.Type.In(model.END):
			return
		default:
			err = fmt.Errorf("unknow clause : %s", n.Origin)
			return
		}
	}
	if s.Current().Type == model.END {
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
		switch {
		case n.Type.In(model.SQL, model.DATA, model.VARIABLE):
			strClause := s.parseSQL(name)
			res.Value = append(res.Value, strClause)
			s.appendSetValue(name, strClause.String())
		case n.Type.In(model.IF):
			var ifClause IfClause
			ifClause, err = s.parseIF()
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendSetValue(name, ifClause.String())
		case n.Type.In(model.END):
			return
		default:
			err = fmt.Errorf("unknow clause : %s", n.Origin)
			return
		}
	}
	if s.Current().Type == model.END {
		return
	}
	err = fmt.Errorf("incomplete SQL,set not end")
	return
}

// parseWhere parse where clause, the clause' type must be one of if, SQL condition
func (s *Slices) parseRange() (res RangeClause, err error) {
	slice := s.Current()
	name := s.GetName(slice.Type)
	s.CreateRange(name, slice.Value)

	res.VarName = name
	res.Type = slice.Type
	for s.HasMore() {
		n := s.Next()
		switch {
		case n.Type.In(model.SQL, model.DATA, model.VARIABLE, model.RANGEBODY):
			strClause := s.parseSQL(name)
			res.Value = append(res.Value, strClause)
			s.appendRangeValue(name, strClause.String())
		case n.Type.In(model.IF):
			var ifClause IfClause
			ifClause, err = s.parseIF()
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendRangeValue(name, ifClause.String())
		case n.Type.In(model.SET):
			var setClause SetClause
			setClause, err = s.parseSet()
			if err != nil {
				return
			}
			res.Value = append(res.Value, setClause)
			s.appendRangeValue(name, setClause.String())
		case n.Type.In(model.WHERE):
			var whereClause WhereClause
			whereClause, err = s.parseWhere()
			if err != nil {
				return
			}
			res.Value = append(res.Value, whereClause)
			s.appendRangeValue(name, whereClause.String())
		case n.Type.In(model.END):
			s.tmpl = append(s.tmpl, "}")
			return
		default:
			err = fmt.Errorf("unknow clause : %s", n.Origin)
			return
		}
	}
	if s.Current().Type == model.END {
		return
	}
	err = fmt.Errorf("incomplete SQL,where not end")
	return
}

// parseSQL parse sql condition, the clause' type must be one of SQL condition, VARIABLE, Data
func (s *Slices) parseSQL(name string) (res SQLClause) {
	res.VarName = name
	res.Type = model.SQL
	for slice := s.Current(); ; slice = s.Next() {
		switch {
		case slice.Type.In(model.SQL, model.VARIABLE, model.DATA, model.RANGEBODY):
			if slice.Type.In(model.RANGEBODY) {
				result := strings.Replace(slice.Value, "$", "%v", 1)
				result = result[1 : len(result)-1]
				if i := strings.Index(result, "%"); i != -1 {
					result = fmt.Sprintf(`fmt.Sprintf("%s", %s)`, result[:i+2], result[i+2:])
				}
				res.Value = append(res.Value, result)
			} else {
				res.Value = append(res.Value, slice.Value)
			}
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
	Type    model.Status
	Value   string
	IsArray bool
}

func checkFragment(s string, params []parser.Param) (f fragment, err error) {
	f = fragment{Value: strings.Trim(s, " ")}
	str := strings.ToLower(strings.Trim(s, " "))
	switch str {
	case isDigit(str):
		f.Type |= model.DIGIT
	case isRangeVar(str):
		f.Type |= model.RANGEVAR | model.DIGIT
	case "&&", "||":
		f.Type |= model.LOGICAL
	case ">", "<", ">=", "<=", "==", "!=":
		f.Type |= model.EXPRESSION
	case "=":
		f.Type |= model.ASSIGN
	case "end":
		f.Type |= model.END
	case "if":
		f.Type |= model.IF
	case "set":
		f.Type |= model.SET
	case "else":
		f.Type |= model.ELSE
	case "where":
		f.Type |= model.WHERE
	case "range":
		f.Type |= model.RANGE
	case "true", "false":
		f.Type |= model.BOOL
	case "nil":
		f.Type |= model.NIL
	default:
		f.fragmentByParams(params)
		if f.Type.In(model.UNKNOWN) {
			err = fmt.Errorf("unknow parameter: %s", s)
		}
	}
	return
}

func (f *fragment) fragmentByParams(params []parser.Param) {
	f.Value = strings.Replace(f.Value, "@", "", 1)
	for _, param := range params {
		if param.Name == f.Value {
			f.IsArray = param.IsArray
			switch param.Type {
			case "bool":
				f.Type |= model.BOOL
				return
			case "int", "int8", "int16", "int32", "int64", "float32", "float64":
				f.Type |= model.DIGIT
				return
			case "string":
				if f.IsArray {
					f.Type |= model.ARRAY
				} else {
					f.Type |= model.STRING
				}
				return
			case "Time":
				f.Type |= model.TIME
			default:
				f.Type |= model.OTHER
			}
		}
	}
}

func splitTemplate(tmpl string, params []parser.Param) (fragList []fragment, err error) {
	var buf model.SQLBuffer
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
					fragList = append(fragList, fragment{Type: model.STRING, Value: buf.Dump()})
					break
				}
			}
		case ' ', ',':
			if sqlClause := buf.Dump(); sqlClause != "" {
				f, err = checkFragment(sqlClause, params)
				if err != nil {
					return nil, err
				}
				if f.Type.In(model.RANGEVAR, model.DIGIT) && len(fragList) > 0 && fragList[len(fragList)-1].Type.In(model.RANGEVAR, model.DIGIT) {
					f.Type &^= model.DIGIT
				}
				fragList = append(fragList, f)
			}
		case '>', '<', '=', '!':
			if i < 1 {
				return nil, fmt.Errorf("wrong code expression:%s", tmpl)
			}
			if tmpl[i-1] == ':' {
				buf.Truncate(buf.Len() - 1)
			}
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
					Type:  model.LOGICAL,
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
		switch {
		case list[i].Type.In(model.IF, model.ELSE, model.END, model.BOOL, model.LOGICAL, model.WHERE, model.SET, model.RANGE, model.ASSIGN):
			continue
		case list[i].Type.In(model.DIGIT, model.STRING, model.OTHER, model.TIME, model.NIL, model.RANGEVAR):
			if list[0].Type.In(model.RANGE) {
				if !list[len(list)-1].Type.In(model.ARRAY) {
					return fmt.Errorf("cannot range over %s", fragmentToString(list[len(list)-1:]))
				}
				return nil
			}
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
	if expr[1].Type&model.EXPRESSION == 0 {
		return false
	}
	//Only arrays can be compared with nil
	if expr[0].Type == model.NIL || expr[2].Type == model.NIL {
		return expr[0].IsArray || expr[2].IsArray
	}

	return expr[0].Type.In(model.RANGEVAR) || expr[0].Type&expr[2].Type != 0
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
			part.Type = model.IF
			part.Value = strings.Join(values[1:], " ")
			return
		}
	case "else":
		if len(values) == 1 {
			part.Type = model.ELSE
			return
		} else {
			if strings.ToLower(values[1]) == "if" && len(values) > 2 {
				part.Value = strings.Join(values[2:], " ")
				part.Type = model.ELSEIF
				return
			}
		}
	case "where":
		part.Type = model.WHERE
		return
	case "set":
		part.Type = model.SET
		return
	case "range":
		part.Type = model.RANGE
		var (
			idx string
			res string
		)

		for _, e := range list {
			if e.Type.In(model.RANGEVAR) {
				if e.Type.In(model.DIGIT) {
					idx = fmt.Sprintf("%s, ", strings.ReplaceAll(e.Value, "$", ""))
				} else {
					res = fmt.Sprintf("%s", strings.ReplaceAll(e.Value, "$", ""))
				}
			}
		}
		part.Value = fmt.Sprintf("for %s %s := range %s {", idx, res, list[len(list)-1].Value)
		return
	case "end":
		part.Type = model.END
		return
	}

	err = rangeFragmentToSlice(list, &part)
	return
}

func rangeFragmentToSlice(list []fragment, part *slice) error {
	rangeExp := false
	for _, t := range list {
		if t.Type.In(model.RANGEVAR) {
			rangeExp = true
		}
		if t.Type.In(model.EXPRESSION) {
			part.Type |= model.IF
		}
	}
	if !rangeExp {
		return fmt.Errorf("syntax error:%s", part.Origin)
	}
	part.Type |= model.RANGEBODY
	part.Value = strconv.Quote(part.Origin)
	return nil
}
