package field

import (
	"sort"
	"strings"
)

// TODO implement unit tests for tags

const (
	// TagKeyGorm is the struct-tag key containing GORM directives.
	TagKeyGorm = "gorm"
	// TagKeyJson is the struct-tag key containing a JSON field name.
	// The historical spelling is retained for API compatibility.
	TagKeyJson = "json"

	// TagKeyGormColumn names the GORM column directive.
	TagKeyGormColumn = "column"
	// TagKeyGormType names the GORM data-type directive.
	TagKeyGormType = "type"
	// TagKeyGormPrimaryKey names the GORM primary-key directive.
	TagKeyGormPrimaryKey = "primaryKey"
	// TagKeyGormAutoIncrement names the GORM auto-increment directive.
	TagKeyGormAutoIncrement = "autoIncrement"
	// TagKeyGormNotNull names the GORM not-null directive.
	TagKeyGormNotNull = "not null"
	// TagKeyGormUniqueIndex names the GORM unique-index directive.
	TagKeyGormUniqueIndex = "uniqueIndex"
	// TagKeyGormIndex names the GORM index directive.
	TagKeyGormIndex = "index"
	// TagKeyGormDefault names the GORM default-value directive.
	TagKeyGormDefault = "default"
	// TagKeyGormComment names the GORM column-comment directive.
	TagKeyGormComment = "comment"
)

var (
	tagKeyPriorities = map[string]int16{
		TagKeyGorm: 100,
		TagKeyJson: 99,

		TagKeyGormColumn:        10,
		TagKeyGormType:          9,
		TagKeyGormPrimaryKey:    8,
		TagKeyGormAutoIncrement: 7,
		TagKeyGormNotNull:       6,
		TagKeyGormUniqueIndex:   5,
		TagKeyGormIndex:         4,
		TagKeyGormDefault:       3,
		TagKeyGormComment:       0,
	}
)

// TagBuilder serializes a collection of struct-tag components.
type TagBuilder interface {
	// Build returns the tag text without surrounding backticks.
	Build() string
}

// Tag maps struct-tag keys to their quoted values.
type Tag map[string]string

// Set stores value for key and returns tag for fluent updates.
func (tag Tag) Set(key, value string) Tag {
	tag[key] = value
	return tag
}

// Remove deletes key and returns tag for fluent updates.
func (tag Tag) Remove(key string) Tag {
	delete(tag, key)
	return tag
}

// Build renders struct tags in deterministic priority order.
func (tag Tag) Build() string {
	if len(tag) == 0 {
		return ""
	}

	tags := make([]string, 0, len(tag))
	for _, k := range tag.keys() {
		v := tag[k]
		if k == "" {
			continue
		}
		tags = append(tags, k+":\""+v+"\"")
	}
	return strings.Join(tags, " ")
}

func (tag Tag) keys() []string {
	if len(tag) == 0 {
		return nil
	}

	keys := make([]string, 0, len(tag))
	for k := range tag {
		keys = append(keys, k)
	}
	return keySort(keys)
}

// GormTag maps GORM directive names to their optional values.
type GormTag map[string][]string

// Append adds values to key without replacing existing values.
func (tag GormTag) Append(key string, values ...string) GormTag {
	if _, ok := tag[key]; ok {
		tag[key] = append(tag[key], values...)
	} else {
		tag[key] = values
	}
	return tag
}

// Set replaces the values associated with key.
func (tag GormTag) Set(key string, values ...string) GormTag {
	tag[key] = values
	return tag
}

// Remove deletes key and returns tag for fluent updates.
func (tag GormTag) Remove(key string) GormTag {
	delete(tag, key)
	return tag
}

// Build renders directives as a deterministic semicolon-separated value.
func (tag GormTag) Build() string {
	if len(tag) == 0 {
		return ""
	}

	tags := make([]string, 0, len(tag))
	for _, k := range tag.keys() {
		vs := tag[k]
		if len(vs) == 0 && k == "" {
			continue
		}
		if len(vs) == 0 {
			tags = append(tags, k)
			continue
		}
		for _, v := range vs {
			if k == "" && v == "" {
				continue
			}
			tv := make([]string, 0, 2)
			if k != "" {
				tv = append(tv, k)
			}
			if v != "" {
				tv = append(tv, v)
			}
			tags = append(tags, strings.Join(tv, ":"))
		}
	}

	return strings.Join(tags, ";")
}

func (tag GormTag) keys() []string {
	if len(tag) == 0 {
		return nil
	}

	keys := make([]string, 0, len(tag))
	for k := range tag {
		keys = append(keys, k)
	}
	return keySort(keys)
}

func keySort(keys []string) []string {
	if len(keys) == 0 {
		return keys
	}

	sort.Slice(keys, func(i, j int) bool {
		if tagKeyPriorities[keys[i]] == tagKeyPriorities[keys[j]] {
			return keys[i] <= keys[j]
		}
		return tagKeyPriorities[keys[i]] > tagKeyPriorities[keys[j]]
	})
	return keys
}
