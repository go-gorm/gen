package check

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gen/internal/model"
	"gorm.io/gen/internal/parser"
)

// Clause a symbol of clause, it can be sql condition clause, if clause, where clause, set clause and else clause
type Clause interface {
	String() string
	Creat() string
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
	return strings.ReplaceAll(strings.Join(s.Value, "+"), `"+"`, "")
}

func (s SQLClause) Creat() string {
	return fmt.Sprintf("%s.WriteString(%s)", s.VarName, s.String())
}

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

func (i IfClause) Creat() string {
	return fmt.Sprintf("%s {", i.String())
}

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

func (e ElseClause) Creat() string {
	return fmt.Sprintf("} %s {", e.String())
}

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
func (w WhereClause) Creat() string {
	return fmt.Sprintf("var %s strings.Builder", w.VarName)
}
func (w WhereClause) Finish(name string) string {
	return fmt.Sprintf("%s.WriteString(%s)", name, w.String())
}

// SetClause set clause
type SetClause struct {
	clause
	Value []Clause
}

func (s SetClause) String() string {
	return fmt.Sprintf("helper.SetTrim(%s.String())", s.VarName)
}

func (s SetClause) Creat() string {
	return fmt.Sprintf("var %s strings.Builder", s.VarName)
}

func (s SetClause) Finish() string {
	return fmt.Sprintf("%s.WriteString(%s)", s.VarName, s.String())
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
func (f ForClause) Creat() string {
	return f.String()
}

func (f ForClause) Finish() string {
	return "}"
}

// Sections split sql into chunks
type Sections struct {
	members      []section
	Tmpl         []string
	currentIndex int
	ClauseTotal  map[model.Status]int
	forValue     []ForRange
}

// NewSections create and initialize Sections
func NewSections() *Sections {
	return &Sections{
		ClauseTotal: map[model.Status]int{
			model.WHERE: 0,
			model.SET:   0,
		},
	}
}

// next: return next section and increase index by 1
func (s *Sections) next() section {
	if s.currentIndex < len(s.members)-1 {
		s.currentIndex++
		return s.members[s.currentIndex]
	}
	return section{Type: model.END}
}

// SubIndex take index one step back
func (s *Sections) SubIndex() {
	s.currentIndex--
}

// HasMore: is has more section
func (s *Sections) HasMore() bool {
	return s.currentIndex < len(s.members)-1
}

// IsNull whether section is empty
func (s *Sections) IsNull() bool {
	return len(s.members) == 0
}

// current return current section
func (s *Sections) current() section {
	return s.members[s.currentIndex]
}

func (s *Sections) tmplAppend(value string) {
	s.Tmpl = append(s.Tmpl, value)
}

func (s *Sections) isInForValue(value string) (ForRange, bool) {
	valueList := strings.Split(value, ".")
	for _, v := range s.forValue {
		if v.value == valueList[0] {
			if len(valueList) > 1 {
				v.suffix = "." + strings.Join(valueList[1:], ".")
			}
			return v, true
		}
	}
	return ForRange{}, false
}
func (s *Sections) hasSameName(value string) bool {
	for _, p := range s.members {
		if p.Type == model.FOR && p.ForRange.value == value {
			return true
		}
	}
	return false

}

// BuildSQLClause sql sections and append to tmpl, return a Clause array
func (s *Sections) BuildSQLClause() ([]Clause, error) {
	if s.IsNull() {
		return nil, fmt.Errorf("sql is null")
	}
	name := "generateSQL"
	res := make([]Clause, 0, len(s.members))
	for {
		c := s.current()
		switch c.Type {
		case model.SQL, model.DATA, model.VARIABLE:
			sqlClause := s.parseSQL(name)
			res = append(res, sqlClause)
			s.tmplAppend(sqlClause.Finish())
		case model.IF:
			ifClause, err := s.parseIF(name)
			if err != nil {
				return nil, err
			}
			res = append(res, ifClause)
			s.tmplAppend(ifClause.Finish())
		case model.WHERE:
			whereClause, err := s.parseWhere()
			if err != nil {
				return nil, err
			}
			res = append(res, whereClause)
			s.tmplAppend(whereClause.Finish(name))
		case model.SET:
			setClause, err := s.parseSet()
			if err != nil {
				return nil, err
			}
			res = append(res, setClause)
			s.tmplAppend(setClause.Finish())
		case model.FOR:
			forClause, err := s.parseFor(name)
			_, _ = forClause, err
			if err != nil {
				return nil, err
			}
			res = append(res, forClause)
			s.tmplAppend(forClause.Finish())
		case model.END:
		default:
			return nil, fmt.Errorf("unknow clause:%s", c.Origin)
		}
		if !s.HasMore() {
			break
		}
		c = s.next()
	}
	return res, nil
}

// parseIF parse if clause
func (s *Sections) parseIF(name string) (res IfClause, err error) {
	c := s.current()
	res.slice = c

	s.tmplAppend(res.Creat())
	if !s.HasMore() {
		return
	}
	c = s.next()
	for {
		switch c.Type {
		case model.SQL, model.DATA, model.VARIABLE:
			sqlClause := s.parseSQL(name)
			res.Value = append(res.Value, sqlClause)
			s.tmplAppend(sqlClause.Finish())
		case model.IF:
			var ifClause IfClause
			ifClause, err = s.parseIF(name)
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.tmplAppend(ifClause.Finish())
		case model.WHERE:
			var whereClause WhereClause
			whereClause, err = s.parseWhere()
			if err != nil {
				return
			}
			res.Value = append(res.Value, whereClause)
			s.tmplAppend(whereClause.Finish(name))
		case model.SET:
			var setClause SetClause
			setClause, err = s.parseSet()
			if err != nil {
				return
			}
			res.Value = append(res.Value, setClause)
			s.tmplAppend(setClause.Finish())
		case model.ELSE:
			var elseClause ElseClause
			elseClause, err = s.parseElSE(name)
			if err != nil {
				return
			}
			res.Value = append(res.Value, elseClause)
		case model.FOR:
			var forClause ForClause
			forClause, err = s.parseFor(name)
			if err != nil {
				return
			}
			res.Value = append(res.Value, forClause)
			s.tmplAppend(res.Finish())
		case model.END:
			return
		default:
			err = fmt.Errorf("unknow clause : %s", c.Origin)
			return
		}
		if !s.HasMore() {
			break
		}
		c = s.next()
	}
	if c.isEnd() {
		err = fmt.Errorf("incomplete SQL,if not end")
	}
	return
}

// parseElSE parse else clause, the clause' type must be one of if, where, set, SQL condition
func (s *Sections) parseElSE(name string) (res ElseClause, err error) {
	res.slice = s.current()
	s.tmplAppend(res.Creat())

	if !s.HasMore() {
		return
	}
	c := s.next()
	for {
		switch c.Type {
		case model.SQL, model.DATA, model.VARIABLE:
			sqlClause := s.parseSQL(name)
			res.Value = append(res.Value, sqlClause)
			s.tmplAppend(sqlClause.Creat())
		case model.IF:
			var ifClause IfClause
			ifClause, err = s.parseIF(name)
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.tmplAppend(ifClause.Finish())
		case model.WHERE:
			var whereClause WhereClause
			whereClause, err = s.parseWhere()
			if err != nil {
				return
			}
			res.Value = append(res.Value, whereClause)
			s.tmplAppend(whereClause.Finish(name))
		case model.SET:
			var setClause SetClause
			setClause, err = s.parseSet()
			if err != nil {
				return
			}
			res.Value = append(res.Value, setClause)
			s.tmplAppend(setClause.Finish())
		case model.ELSE:
			var elseClause ElseClause
			elseClause, err = s.parseElSE(name)
			if err != nil {
				return
			}
			res.Value = append(res.Value, elseClause)
		case model.FOR:
			var forClause ForClause
			forClause, err = s.parseFor(name)
			if err != nil {
				return
			}
			res.Value = append(res.Value, forClause)
			s.tmplAppend(forClause.Finish())
		default:
			s.SubIndex()
			return
		}
		if !s.HasMore() {
			break
		}
		c = s.next()
	}
	return
}

// parseWhere parse where clause, the clause' type must be one of if, SQL condition
func (s *Sections) parseWhere() (res WhereClause, err error) {
	c := s.current()
	res.VarName = s.GetName(c.Type)
	s.tmplAppend(res.Creat())
	res.Type = c.Type

	if !s.HasMore() {
		return
	}
	c = s.next()
	for {
		switch c.Type {
		case model.SQL, model.DATA, model.VARIABLE:
			sqlClause := s.parseSQL(res.VarName)
			res.Value = append(res.Value, sqlClause)
			s.tmplAppend(sqlClause.Finish())
		case model.IF:
			var ifClause IfClause
			ifClause, err = s.parseIF(res.VarName)
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.tmplAppend(ifClause.Finish())
		case model.FOR:
			var forClause ForClause
			forClause, err = s.parseFor(res.VarName)
			if err != nil {
				return
			}
			res.Value = append(res.Value, forClause)
			s.tmplAppend(forClause.Finish())
		case model.END:
			return
		default:
			err = fmt.Errorf("unknow clause : %s", c.Origin)
			return
		}
		if !s.HasMore() {
			break
		}
		c = s.next()
	}
	if c.isEnd() {
		return
	}
	err = fmt.Errorf("incomplete SQL,where not end")
	return
}

// parseSet parse set clause, the clause' type must be one of if, SQL condition
func (s *Sections) parseSet() (res SetClause, err error) {
	c := s.current()
	res.VarName = s.GetName(c.Type)
	s.tmplAppend(res.Creat())
	if !s.HasMore() {
		return
	}
	c = s.next()

	res.Type = c.Type
	for {
		switch c.Type {
		case model.SQL, model.DATA, model.VARIABLE:
			sqlClause := s.parseSQL(res.VarName)
			res.Value = append(res.Value, sqlClause)
			s.tmplAppend(sqlClause.Finish())
		case model.IF:
			var ifClause IfClause
			ifClause, err = s.parseIF(res.VarName)
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.tmplAppend(ifClause.Finish())
		case model.FOR:
			var forClause ForClause
			forClause, err = s.parseFor(res.VarName)
			if err != nil {
				return
			}
			res.Value = append(res.Value, forClause)
			s.tmplAppend(forClause.Finish())
		case model.END:
			return
		default:
			err = fmt.Errorf("unknow clause : %s", c.Origin)
			return
		}
		if !s.HasMore() {
			break
		}
		c = s.next()
	}
	if c.isEnd() {
		err = fmt.Errorf("incomplete SQL,set not end")
	}
	return
}
func (s *Sections) parseFor(name string) (res ForClause, err error) {
	c := s.current()
	res.forSlice = c
	s.tmplAppend(res.Creat())
	s.forValue = append(s.forValue, res.forSlice.ForRange)

	if !s.HasMore() {
		return
	}
	c = s.next()
	for {
		switch c.Type {
		case model.SQL, model.DATA, model.VARIABLE:
			strClause := s.parseSQL(name)
			res.Value = append(res.Value, strClause)
			s.tmplAppend(fmt.Sprintf("%s.WriteString(%s)", name, strClause.String()))
		case model.IF:
			var ifClause IfClause
			ifClause, err = s.parseIF(name)
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.tmplAppend(ifClause.Finish())
		case model.FOR:
			var forClause ForClause
			forClause, err = s.parseFor(name)
			if err != nil {
				return
			}
			res.Value = append(res.Value, forClause)
			s.tmplAppend(forClause.Finish())
		case model.END:
			s.forValue = s.forValue[:len(s.forValue)-1]
			return
		default:
			err = fmt.Errorf("unknow clause : %s", c.Origin)
			return
		}
		if !s.HasMore() {
			break
		}
		c = s.next()
	}
	if c.isEnd() {
		err = fmt.Errorf("incomplete SQL,set not end")
	}
	return
}

// parseSQL parse sql condition, the clause' type must be one of SQL condition, VARIABLE, Data
func (s *Sections) parseSQL(name string) (res SQLClause) {
	res.VarName = name
	res.Type = model.SQL
	for {
		c := s.current()
		switch c.Type {
		case model.SQL:
			res.Value = append(res.Value, c.Value)
		case model.VARIABLE:
			res.Value = append(res.Value, c.Value)
		case model.DATA:
			forKeyValue, isInForRange := s.isInForValue(c.Value)
			if isInForRange {
				forDataName := fmt.Sprintf("%sFor%s_%%d", strings.Replace(c.Value, ".", "", -1), strings.Title(name))
				c.Value = fmt.Sprintf("fmt.Sprintf(\"@%s\",%s)", forDataName, forKeyValue.key)
				s.tmplAppend(fmt.Sprintf("params[fmt.Sprintf(\"%s\",%s)]=%s%s", forDataName, forKeyValue.key, forKeyValue.value, forKeyValue.suffix))
			} else {
				c.Value = strconv.Quote("@" + c.Value)
			}
			res.Value = append(res.Value, c.Value)
		default:
			s.SubIndex()
			return
		}
		if !s.HasMore() {
			return
		}
		c = s.next()
	}
}

// checkSQLVar check sql variable by for loops value and external params
func (s *Sections) checkSQLVar(param string, status model.Status, method *InterfaceMethod) (result section, err error) {
	paramName := strings.Split(param, ".")[0]
	for index, part := range s.members {
		if part.Type == model.FOR && part.ForRange.value == paramName {
			switch status {
			case model.DATA:
				method.HasForParams = true
				if part.ForRange.key == "_" {
					s.members[index].SetForRangeKey("_index")
				}
			case model.VARIABLE:
				param = fmt.Sprintf("%s.Quote(%s)", method.S, param)
			}
			result = section{
				Type:  status,
				Value: param,
			}
			return
		}

	}

	return method.checkSQLVarByParams(param, status)
}

// GetName ...
func (s *Sections) GetName(status model.Status) string {
	switch status {
	case model.WHERE:
		defer func() { s.ClauseTotal[model.WHERE]++ }()
		return fmt.Sprintf("whereClause%d", s.ClauseTotal[model.WHERE])
	case model.SET:
		defer func() { s.ClauseTotal[model.SET]++ }()
		return fmt.Sprintf("setClause%d", s.ClauseTotal[model.SET])
	default:
		return "generateSQL"
	}
}

// sql fragment
type fragment struct {
	Type    model.Status
	Value   string
	IsArray bool
}

func (f *fragment) fragmentByParams(params []parser.Param) (str string) {
	str = strings.ToLower(f.Value)
	for _, param := range params {
		if param.Name == f.Value {
			f.IsArray = param.IsArray
			switch param.Type {
			case "bool":
				f.Type = model.BOOL
				return
			case "int", "int8", "int16", "int32", "int64", "float32", "float64":
				f.Type = model.DIGIT
				return
			case "string":
				f.Type = model.STRING
				return
			case "Time":
				f.Type = model.TIME
				return
			default:
				f.Type = model.OTHER
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
		switch list[i].Type {
		case model.IF, model.ELSE, model.END, model.BOOL, model.LOGICAL, model.WHERE, model.SET:
			continue
		case model.DIGIT, model.STRING, model.OTHER, model.TIME, model.NIL:
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
	if expr[1].Type != model.EXPRESSION {
		return false
	}
	//Only arrays can be compared with nil
	if expr[0].Type == model.NIL || expr[2].Type == model.NIL {
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
