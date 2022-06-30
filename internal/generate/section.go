package generate

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gen/internal/model"
)

// NewSection create and initialize Sections
func NewSection() *Section {
	return &Section{
		ClauseTotal: map[model.Status]int{
			model.WHERE: 0,
			model.SET:   0,
		},
	}
}

// Section split sql into chunks
type Section struct {
	members      []section
	Tmpls        []string
	currentIndex int
	ClauseTotal  map[model.Status]int
	forValue     []ForRange
}

// next return next section and increase index by 1
func (s *Section) next() section {
	if s.currentIndex < len(s.members)-1 {
		s.currentIndex++
		return s.members[s.currentIndex]
	}
	return section{Type: model.END}
}

// SubIndex take index one step back
func (s *Section) SubIndex() {
	s.currentIndex--
}

// HasMore is has more section
func (s *Section) HasMore() bool {
	return s.currentIndex < len(s.members)-1
}

// IsNull whether section is empty
func (s *Section) IsNull() bool {
	return len(s.members) == 0
}

// current return current section
func (s *Section) current() section {
	return s.members[s.currentIndex]
}

func (s *Section) appendTmpl(value string) {
	s.Tmpls = append(s.Tmpls, value)
}

func (s *Section) isInForValue(value string) (ForRange, bool) {
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

func (s *Section) hasSameName(value string) bool {
	for _, p := range s.members {
		if p.Type == model.FOR && p.ForRange.value == value {
			return true
		}
	}
	return false
}

// BuildSQL sql sections and append to tmpl, return a Clause array
func (s *Section) BuildSQL() ([]Clause, error) {
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
			s.appendTmpl(sqlClause.Finish())
		case model.IF:
			ifClause, err := s.parseIF(name)
			if err != nil {
				return nil, err
			}
			res = append(res, ifClause)
			s.appendTmpl(ifClause.Finish())
		case model.WHERE:
			whereClause, err := s.parseWhere()
			if err != nil {
				return nil, err
			}
			res = append(res, whereClause)
			s.appendTmpl(whereClause.Finish(name))
		case model.SET:
			setClause, err := s.parseSet()
			if err != nil {
				return nil, err
			}
			res = append(res, setClause)
			s.appendTmpl(setClause.Finish(name))
		case model.FOR:
			forClause, err := s.parseFor(name)
			_, _ = forClause, err
			if err != nil {
				return nil, err
			}
			res = append(res, forClause)
			s.appendTmpl(forClause.Finish())
		case model.END:
		default:
			return nil, fmt.Errorf("unknow clause:%s", c.Value)
		}
		if !s.HasMore() {
			break
		}
		c = s.next()
	}
	return res, nil
}

