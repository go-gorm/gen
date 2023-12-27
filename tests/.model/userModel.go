package _model

import (
	"time"
)

const TableNameUser = "users"

// User mapped from table <users>
type User struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement:true;comment:id" json:"-"`                                         // id
	CreatedAt    *time.Time `gorm:"column:created_at;comment:创建时间" json:"-"`                                                             // 创建时间
	Name         *string    `gorm:"column:name;index:idx_name,priority:1;index:idx_name_company_id,priority:1;comment:oneline" json:"-"` // oneline
	Address      *string    `gorm:"column:address;comment:地址" json:"-"`                                                                  // 地址
	RegisterTime *time.Time `gorm:"column:register_time;comment:注册时间" json:"-"`                                                          // 注册时间
	/*
		multiline
		line1
		line2
	*/
	Alive      *bool   `gorm:"column:alive;comment:multiline\nline1\nline2" json:"-"`
	CompanyID  *int64  `gorm:"column:company_id;index:idx_name_company_id,priority:2;default:666;comment:公司id" json:"-"` // 公司id
	PrivateURL *string `gorm:"column:private_url;default:https://a.b.c ;comment:私人地址" json:"-"`                          // 私人地址
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
