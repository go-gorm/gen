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

func isStringEnd(b byte) bool {
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

func isDBUndefined(db *gorm.DB) bool {
	return db == nil
}

func getPackageName(fullName string) string {
	return strings.Split(fullName, ".")[0]
}

func isDigit(str string) bool {
	for _, x := range str {
		if !unicode.IsDigit(x) {
			return false
		}
	}
	return true
}

func allowType(typ string) bool {
	switch typ {
	case "string", "bytes":
		return true
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return true
	case "float64", "float32":
		return true
	case "bool":
		return true
	case "time.Time":
		return true
	default:
		return false
	}
}

func stringHasMore(i int, str string) bool {
	return i < len(str)-1
}

func DelPointer(name string) string {
	return strings.TrimLeft(name, "*")
}

func GetSimpleName(s string) string {
	return string(strings.ToLower(s)[0])
}
