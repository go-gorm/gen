package model

import "gorm.io/gorm"

// some struct implement manually

// Customer a struct mapping to table customers
type Customer struct {
	gorm.Model

	Name    string  `gorm:"type:varchar(100);not null"`
	Age     int     `gorm:"type:int"`
	Phone   string  `gorm:"type:varchar(11)"`
	Address string  `gorm:"type:text"`
	Amount  float64 `gorm:"type:float"`
}

func (Customer) TableName() string {
	return "customers"
}
