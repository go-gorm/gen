package field

import (
	"sort"
	"strings"
)

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

func NewTag() Tag {
	return Tag{}
}

func (tag Tag) Set(key, value string) {
	tag[key] = value
}

func (tag Tag) Remove(key string) {
	delete(tag, key)
}

func (tag Tag) Build() string {
	if tag == nil || len(tag) == 0 {
		return ""
	}

	tags := make([]string, 0, len(tag))
	keys := tagKeySort(tag)
	for _, k := range keys {
		v := tag[k]
		if k == "" || v == "" {
			continue
		}
		tags = append(tags, k+":\""+v+"\"")
	}
	return strings.Join(tags, " ")
}

type GormTag Tag

func NewGormTag() GormTag {
	return GormTag{}
}

func (tag GormTag) Set(key, value string) {
	tag[key] = value
}

func (tag GormTag) Remove(key string) {
	delete(tag, key)
}

func (tag GormTag) Build() string {
	if tag == nil || len(tag) == 0 {
		return ""
	}
	tags := make([]string, 0, len(tag))
	keys := tagKeySort(Tag(tag))
	for _, k := range keys {
		v := tag[k]
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

	return strings.Join(tags, ";")
}

func tagKeySort(tag Tag) []string {
	keys := make([]string, 0, len(tag))
	if len(tag) == 0 {
		return keys
	}
	for k, _ := range tag {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if tagKeyPriorities[keys[i]] == tagKeyPriorities[keys[j]] {
			return keys[i] <= keys[j]
		}
		return tagKeyPriorities[keys[i]] > tagKeyPriorities[keys[j]]
	})
	return keys
}
