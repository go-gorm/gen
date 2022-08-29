// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameLanguage = "languages"

// Language mapped from table <languages>
type Language struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"-"`
	CreatedAt *time.Time     `gorm:"column:created_at" json:"-"`
	UpdatedAt *time.Time     `gorm:"column:updated_at" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index:idx_languages_deleted_at,priority:1" json:"-"`
	Name      *string        `gorm:"column:name" json:"-"`
}

// TableName Language's table name
func (*Language) TableName() string {
	return TableNameLanguage
}