package generate

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gen/internal/model"
	"gorm.io/gen/internal/parser"
)

// InterfaceMethod interface's method
type InterfaceMethod struct { // feature will replace InterfaceMethod to parser.Method
	Doc           string         // comment
	S             string         // First letter of
	OriginStruct  parser.Param   // origin struct name
	TargetStruct  string         // generated query struct bane
	MethodName    string         // generated function name
	Params        []parser.Param // function input params
	Result        []parser.Param // function output params
	ResultData    parser.Param   // output data
	Section       *Section       // Parse split SQL into sections
	SQLParams     []parser.Param // variable in sql need function input
	SQLString     string         // SQL
	GormOption    string         // gorm execute method Find or Exec or Take
	Table         string         // specified by user. if empty, generate it with gorm
	InterfaceName string         // origin interface name
	Package       string         // interface package name
	HasForParams  bool           //
}

// FuncSign function signature
func (m *InterfaceMethod) FuncSign() string {
	return fmt.Sprintf("%s(%s) (%s)", m.MethodName, m.GetParamInTmpl(), m.GetResultParamInTmpl())
}

// HasSQLData has variable or for params will creat params map
func (m *InterfaceMethod) HasSQLData() bool {
	return len(m.SQLParams) > 0 || m.HasForParams
}

// HasGotPoint parameter has pointer or not
func (m *InterfaceMethod) HasGotPoint() bool {
	return !m.HasNeedNewResult()
}

// HasNeedNewResult need pointer or not
func (m *InterfaceMethod) HasNeedNewResult() bool {
	return !m.ResultData.IsArray && ((m.ResultData.IsNull() && m.ResultData.IsTime()) || m.ResultData.IsMap())
}

// GormRunMethodName return single data use Take() return array use Find
func (m *InterfaceMethod) GormRunMethodName() string {
	if m.ResultData.IsArray {
		return "Find"
	}
	return "Take"
}

// ReturnRowsAffected return rows affected
func (m *InterfaceMethod) ReturnRowsAffected() bool {
	for _, res := range m.Result {
		if res.Name == "rowsAffected" {
			return true
		}
	}
	return false
}

// ReturnError return error
func (m *InterfaceMethod) ReturnError() bool {
	for _, res := range m.Result {
		if res.IsError() {
			return true
		}
	}
	return false
}

// IsRepeatFromDifferentInterface check different interface has same mame method
func (m *InterfaceMethod) IsRepeatFromDifferentInterface(newMethod *InterfaceMethod) bool {
	return m.MethodName == newMethod.MethodName && m.InterfaceName != newMethod.InterfaceName && m.TargetStruct == newMethod.TargetStruct
}

// IsRepeatFromSameInterface check different interface has same mame method
func (m *InterfaceMethod) IsRepeatFromSameInterface(newMethod *InterfaceMethod) bool {
	return m.MethodName == newMethod.MethodName && m.InterfaceName == newMethod.InterfaceName && m.TargetStruct == newMethod.TargetStruct
}

//GetParamInTmpl return param list
func (m *InterfaceMethod) GetParamInTmpl() string {
	return paramToString(m.Params)
}

// GetResultParamInTmpl return result list
func (m *InterfaceMethod) GetResultParamInTmpl() string {
	return paramToString(m.Result)
}

// SQLParamName sql param map key,
func (m *InterfaceMethod) SQLParamName(param string) string {
	return strings.Replace(param, ".", "", -1)
}

// paramToString param list to string used in tmpl
func paramToString(params []parser.Param) string {
	var res []string
	for _, param := range params {
		res = append(res, param.TmplString())
	}
	return strings.Join(res, ",")
}

// DocComment return comment sql add "//" every line
func (m *InterfaceMethod) DocComment() string {
	return strings.Replace(strings.TrimSpace(m.Doc), "\n", "\n//", -1)
}

// checkParams check all parameters
func (m *InterfaceMethod) checkMethod(methods []*InterfaceMethod, s *QueryStructMeta) (err error) {
	if model.GormKeywords.FullMatch(m.MethodName) {
		return fmt.Errorf("can not use keyword as method name:%s", m.MethodName)
	}
	// TODO check methods Always empty?
	for _, method := range methods {
		if m.IsRepeatFromDifferentInterface(method) {
			return fmt.Errorf("can not generate method with the same name from different interface:[%s.%s] and [%s.%s]",
				m.InterfaceName, m.MethodName, method.InterfaceName, method.MethodName)
		}
	}
	for _, f := range s.Fields {
		if f.Name == m.MethodName {
			return fmt.Errorf("can not generate method same name with struct field:[%s.%s] and [%s.%s]",
				m.InterfaceName, m.MethodName, s.ModelStructName, f.Name)
		}
	}

	return nil
}

