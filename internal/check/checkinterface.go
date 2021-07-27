package check

import (
	"bytes"
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
	SqlList       []string
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
	retList := make([]parser.Param, len(result))
	for i, param := range result {
		if param.Package == "UNDEFINED" {
			param.Package = f.OriginStruct.Package
		}
		switch {
		case param.IsError():
			param.SetName("err")
			f.ExecuteResult = "err"
		case param.Eq(f.OriginStruct) || param.IsGenT():
			param.Type = f.OriginStruct.Type
			param.Package = f.OriginStruct.Package
			param.SetName("result")
			param.IsPointer = true
			f.ResultData = param
		case param.AllowType(), param.IsTime():
			param.SetName("result")
			f.ResultData = param
		default:
			return fmt.Errorf("illegal parameter：%s.%s on struct %s.%s generated method %s \n ", param.Package, param.Type, f.OriginStruct.Package, f.OriginStruct.Type, f.MethodName)
		}
		retList[i] = param
	}
	f.Result = retList
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
		docString = docString[1 : len(f.SqlString)-1]
	}
	return docString
}

type sql struct{ bytes.Buffer }

func (s *sql) WriteSql(b byte) {
	switch b {
	case '\n', '\t', ' ':
		if s.Len() == 0 || s.Bytes()[s.Len()-1] != ' ' {
			_ = s.WriteByte(' ')
		}
	default:
		_ = s.WriteByte(b)
	}

}

// sqlStateCheck check sql with an adeterministic finite automaton
func (f *InterfaceMethod) sqlStateCheck() (err error) {
	sqlString := f.SqlString + " "
	result := NewSlices()
	var out sql

	for i := 0; i < len(sqlString); i++ {
		b := sqlString[i]
		switch b {
		case '"':
			_ = out.WriteByte(sqlString[i])
			for {
				i++
				if !stringHasMore(i, sqlString) {
					return fmt.Errorf("incomplete SQL:%s", sqlString)
				}
				_ = out.WriteByte(sqlString[i])
				if sqlString[i] == '"' && sqlString[i-1] != '\\' {
					break
				}
			}

		case '{', '@':
			sqlClause := out.String()
			if strings.TrimSpace(sqlClause) != "" {
				result.slices = append(result.slices, slice{
					Type:  SQL,
					Value: strconv.Quote(out.String()),
				})

			}
			out.Reset()

			if !stringHasMore(i+1, sqlString) {
				return fmt.Errorf("incomplete SQL:%s", sqlString)
			}
			if b == '{' && sqlString[i+1] == '{' {
				i += 2
				for {
					if !stringHasMore(i, sqlString) {
						return fmt.Errorf("incomplete SQL:%s", sqlString)
					}
					if sqlString[i] == '"' {
						_ = out.WriteByte(sqlString[i])
						for {
							i++
							if !stringHasMore(i, sqlString) {
								return fmt.Errorf("incomplete SQL:%s", sqlString)
							}
							_ = out.WriteByte(sqlString[i])
							if sqlString[i] == '"' && sqlString[i-1] != '\\' {
								i++
								break
							}
						}

					}
					if sqlString[i] == '}' && sqlString[i+1] == '}' {
						part, err := checkTemplate(out.String(), f.Params)
						if err != nil {
							err := fmt.Errorf("sql [%s] dynamic template %s err:%s  ", sqlString, out.String(), err)
							return err
						}
						result.slices = append(result.slices, part)
						i++
						out.Reset()
						break
					}
					out.WriteSql(sqlString[i])
					i++

				}

			}
			if b == '@' {
				t := DATA
				i++
				if sqlString[i] == '@' {
					t = VARIABLE
					i++
				}
				for {
					if i == len(sqlString) || isStringEnd(sqlString[i]) {
						varString := out.String()
						s, err := f.isParamInMethod(varString, t)
						if err != nil {
							return fmt.Errorf("sql [%s] varable %s err:%s", sqlString, varString, err)
						}

						result.slices = append(result.slices, s)
						out.Reset()
						i--
						break
					}

					out.WriteSql(sqlString[i])
					i++
				}
			}
		default:
			out.WriteSql(b)
		}
	}
	if strings.Trim(out.String(), " ") != "" {
		result.slices = append(result.slices, slice{
			Type:  SQL,
			Value: strconv.Quote(out.String()),
		})
	}

	_, err = result.parse()
	if err != nil {
		return fmt.Errorf("sql [%s] parser err:%s", sqlString, err)
	}
	f.SqlList = result.tmpl
	return
}

// isParamInMethod if sql's param is not external(except table), generate by gorm
func (f *InterfaceMethod) isParamInMethod(param string, s Status) (result slice, err error) {
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
				str = fmt.Sprintf("helper.Quote(%s)", param)
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
