package check

import (
	"fmt"
	"strings"

	"gorm.io/gen/internal/model"
)

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
func (f *ForRange) DataValue(dataName, clauseName string) string {
	return f.mapIndexName("@", dataName, clauseName)
}
func (f *ForRange) appendDataToParams(dataName, clauseName string) string {
	return fmt.Sprintf("params[%s]=%s%s", f.mapIndexName("", dataName, clauseName), f.value, f.suffix)
}

type section struct {
	Type      model.Status
	Value     string
	ForRange  ForRange
	SQLSlice  *Sections
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
	tmpl := strings.TrimSpace(s.Value)
	s.splitList = strings.FieldsFunc(tmpl, func(r rune) bool {
		return r == ':' || r == ' ' || r == '=' || r == ','
	})
}

func (s *section) checkTemple() error {
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
			return fmt.Errorf("for range syntax error:%s", s.Value)
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