// checkParams check all parameters
func (m *InterfaceMethod) checkParams(params []parser.Param) (err error) {
	paramList := make([]parser.Param, len(params))
	for i, param := range params {
		switch {
		case param.Package == "UNDEFINED":
			param.Package = m.Package
		case param.IsMap() || param.IsGenM() || param.IsError() || param.IsNull():
			return fmt.Errorf("type error on interface [%s] param: [%s]", m.InterfaceName, param.Name)
		case param.IsGenT():
			param.Type = m.OriginStruct.Type
			param.Package = m.OriginStruct.Package
		}
		paramList[i] = param
	}
	m.Params = paramList
	return
}

// checkResult check all parameters and replace gen.T by target structure. Parameters must be one of int/string/struct/map
func (m *InterfaceMethod) checkResult(result []parser.Param) (err error) {
	resList := make([]parser.Param, len(result))
	var hasError bool
	for i, param := range result {
		if param.Package == "UNDEFINED" {
			param.Package = m.Package
		}
		if param.IsGenM() {
			param.Type = "map[string]interface{}"
			param.Package = ""
		}
		switch {
		case param.InMainPkg():
			return fmt.Errorf("query method cannot return struct of main package in [%s.%s]", m.InterfaceName, m.MethodName)
		case param.IsError():
			if hasError {
				return fmt.Errorf("query method cannot return more than 1 error value in [%s.%s]", m.InterfaceName, m.MethodName)
			}
			param.SetName("err")
			hasError = true
		case param.Eq(m.OriginStruct) || param.IsGenT():
			if !m.ResultData.IsNull() {
				return fmt.Errorf("query method cannot return more than 1 data value in [%s.%s]", m.InterfaceName, m.MethodName)
			}
			param.SetName("result")
			param.Type = m.OriginStruct.Type
			param.Package = m.OriginStruct.Package
			param.IsPointer = true
			m.ResultData = param
		case param.IsInterface():
			return fmt.Errorf("query method can not return interface in [%s.%s]", m.InterfaceName, m.MethodName)
		case param.IsGenRowsAffected():
			param.Type = "int64"
			param.Package = ""
			param.SetName("rowsAffected")
			m.GormOption = "Exec"
		default:
			if !m.ResultData.IsNull() {
				return fmt.Errorf("query method cannot return more than 1 data value in [%s.%s]", m.InterfaceName, m.MethodName)
			}
			if param.Package == "" && !(param.IsBaseType() || param.IsMap() || param.IsTime()) {
				param.Package = m.Package
			}
			param.SetName("result")
			m.ResultData = param
		}
		resList[i] = param
	}
	m.Result = resList
	return
}

// checkSQL get sql from comment and check it
func (m *InterfaceMethod) checkSQL() (err error) {
	m.SQLString = m.parseDocString()
	if err = m.sqlStateCheckAndSplit(); err != nil {
		err = fmt.Errorf("interface %s member method %s check sql err:%w", m.InterfaceName, m.MethodName, err)
	}
	return
}

func (m *InterfaceMethod) parseDocString() string {
	docString := strings.TrimSpace(m.getSQLDocString())
	switch {
	case strings.HasPrefix(strings.ToLower(docString), "sql("):
		docString = docString[4 : len(docString)-1]
		m.GormOption = "Raw"
		if m.ResultData.IsNull() {
			m.GormOption = "Exec"
		}
	case strings.HasPrefix(strings.ToLower(docString), "where("):
		docString = docString[6 : len(docString)-1]
		m.GormOption = "Where"
	default:
		m.GormOption = "Raw"
		if m.ResultData.IsNull() {
			m.GormOption = "Exec"
		}
	}

	// if wrapped by ", trim it
	if strings.HasPrefix(docString, `"`) && strings.HasSuffix(docString, `"`) {
		docString = docString[1 : len(docString)-1]
	}
	return docString
}

func (m *InterfaceMethod) getSQLDocString() string {
	docString := strings.TrimSpace(m.Doc)
	/*
		// methodName descriptive message
		// (this blank line is needed)
		// sql
	*/
	if index := strings.Index(docString, "\n\n"); index != -1 {
		if strings.Contains(docString[index+2:], m.MethodName) {
			docString = docString[:index]
		} else {
			docString = docString[index+2:]
		}
	}
	/* //methodName sql */
	docString = strings.TrimPrefix(docString, m.MethodName)
	// TODO: using sql key word to split comment
	return docString
}

