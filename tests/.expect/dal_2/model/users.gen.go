// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameUser = "users"

// User mapped from table <users>
type User struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"-"`
	Name      *string        `gorm:"column:name" json:"-"`
	Age       *string        `gorm:"column:age;index:idx_age,priority:1" json:"-"`
	Address   *string        `gorm:"column:address" json:"-"`
	Role      *string        `gorm:"column:role" json:"-"`
	CreatedAt *time.Time     `gorm:"column:created_at" json:"-"`
	UpdatedAt *time.Time     `gorm:"column:updated_at" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index:idx_users_deleted_at,priority:1" json:"-"`
	Remark    *string        `gorm:"column:remark" json:"-"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
