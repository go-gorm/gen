package model

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	// ErrInvalidSyntax is returned when the StructTag syntax is invalid.
	ErrInvalidSyntax = errors.New("invalid syntax for key value pair")

	// ErrInvalidKey is returned if a key is containing invalid characters or
	// is missing.
	ErrInvalidKey = errors.New("invalid key")

	// ErrInvalidValue is returned if a value is not qouted.
	ErrInvalidValue = errors.New("invalid value")

	// ErrInvalidSeparator is returned if comma is used as separator.
	ErrInvalidSeparator = errors.New("invalid separator, key value pairs should be separated by spaces")
)

// Tag is just a map of key value pairs.
type Tag map[string]string

func Merge(tags ...Tag) Tag {
	for _, t := range tags {
		for k, v := range t {
			tags[0][k] = v
		}
	}

	return tags[0]
}

func (t Tag) StructTag() reflect.StructTag {
	var s string
	for k, v := range t {
		s += fmt.Sprintf(`%s:"%s" `, k, v)
	}
	return reflect.StructTag(strings.TrimSpace(s))
}

func Parse(st reflect.StructTag) Tag {
	tag := Tag{}

	for st != "" {
		i := 0
		for i < len(st) && st[i] == ' ' {
			i++
		}

		st = st[i:]
		if st == "" {
			break
		}

		i = 0
		for i < len(st) && st[i] > ' ' && st[i] != ':' && st[i] != '"' && st[i] != 0x7f {
			if st[i] == ',' {
				return tag
			}
			i++
		}

		if i == 0 {
			return tag
		}

		if i+1 >= len(st) || st[i] != ':' {
			return tag
		}

		if st[i+1] != '"' {
			return tag
		}

		key := string(st[:i])
		st = st[i+1:]

		i = 1
		for i < len(st) && st[i] != '"' {
			if st[i] == '\\' {
				i++
			}
			i++
		}

		if i >= len(st) {
			return tag
		}

		qvalue := string(st[:i+1])
		st = st[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			return tag
		}

		tag[key] = value
	}

	return tag
}
