package check

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gen/internal/parser"
)

type slice struct {
	Type   Status
	Value  string
	Origin string
}

// InterfaceMethod interface's method
type InterfaceMethod struct {
	Doc           string
	S             string
	OriginStruct  parser.Param
	TargetStruct  string
	MethodName    string
	Params        []parser.Param
	Result        []parser.Param
	ResultData    parser.Param
	ExecuteResult string
	SqlTmplList   []string
	SqlData       []string
	SqlString     string
	GormOption    string
	Table         string // specified by user. if empty, generate it with gorm
	InterfaceName string
	Package       string
}

// HasSqlData has variable or not
func (m *InterfaceMethod) HasSqlData() bool {
	return len(m.SqlData) > 0
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

// paramToString param list to string used in tmpl
func paramToString(params []parser.Param) string {
	var res []string
	for _, param := range params {
		tmplString := fmt.Sprintf("%s ", param.Name)
		if param.IsArray {
			tmplString += "[]"
		}
		if param.IsPointer {
			tmplString += "*"
		}
		if param.Package != "" {
			tmplString += fmt.Sprintf("%s.", param.Package)
		}
		tmplString += param.Type
		res = append(res, tmplString)
	}
	return strings.Join(res, ",")
}

// checkParams check all parameters
func (m *InterfaceMethod) checkMethod(methods []*InterfaceMethod, s *BaseStruct) (err error) {
	for _, method := range methods {
		if m.IsRepeatFromDifferentInterface(method) {
			return fmt.Errorf("can not generate method with the same name from different interface:[%s.%s] and [%s.%s]",
				m.InterfaceName, m.MethodName, method.InterfaceName, method.MethodName)
		}
	}

	for _, member := range s.Members {
		if member.Name == m.MethodName {
			return fmt.Errorf("can not generate method same name with struct member:[%s.%s] and [%s.%s]",
				m.InterfaceName, m.MethodName, s.StructName, member.Name)
		}
	}

	return nil
}

// checkParams check all parameters
func (m *InterfaceMethod) checkParams(params []parser.Param) (err error) {
	paramList := make([]parser.Param, len(params))
	for i, param := range params {
		if param.Package == "UNDEFINED" {
			param.Package = m.OriginStruct.Package
		}
		if param.IsMap() || param.IsGenM() || param.IsError() || param.IsNull() {
			return fmt.Errorf("type error on interface [%s] param: [%s]", m.InterfaceName, param.Name)
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
			m.ExecuteResult = "err"
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
		default:
			if !m.ResultData.IsNull() {
				return fmt.Errorf("query method cannot return more than 1 data value in [%s.%s]", m.InterfaceName, m.MethodName)
			}
			if param.Package == "" && !(param.AllowType() || param.IsMap() || param.IsTime()) {
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
	m.SqlString = m.parseDocString()
	if err = m.sqlStateCheck(); err != nil {
		err = fmt.Errorf("interface %s member method %s check sql err:%w", m.InterfaceName, m.MethodName, err)
	}
	return
}

func (m *InterfaceMethod) parseDocString() string {
	docString := strings.TrimSpace(m.Doc)
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

// sqlStateCheck check sql with an adeterministic finite automaton
func (m *InterfaceMethod) sqlStateCheck() error {
	sqlString := m.SqlString
	result := NewSlices()
	var buf sql
	for i := 0; !strOutrange(i, sqlString); i++ {
		b := sqlString[i]
		switch b {
		case '"':
			_ = buf.WriteByte(sqlString[i])
			for i++; ; i++ {
				if strOutrange(i, sqlString) {
					return fmt.Errorf("incomplete SQL:%s", sqlString)
				}
				_ = buf.WriteByte(sqlString[i])
				if sqlString[i] == '"' && sqlString[i-1] != '\\' {
					break
				}
			}
		case '{', '@':
			if sqlClause := buf.Dump(); strings.TrimSpace(sqlClause) != "" {
				result.slices = append(result.slices, slice{
					Type:  SQL,
					Value: strconv.Quote(sqlClause),
				})
			}

			if strOutrange(i+1, sqlString) {
				return fmt.Errorf("incomplete SQL:%s", sqlString)
			}
			if b == '{' && sqlString[i+1] == '{' {
				for i += 2; ; i++ {
					if strOutrange(i, sqlString) {
						return fmt.Errorf("incomplete SQL:%s", sqlString)
					}
					if sqlString[i] == '"' {
						_ = buf.WriteByte(sqlString[i])
						for i++; ; i++ {
							if strOutrange(i, sqlString) {
								return fmt.Errorf("incomplete SQL:%s", sqlString)
							}
							_ = buf.WriteByte(sqlString[i])
							if sqlString[i] == '"' && sqlString[i-1] != '\\' {
								break
							}
						}
						i++
					}

					if strOutrange(i+1, sqlString) {
						return fmt.Errorf("incomplete SQL:%s", sqlString)
					}
					if sqlString[i] == '}' && sqlString[i+1] == '}' {
						i++

						sqlClause := buf.Dump()
						part, err := checkTemplate(sqlClause, m.Params)
						if err != nil {
							return fmt.Errorf("sql [%s] dynamic template %s err:%w", sqlString, sqlClause, err)
						}
						result.slices = append(result.slices, part)
						break
					}
					buf.WriteSql(sqlString[i])
				}
			}
			if b == '@' {
				i++
				status := DATA
				if sqlString[i] == '@' {
					i++
					status = VARIABLE
				}
				for ; ; i++ {
					if strOutrange(i, sqlString) || isEnd(sqlString[i]) {
						varString := buf.Dump()
						params, err := m.methodParams(varString, status)
						if err != nil {
							return fmt.Errorf("sql [%s] varable %s err:%s", sqlString, varString, err)
						}
						result.slices = append(result.slices, params)
						i--
						break
					}
					buf.WriteSql(sqlString[i])
				}
			}
		default:
			buf.WriteSql(b)
		}
	}
	if sqlClause := buf.Dump(); strings.TrimSpace(sqlClause) != "" {
		result.slices = append(result.slices, slice{
			Type:  SQL,
			Value: strconv.Quote(sqlClause),
		})
	}

	_, err := result.parse()
	if err != nil {
		return fmt.Errorf("sql [%s] parser err:%w", sqlString, err)
	}
	m.SqlTmplList = result.tmpl
	return nil
}

// methodParams return extrenal parameters, table name
func (m *InterfaceMethod) methodParams(param string, s Status) (result slice, err error) {
	for _, p := range m.Params {
		if p.Name == param {
			var str string
			switch s {
			case DATA:
				str = fmt.Sprintf("\"@%s\"", param)
				if ok := m.isParamExist(param); !ok {
					m.SqlData = append(m.SqlData, param)
				}
			case VARIABLE:
				if p.Type != "string" {
					err = fmt.Errorf("variable name must be string :%s type is %s", param, p.Type)
				}
				str = fmt.Sprintf("%s.Quote(%s)", m.S, param)
			}
			result = slice{
				Type:  s,
				Value: str,
			}
			return
		}
	}
	if param == "table" {
		result = slice{
			Type:  SQL,
			Value: strconv.Quote(m.Table),
		}
		return
	}
	return result, fmt.Errorf("unknow variable param:%s", param)
}

// isParamExist check param duplicate
func (m *InterfaceMethod) isParamExist(paramName string) bool {
	for _, param := range m.SqlData {
		if param == paramName {
			return true
		}
	}
	return false
}

// checkTemplate check sql template's syntax (check if/else/where/set)
func checkTemplate(tmpl string, params []parser.Param) (result slice, err error) {
	fragmentList, err := splitTemplate(tmpl, params)
	if err != nil {
		return
	}
	err = checkTempleFragmentValid(fragmentList)
	if err != nil {
		return
	}
	return fragmentToSLice(fragmentList)
}
