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
	MethodStruct  string
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
}

// HasSqlData has variable or not
func (f *InterfaceMethod) HasSqlData() bool {
	return len(f.SqlData) > 0
}

// HasGotPoint parameter has pointer or not
func (f *InterfaceMethod) HasGotPoint() bool {
	return !f.HasNeedNewResult()
}

// HasNeedNewResult need pointer or not
func (f *InterfaceMethod) HasNeedNewResult() bool {
	return !f.ResultData.IsArray && f.ResultData.IsNull() && f.ResultData.IsTime()
}

// checkParams check all parameters
func (f *InterfaceMethod) checkParams(params []parser.Param) (err error) {
	paramList := make([]parser.Param, len(params))
	for i, r := range params {
		if r.Package == "UNDEFINED" {
			r.Package = f.OriginStruct.Package
		}
		paramList[i] = r
	}
	f.Params = paramList
	return
}

// checkResult check all parameters and replace gen.T by target structure. Parameters must be one of int/string/struct
func (f *InterfaceMethod) checkResult(result []parser.Param) (err error) {
	resList := make([]parser.Param, len(result))
	for i, param := range result {
		if param.Package == "UNDEFINED" {
			param.Package = f.OriginStruct.Package
		}
		switch {
		case param.IsError():
			param.SetName("err")
			f.ExecuteResult = "err"
		case param.Eq(f.OriginStruct) || param.IsGenT():
			param.SetName("result")
			param.Type = f.OriginStruct.Type
			param.Package = f.OriginStruct.Package
			param.IsPointer = true
			f.ResultData = param
		case param.AllowType(), param.IsTime():
			param.SetName("result")
			f.ResultData = param
		default:
			return fmt.Errorf("illegal parameter：%s.%s on struct %s.%s generated method %s", param.Package, param.Type, f.OriginStruct.Package, f.OriginStruct.Type, f.MethodName)
		}
		resList[i] = param
	}
	f.Result = resList
	return
}

// checkSQL get sql from comment and check it
func (f *InterfaceMethod) checkSQL() (err error) {
	f.SqlString = f.parseDocString()
	if err = f.sqlStateCheck(); err != nil {
		err = fmt.Errorf("interface %s member method %s check sql err:%w", f.InterfaceName, f.MethodName, err)
	}
	return
}

func (f *InterfaceMethod) parseDocString() string {
	docString := strings.TrimSpace(f.Doc)
	switch {
	case strings.HasPrefix(strings.ToLower(docString), "sql("):
		docString = docString[4 : len(docString)-1]
		f.GormOption = "Raw"
		if f.ResultData.IsNull() {
			f.GormOption = "Exec"
		}
	case strings.HasPrefix(strings.ToLower(docString), "where("):
		docString = docString[6 : len(docString)-1]
		f.GormOption = "Where"
	default:
		f.GormOption = "Raw"
		if f.ResultData.IsNull() {
			f.GormOption = "Exec"
		}
	}

	// if wrapped by ", trim it
	if strings.HasPrefix(docString, `"`) && strings.HasSuffix(docString, `"`) {
		docString = docString[1 : len(docString)-1]
	}
	return docString
}

// sqlStateCheck check sql with an adeterministic finite automaton
func (f *InterfaceMethod) sqlStateCheck() error {
	sqlString := f.SqlString
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
						part, err := checkTemplate(sqlClause, f.Params)
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
						params, err := f.methodParams(varString, status)
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
	f.SqlTmplList = result.tmpl
	return nil
}

// methodParams return extrenal parameters, table name
func (f *InterfaceMethod) methodParams(param string, s Status) (result slice, err error) {
	for _, p := range f.Params {
		if p.Name == param {
			var str string
			switch s {
			case DATA:
				str = fmt.Sprintf("\"@%s\"", param)
			case VARIABLE:
				if p.Type != "string" {
					err = fmt.Errorf("variable name must be string :%s type is %s", param, p.Type)
				}
				str = fmt.Sprintf("%s.Quote(%s)", f.S, param)
			}
			f.SqlData = append(f.SqlData, param)
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
			Value: strconv.Quote(f.Table),
		}
		return
	}
	return result, fmt.Errorf("unknow variable param:%s", param)
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

// DupAppend append methon and duplicate
func (f *InterfaceMethod) DupAppend(methods []*InterfaceMethod) []*InterfaceMethod {
	for index, method := range methods {
		if method.MethodName == f.MethodName {
			methods[index] = f
			return methods
		}
	}
	return append(methods, f)
}
