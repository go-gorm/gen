package check

import (
	"fmt"
	"gorm.io/gen/internal/model"
	"gorm.io/gen/internal/parser"
	"strings"
)

type ForRange struct {
	key    string
	value  string
	suffix string
}

type section struct {
	Type     model.Status
	Value    string
	Origin   string
	fList    []fragment
	ForRange ForRange
	SQLSlice *Sections
}

func (s *section) isEnd() bool {
	return s.Type == model.END
}

func (s *section) String() string {
	var values []string
	if len(s.fList) == 0 {
		return ""
	}
	for _, t := range s.fList {
		values = append(values, t.Value)
	}
	return strings.Join(values, " ")
}

func (s *section) checkFragment(str string, params []parser.Param) (f fragment, err error) {
	digital := func(str string) string {
		if isDigit(str) {
			return "<integer>"
		}
		return str
	}

	f = fragment{Type: model.UNKNOWN, Value: strings.Trim(str, " ")}
	str = strings.ToLower(strings.Trim(str, " "))
	switch digital(str) {
	case "<integer>":
		f.Type = model.INT
	case "&&", "||":
		f.Type = model.LOGICAL
	case ">", "<", ">=", "<=", "==", "!=", ":=":
		f.Type = model.EXPRESSION
	case "end":
		f.Type = model.END
		s.Type = model.END
	case "if":
		f.Type = model.IF
		if s.Type != model.ELSE {
			s.Type = model.IF
		}
	case "set":
		f.Type = model.SET
		s.Type = model.SET
	case "else":
		f.Type = model.ELSE
		s.Type = model.ELSE
	case "where":
		f.Type = model.WHERE
		s.Type = model.WHERE
	case "true", "false":
		f.Type = model.BOOL
	case "for":
		f.Type = model.FOR
		s.Type = model.FOR
	case "range":
		f.Type = model.RANGE
	case "nil":
		f.Type = model.NIL
	case f.fragmentByForRange(s):
		if f.Type == model.VALUE {
			s.fList = append(s.fList, fragment{
				Type:  model.OTHER,
				Value: ",",
			})
		}
	case f.fragmentByParams(params):
	default:
		if f.Type != model.UNKNOWN {
			return f, fmt.Errorf("template unknow string: %s", str)
		}

	}
	return
}

func (s *section) splitTemplate(tmpl string, params []parser.Param) (err error) {
	var buf model.SQLBuffer
	var f fragment
	for i := 0; !strOutrange(i, tmpl); i++ {
		switch tmpl[i] {
		case '"':
			_ = buf.WriteByte(tmpl[i])
			for i++; ; i++ {
				if strOutrange(i, tmpl) {
					return fmt.Errorf("incomplete code:%s", tmpl)
				}
				_ = buf.WriteByte(tmpl[i])

				if tmpl[i] == '"' && tmpl[i-1] != '\\' {
					s.fList = append(s.fList, fragment{Type: model.STRING, Value: buf.Dump()})
					break
				}
			}
		case ' ', ',':
			if sqlClause := buf.Dump(); sqlClause != "" {
				f, err = s.checkFragment(sqlClause, params)
				if err != nil {
					return err
				}
				s.fList = append(s.fList, f)
			}

		case '>', '<', '=', '!', ':':
			if sqlClause := buf.Dump(); sqlClause != "" {
				f, err = s.checkFragment(sqlClause, params)
				if err != nil {
					return err
				}
				s.fList = append(s.fList, f)
			}

			_ = buf.WriteByte(tmpl[i])

			if strOutrange(i+1, tmpl) {
				return fmt.Errorf("incomplete code:%s", tmpl)
			}
			if tmpl[i+1] == '=' {
				_ = buf.WriteByte(tmpl[i+1])
				i++
			}

			f, err = s.checkFragment(buf.Dump(), params)
			if err != nil {
				return err
			}
			s.fList = append(s.fList, f)
		case '&', '|':
			if strOutrange(i+1, tmpl) {
				return fmt.Errorf("incomplete code:%s", tmpl)
			}

			if tmpl[i+1] == tmpl[i] {
				i++

				if sqlClause := buf.Dump(); sqlClause != "" {
					f, err = s.checkFragment(sqlClause, params)
					if err != nil {
						return err
					}
					s.fList = append(s.fList, f)
				}

				// write && or ||
				s.fList = append(s.fList, fragment{
					Type:  model.LOGICAL,
					Value: tmpl[i-1 : i+1],
				})
			}
		default:
			_ = buf.WriteByte(tmpl[i])
		}
	}

	if sqlClause := buf.Dump(); sqlClause != "" {
		f, err = s.checkFragment(sqlClause, params)
		if err != nil {
			return err
		}
		s.fList = append(s.fList, f)
	}
	return nil
}

// check validition of clause's value
func (s *section) checkTempleFragmentValid() error {
	if s.Type == model.FOR {
		switch {
		case !s.isForFormat():
			return fmt.Errorf("for range syntax error:%s ", s.String())
		case !s.fList[6].IsArray:
			return fmt.Errorf("for loop argument [%s] is not array ", s.fList[6].Value)
		case s.SQLSlice.hasSameName(s.ForRange.value):
			return fmt.Errorf("cannot use the same value name in different for loops")
		}
	}
	s.Value = s.String()
	return nil
}

func (s *section) SetForRangeKey(key string) {
	s.ForRange.key = key
	s.fList[1].Value = key
	s.Value = s.String()
}

// isExpressionValid  check express valid
func (s *section) isExpressionValid(expr []fragment) bool {
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
func (s *section) isForFormat() bool {
	switch {
	case len(s.fList) != 7:
		return false
	case s.fList[1].Type != model.KEY:
		return false
	case s.fList[2].Value != ",":
		return false
	case s.fList[3].Type != model.VALUE:
		return false
	case s.fList[4].Value != ":=":
		return false
	case s.fList[5].Type != model.RANGE:
		return false
	}
	return true
}
