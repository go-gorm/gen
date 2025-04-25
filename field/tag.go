package field

import (
	"sort"
	"strings"
)

// TODO implement unit tests for tags

const (
	TagKeyGorm = "gorm"
	TagKeyJson = "json"

	//gorm tag
	TagKeyGormColumn        = "column"
	TagKeyGormType          = "type"
	TagKeyGormPrimaryKey    = "primaryKey"
	TagKeyGormAutoIncrement = "autoIncrement"
	TagKeyGormNotNull       = "not null"
	TagKeyGormUniqueIndex   = "uniqueIndex"
	TagKeyGormIndex         = "index"
	TagKeyGormDefault       = "default"
	TagKeyGormComment       = "comment"
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

type TagBuilder interface {
	Build() string
}

type Tag map[string]string

func (tag Tag) Set(key, value string) Tag {
	tag[key] = value
	return tag
}

func (tag Tag) Remove(key string) Tag {
	delete(tag, key)
	return tag
}

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

type GormTag map[string][]string

func (tag GormTag) Append(key string, values ...string) GormTag {
	if _, ok := tag[key]; ok {
		tag[key] = append(tag[key], values...)
	} else {
		tag[key] = values
	}
	return tag
}

func (tag GormTag) Set(key string, values ...string) GormTag {
	tag[key] = values
	return tag
}

func (tag GormTag) Remove(key string) GormTag {
	delete(tag, key)
	return tag
}

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