// parseIF parse if clause
func (s *Section) parseIF(name string) (res IfClause, err error) {
	c := s.current()
	res.slice = c

	s.appendTmpl(res.Create())
	if !s.HasMore() {
		return
	}
	c = s.next()
	for {
		switch c.Type {
		case model.SQL, model.DATA, model.VARIABLE:
			sqlClause := s.parseSQL(name)
			res.Value = append(res.Value, sqlClause)
			s.appendTmpl(sqlClause.Finish())
		case model.IF:
			var ifClause IfClause
			ifClause, err = s.parseIF(name)
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendTmpl(ifClause.Finish())
		case model.WHERE:
			var whereClause WhereClause
			whereClause, err = s.parseWhere()
			if err != nil {
				return
			}
			res.Value = append(res.Value, whereClause)
			s.appendTmpl(whereClause.Finish(name))
		case model.SET:
			var setClause SetClause
			setClause, err = s.parseSet()
			if err != nil {
				return
			}
			res.Value = append(res.Value, setClause)
			s.appendTmpl(setClause.Finish(name))
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
			s.appendTmpl(res.Finish())
		case model.END:
			return
		default:
			err = fmt.Errorf("unknow clause : %s", c.Value)
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
func (s *Section) parseElSE(name string) (res ElseClause, err error) {
	res.slice = s.current()
	s.appendTmpl(res.Create())

	if !s.HasMore() {
		return
	}
	c := s.next()
	for {
		switch c.Type {
		case model.SQL, model.DATA, model.VARIABLE:
			sqlClause := s.parseSQL(name)
			res.Value = append(res.Value, sqlClause)
			s.appendTmpl(sqlClause.Create())
		case model.IF:
			var ifClause IfClause
			ifClause, err = s.parseIF(name)
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendTmpl(ifClause.Finish())
		case model.WHERE:
			var whereClause WhereClause
			whereClause, err = s.parseWhere()
			if err != nil {
				return
			}
			res.Value = append(res.Value, whereClause)
			s.appendTmpl(whereClause.Finish(name))
		case model.SET:
			var setClause SetClause
			setClause, err = s.parseSet()
			if err != nil {
				return
			}
			res.Value = append(res.Value, setClause)
			s.appendTmpl(setClause.Finish(name))
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
			s.appendTmpl(forClause.Finish())
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
func (s *Section) parseWhere() (res WhereClause, err error) {
	c := s.current()
	res.VarName = s.GetName(c.Type)
	s.appendTmpl(res.Create())
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
			s.appendTmpl(sqlClause.Finish())
		case model.IF:
			var ifClause IfClause
			ifClause, err = s.parseIF(res.VarName)
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendTmpl(ifClause.Finish())
		case model.FOR:
			var forClause ForClause
			forClause, err = s.parseFor(res.VarName)
			if err != nil {
				return
			}
			res.Value = append(res.Value, forClause)
			s.appendTmpl(forClause.Finish())
		case model.END:
			return
		default:
			err = fmt.Errorf("unknow clause : %s", c.Value)
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
func (s *Section) parseSet() (res SetClause, err error) {
	c := s.current()
	res.VarName = s.GetName(c.Type)
	s.appendTmpl(res.Create())
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
			s.appendTmpl(sqlClause.Finish())
		case model.IF:
			var ifClause IfClause
			ifClause, err = s.parseIF(res.VarName)
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendTmpl(ifClause.Finish())
		case model.FOR:
			var forClause ForClause
			forClause, err = s.parseFor(res.VarName)
			if err != nil {
				return
			}
			res.Value = append(res.Value, forClause)
			s.appendTmpl(forClause.Finish())
		case model.END:
			return
		default:
			err = fmt.Errorf("unknow clause : %s", c.Value)
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

func (s *Section) parseFor(name string) (res ForClause, err error) {
	c := s.current()
	res.forSlice = c
	s.appendTmpl(res.Create())
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
			s.appendTmpl(fmt.Sprintf("%s.WriteString(%s)", name, strClause.String()))
		case model.IF:
			var ifClause IfClause
			ifClause, err = s.parseIF(name)
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendTmpl(ifClause.Finish())
		case model.FOR:
			var forClause ForClause
			forClause, err = s.parseFor(name)
			if err != nil {
				return
			}
			res.Value = append(res.Value, forClause)
			s.appendTmpl(forClause.Finish())
		case model.END:
			s.forValue = s.forValue[:len(s.forValue)-1]
			return
		default:
			err = fmt.Errorf("unknow clause : %s", c.Value)
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
func (s *Section) parseSQL(name string) (res SQLClause) {
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
			forRange, isInForRange := s.isInForValue(c.Value)
			if isInForRange {
				s.appendTmpl(forRange.appendDataToParams(c.Value, name))
				c.Value = forRange.DataValue(c.Value, name)
			} else {
				s.appendTmpl(c.AddDataToParamMap())
				c.Value = strconv.Quote("@" + c.SQLParamName())
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
func (s *Section) checkSQLVar(param string, status model.Status, method *InterfaceMethod) (result section, err error) {
	paramName := strings.Split(param, ".")[0]
	for index, part := range s.members {
		if part.Type == model.FOR && part.ForRange.value == paramName {
			switch status {
			case model.DATA:
				method.HasForParams = true
				if part.ForRange.index == "_" {
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
func (s *Section) GetName(status model.Status) string {
	switch status {
	case model.WHERE:
		defer func() { s.ClauseTotal[model.WHERE]++ }()
		return fmt.Sprintf("whereSQL%d", s.ClauseTotal[model.WHERE])
	case model.SET:
		defer func() { s.ClauseTotal[model.SET]++ }()
		return fmt.Sprintf("setSQL%d", s.ClauseTotal[model.SET])
	default:
		return "generateSQL"
	}
}

// checkTemplate check sql template's syntax (if/else/where/set/for)
func (s *Section) checkTemplate(tmpl string) (part section, err error) {
	part.Value = tmpl
	part.SQLSlice = s
	part.splitTemplate()

	err = part.checkTemplate()

	return
}

type section struct {
	Type      model.Status
	Value     string
	ForRange  ForRange
	SQLSlice  *Section
	splitList []string
}

func (s *section) isEnd() bool {
	return s.Type == model.END
}

func (s *section) String() string {
	if s.Type == model.FOR {
		return s.ForRange.String()
	}
	return s.Value
}

func (s *section) splitTemplate() {
	s.splitList = strings.FieldsFunc(strings.TrimSpace(s.Value), func(r rune) bool {
		return r == ':' || r == ' ' || r == '=' || r == ','
	})
}

func (s *section) checkTemplate() error {
	if len(s.splitList) == 0 {
		return fmt.Errorf("template is null")
	}
	if model.GenKeywords.Contain(s.Value) {
		return fmt.Errorf("template can not use gen keywords")
	}

	err := s.sectionType(s.splitList[0])
	if err != nil {
		return err
	}

	if s.Type == model.FOR {
		if len(s.splitList) != 5 {
			return fmt.Errorf("for range syntax error: %s", s.Value)
		}
		if s.SQLSlice.hasSameName(s.splitList[2]) {
			return fmt.Errorf("cannot use the same value name in different for loops")
		}
		s.ForRange.index = s.splitList[1]
		s.ForRange.value = s.splitList[2]
		s.ForRange.rangeList = s.splitList[4]
	}
	return nil
}

func (s *section) sectionType(str string) error {
	switch str {
	case "if":
		s.Type = model.IF
	case "else":
		s.Type = model.ELSE
	case "for":
		s.Type = model.FOR
	case "where":
		s.Type = model.WHERE
	case "set":
		s.Type = model.SET
	case "end":
		s.Type = model.END
	default:
		return fmt.Errorf("unknown syntax: %s", str)
	}
	return nil
}

func (s *section) SetForRangeKey(key string) {
	s.ForRange.index = key
	s.Value = s.String()
}

func (s *section) AddDataToParamMap() string {
	return fmt.Sprintf("params[%q] = %s", s.SQLParamName(), s.Value)
}

func (s *section) SQLParamName() string {
	return strings.Replace(s.Value, ".", "", -1)
}

// ForRange for range clause for diy method
type ForRange struct {
	index     string
	value     string
	suffix    string
	rangeList string
}

func (f *ForRange) String() string {
	return fmt.Sprintf("for %s, %s := range %s", f.index, f.value, f.rangeList)
}

func (f *ForRange) mapIndexName(prefix, dataName, clauseName string) string {
	return fmt.Sprintf("\"%s%sFor%s_\"+strconv.Itoa(%s)", prefix, strings.Replace(dataName, ".", "", -1), strings.Title(clauseName), f.index)
}

// DataValue return data value
func (f *ForRange) DataValue(dataName, clauseName string) string {
	return f.mapIndexName("@", dataName, clauseName)
}

func (f *ForRange) appendDataToParams(dataName, clauseName string) string {
	return fmt.Sprintf("params[%s]=%s%s", f.mapIndexName("", dataName, clauseName), f.value, f.suffix)
}
