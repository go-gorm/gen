package model

import (
	"regexp"

	"gorm.io/gorm"
)

type SchemaNameOpt func(*gorm.DB) string

// get mysql db' name
var dbNameReg = regexp.MustCompile(`/\w+\??`)

var defaultMysqlSchemaNameOpt = SchemaNameOpt(func(db *gorm.DB) string {
	return db.Migrator().CurrentDatabase()
})

type MemberOpt interface{ Self() func(*Member) *Member }

type ModifyMemberOpt func(*Member) *Member

func (o ModifyMemberOpt) Self() func(*Member) *Member { return o }

type FilterMemberOpt ModifyMemberOpt

func (o FilterMemberOpt) Self() func(*Member) *Member { return o }

type CreateMemberOpt ModifyMemberOpt

func (o CreateMemberOpt) Self() func(*Member) *Member { return o }

func sortOpt(opts []MemberOpt) (modifyOpts []MemberOpt, filterOpts []MemberOpt, createOpts []MemberOpt) {
	for _, opt := range opts {
		switch opt.(type) {
		case ModifyMemberOpt:
			modifyOpts = append(modifyOpts, opt)
		case FilterMemberOpt:
			filterOpts = append(filterOpts, opt)
		case CreateMemberOpt:
			createOpts = append(createOpts, opt)
		}
	}
	return
}
