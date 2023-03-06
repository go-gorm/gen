package field

import "strings"

const TagKeyGorm = "gorm"
const TagKeyJson = "json"

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
	gormTag := tag[TagKeyGorm]
	delete(tag, TagKeyGorm)

	tags := make([]string, 0, len(tag))
	for k, v := range tag {
		if k == "" || v == "" {
			continue
		}
		tags = append(tags, k+":\""+v+"\"")
	}
	if gormTag != "" { //first tag gorm
		tags = append([]string{TagKeyGorm + ":\"" + gormTag + "\""}, tags...)
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
	for k, v := range tag {
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