// sqlStateCheckAndSplit check sql with an adeterministic finite automaton
func (m *InterfaceMethod) sqlStateCheckAndSplit() error {
	sqlString := m.SQLString
	m.Section = NewSection()
	var buf model.SQLBuffer
	for i := 0; !strOutRange(i, sqlString); i++ {
		b := sqlString[i]
		switch b {
		case '"':
			_ = buf.WriteByte(sqlString[i])
			for i++; ; i++ {
				if strOutRange(i, sqlString) {
					return fmt.Errorf("incomplete SQL:%s", sqlString)
				}
				_ = buf.WriteByte(sqlString[i])
				if sqlString[i] == '"' && sqlString[i-1] != '\\' {
					break
				}
			}
		case '\'':
			_ = buf.WriteByte(sqlString[i])
			for i++; ; i++ {
				if strOutRange(i, sqlString) {
					return fmt.Errorf("incomplete SQL:%s", sqlString)
				}
				_ = buf.WriteByte(sqlString[i])
				if sqlString[i] == '\'' && sqlString[i-1] != '\\' {
					break
				}
			}
		case '\\':
			if sqlString[i+1] == '@' {
				i++
				buf.WriteSQL(sqlString[i])
				continue
			}
			buf.WriteSQL(b)
		case '{', '@':
			if sqlClause := buf.Dump(); strings.TrimSpace(sqlClause) != "" {
				m.Section.members = append(m.Section.members, section{
					Type:  model.SQL,
					Value: strconv.Quote(sqlClause),
				})
			}

			if strOutRange(i+1, sqlString) {
				return fmt.Errorf("incomplete SQL:%s", sqlString)
			}
			if b == '{' && sqlString[i+1] == '{' {
				for i += 2; ; i++ {
					if strOutRange(i, sqlString) {
						return fmt.Errorf("incomplete SQL:%s", sqlString)
					}
					if sqlString[i] == '"' {
						_ = buf.WriteByte(sqlString[i])
						for i++; ; i++ {
							if strOutRange(i, sqlString) {
								return fmt.Errorf("incomplete SQL:%s", sqlString)
							}
							_ = buf.WriteByte(sqlString[i])
							if sqlString[i] == '"' && sqlString[i-1] != '\\' {
								break
							}
						}
						i++
					}

					if strOutRange(i+1, sqlString) {
						return fmt.Errorf("incomplete SQL:%s", sqlString)
					}
					if sqlString[i] == '}' && sqlString[i+1] == '}' {
						i++
						sqlClause := buf.Dump()
						part, err := m.Section.checkTemplate(sqlClause)
						if err != nil {
							return fmt.Errorf("sql [%s] dynamic template %s err:%w", sqlString, sqlClause, err)
						}
						m.Section.members = append(m.Section.members, part)
						break
					}
					buf.WriteSQL(sqlString[i])
				}
			}
			if b == '@' {
				i++
				status := model.DATA
				if sqlString[i] == '@' {
					i++
					status = model.VARIABLE
				}
				for ; ; i++ {
					if strOutRange(i, sqlString) || isEnd(sqlString[i]) {
						varString := buf.Dump()
						params, err := m.Section.checkSQLVar(varString, status, m)
						if err != nil {
							return fmt.Errorf("sql [%s] varable %s err:%s", sqlString, varString, err)
						}
						m.Section.members = append(m.Section.members, params)
						i--
						break
					}
					buf.WriteSQL(sqlString[i])
				}
			}
		default:
			buf.WriteSQL(b)
		}
	}
	if sqlClause := buf.Dump(); strings.TrimSpace(sqlClause) != "" {
		m.Section.members = append(m.Section.members, section{
			Type:  model.SQL,
			Value: strconv.Quote(sqlClause),
		})
	}

	return nil
}

// checkSQLVarByParams return external parameters, table name
func (m *InterfaceMethod) checkSQLVarByParams(param string, status model.Status) (result section, err error) {
	for _, p := range m.Params {
		structName := strings.Split(param, ".")[0]
		if p.Name == structName {
			if p.Name != param {
				p = parser.Param{
					Name: param,
					Type: "string",
				}
			}
			switch status {
			case model.DATA:
				if !m.isParamExist(param) {
					m.SQLParams = append(m.SQLParams, p)
				}
			case model.VARIABLE:
				if p.Type != "string" || p.IsArray {
					err = fmt.Errorf("variable name must be string :%s type is %s", param, p.TypeName())
					return
				}
				param = fmt.Sprintf("%s.Quote(%s)", m.S, param)
			}
			result = section{
				Type:  status,
				Value: param,
			}
			return
		}
	}
	if param == "table" {
		result = section{
			Type:  model.SQL,
			Value: strconv.Quote(m.Table),
		}
		return
	}

	return result, fmt.Errorf("unknow variable param:%s", param)
}

// isParamExist check param duplicate
func (m *InterfaceMethod) isParamExist(paramName string) bool {
	for _, param := range m.SQLParams {
		if param.Name == paramName {
			return true
		}
	}
	return false
}
