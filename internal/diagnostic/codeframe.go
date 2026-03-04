package diagnostic

import (
	"bytes"
	"fmt"
	"os"
)

func CodeFrame(src []byte, line, column, context int) string {
	if context < 0 {
		context = 0
	}
	if line < 1 {
		line = 1
	}
	if column < 1 {
		column = 1
	}

	lines := bytes.Split(src, []byte("\n"))
	if len(lines) == 0 {
		return ""
	}
	for len(lines) > 1 && len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	if line > len(lines) {
		line = len(lines)
	}

	start := line - 1 - context
	if start < 0 {
		start = 0
	}
	end := line - 1 + context
	if end >= len(lines) {
		end = len(lines) - 1
	}

	width := digits(end + 1)
	var out bytes.Buffer
	for i := start; i <= end; i++ {
		ln := i + 1
		prefix := fmt.Sprintf("%*d | ", width, ln)
		if i == line-1 {
			prefix = fmt.Sprintf(">%*d | ", width, ln)
		}
		rawLine := string(lines[i])
		expandedLine := expandTabs(rawLine, 4)
		out.WriteString(prefix)
		out.WriteString(expandedLine)
		out.WriteByte('\n')

		if i == line-1 {
			caretPos := caretColumn(rawLine, column, 4)
			out.WriteString(fmt.Sprintf("%*s | ", width+1, ""))
			out.WriteString(spaces(caretPos))
			out.WriteByte('^')
			out.WriteByte('\n')
		}
	}
	return out.String()
}

func CodeFrameFromFile(file string, line, column, context int) (string, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return CodeFrame(b, line, column, context), nil
}

func digits(n int) int {
	if n <= 0 {
		return 1
	}
	d := 0
	for n > 0 {
		n /= 10
		d++
	}
	return d
}

func spaces(n int) string {
	if n <= 0 {
		return ""
	}
	return string(bytes.Repeat([]byte(" "), n))
}

func expandTabs(s string, tabWidth int) string {
	if tabWidth <= 0 {
		tabWidth = 4
	}
	var out bytes.Buffer
	col := 0
	for _, r := range s {
		if r == '\t' {
			spaceCount := tabWidth - (col % tabWidth)
			out.Write(bytes.Repeat([]byte(" "), spaceCount))
			col += spaceCount
			continue
		}
		out.WriteRune(r)
		col++
	}
	return out.String()
}

func caretColumn(raw string, column, tabWidth int) int {
	if column <= 1 {
		return 0
	}
	if tabWidth <= 0 {
		tabWidth = 4
	}
	b := []byte(raw)
	bytePos := column - 1
	if bytePos > len(b) {
		bytePos = len(b)
	}
	prefix := string(b[:bytePos])
	col := 0
	for _, r := range prefix {
		if r == '\t' {
			col += tabWidth - (col % tabWidth)
			continue
		}
		col++
	}
	return col
}
