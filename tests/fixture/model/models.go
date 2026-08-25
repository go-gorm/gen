package model

// Company is the runtime fixture's parent model.
type Company struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Users []User
}

// User is the primary runtime fixture model.
type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Age       int
	Active    bool
	Role      string
	CompanyID uint
	Company   Company
	Orders    []Order
}

// Order is the runtime fixture's child model.
type Order struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint
	Amount int64
	Status string
	User   User
}
