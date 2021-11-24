package check

import (
	"strings"
)

func isCapitalize(s string) bool {
	if len(s) < 1 {
		return false
	}
	b := s[0]
	if b >= 'A' && b <= 'Z' {
		return true
	}
	return false
}

func isEnd(b byte) bool {
	switch {
	case b >= 'a' && b <= 'z':
		return false
	case b >= 'A' && b <= 'Z':
		return false
	case b >= '0' && b <= '9':
		return false
	case b == '-' || b == '_' || b == '.':
		return false
	default:
		return true
	}
}

func getPackageName(fullName string) string {
	return strings.Split(delPointerSym(fullName), ".")[0]
}

func isDigit(str string) string {
	for _, x := range str {
		switch x {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			continue
		case '+', '-', '*', '/', '%':
			continue
		default:
			return ""
		}
	}
	return str
}

func isRangeVar(str string) string {
	if strings.HasPrefix(str, "_") || strings.HasPrefix(str, "$") {
		return ""
	}
	return str
}

func strOutrange(index int, str string) bool {
	return index >= len(str)
}

func delPointerSym(name string) string {
	return strings.TrimLeft(name, "*")
}

func getPureName(s string) string {
	return string(strings.ToLower(delPointerSym(s))[0])
}

// not need capitalize
func getStructName(t string) string {
	list := strings.Split(t, ".")
	return list[len(list)-1]
}

func uncaptialize(s string) string {
	if s == "" {
		return ""
	}

	return strings.ToLower(s[:1]) + s[1:]
}
