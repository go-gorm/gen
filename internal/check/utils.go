package check

import (
	"strings"
	"unicode"

	"gorm.io/gorm"
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

func isDBUnset(db *gorm.DB) bool {
	return db == nil
}

func getPackageName(fullName string) string {
	return strings.Split(delPointerSym(fullName), ".")[0]
}

func isDigit(str string) bool {
	for _, x := range str {
		if !unicode.IsDigit(x) {
			return false
		}
	}
	return true
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

func getNewTypeName(t string) string {
	list := strings.Split(t, ".")
	return strings.Title(list[len(list)-1])
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

func contains(target string, list []string) bool {
	for _, item := range list {
		if target == item {
			return true
		}
	}
	return false
}

func containMultiline(s string) bool {
	return strings.Contains(s, "\n")
}
