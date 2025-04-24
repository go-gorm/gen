package main

import (
	"context"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gorm.io/gen"
	"gorm.io/gen/examples/dal/model" // Adjust path if needed
	"gorm.io/gen/field"
)

// --- Define a custom DAO for the User model --- //

// UserDo embeds GenericDo[model.User, *UserDo] and adds custom methods
type UserDo struct {
	gen.GenericDo[model.User, *UserDo] // Embed the new GenericDo with Self type
	// You can add specific fields for UserDo if needed
}

// --- Remove redundant method overrides --- 
// Where, WithContext, Order etc. are now inherited correctly

// // Where 重载，返回 *UserDo
// defaultWhere := func(u *UserDo, conds ...gen.Condition) *UserDo {
// 	u.GenericDo = *u.GenericDo.Where(conds...)
// 	return u
// }
// func (u *UserDo) Where(conds ...gen.Condition) *UserDo {
// 	return defaultWhere(u, conds...)
// }
// 
// // WithContext 重载，返回 *UserDo
// func (u *UserDo) WithContext(ctx context.Context) *UserDo {
// 	u.GenericDo = *u.GenericDo.WithContext(ctx)
// 	return u
// }
// 
// // Order 重载，返回 *UserDo
// func (u *UserDo) Order(orders ...field.Expr) *UserDo {
// 	u.GenericDo = *u.GenericDo.Order(orders...)
// 	return u
// }

// NewUserDo creates a new UserDo instance
func NewUserDo(db *gorm.DB, opts ...gen.DOOption) *UserDo {
	_do := &UserDo{}

	// Define the clone function for this specific DAO type
	cloneFunc := func(db *gorm.DB) *UserDo {
		newDo := &UserDo{}
		newDo.GenericDo.DO = _do.GenericDo.DO.getInstance(db)
		newDo.GenericDo.cloneFunc = cloneFunc // Important: Set the cloneFunc in the new instance
		newDo.GenericDo.modelType = _do.GenericDo.modelType
		return newDo
	}

	// Initialize the embedded GenericDo, passing the clone function
	_do.GenericDo = gen.NewGenericDo[model.User, *UserDo](db, cloneFunc, opts...)

	return _do
}

// Custom method: FindByUserName finds a user by their name
func (u *UserDo) FindByUserName(name string) (*model.User, error) {
	// Use the embedded GenericDo's methods (now correctly typed)
	// Assuming 'users' table and 'name' column
	// The Where method now returns *UserDo, allowing chaining if needed before First()
	return u.Where(field.NewString("users", "name").Eq(name)).First()
}

// --- Example Usage --- //

func main() {
	// Initialize DB (same as generic_usage_example.go)
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Ensure the table name matches what's used in FindByUserName if not default
	db.Table("users").AutoMigrate(&model.User{}) // Migrate User model
	if err != nil {
		fmt.Println("AutoMigrate failed:", err)
	}

	fmt.Println("--- Extending GenericDo Example (Improved Chaining) --- ")

	// 1. Create an instance of the custom UserDo
	userDAO := NewUserDo(db)

	// 2. Use a method from the embedded GenericDo (chaining now returns *UserDo)
	newUser := &model.User{Name: "Extended User", Age: 35}
	err = userDAO.Create(newUser) // Calls GenericDo.Create
	if err != nil {
		fmt.Println("Failed to create user:", err)
	} else {
		fmt.Printf("Created user: %+v\n", newUser)
	}

	// 3. Use the custom method FindByUserName
	foundUser, err := userDAO.FindByUserName("Extended User")
	if err != nil {
		fmt.Println("Failed to find user by name:", err)
	} else {
		fmt.Printf("Found user using custom method: %+v\n", foundUser)
	}

	// 4. Example of chaining inherited methods (Where returns *UserDo)
	foundUserChained, err := userDAO.Where(field.NewInt("users", "age").Gt(30)).
		Where(field.NewString("users", "name").Eq("Extended User")).
		First()
	if err != nil {
		fmt.Println("Failed to find user with chaining:", err)
	} else {
		fmt.Printf("Found user using chained inherited methods: %+v\n", foundUserChained)
	}

	// 5. Example of finding a non-existent user
	_, err = userDAO.FindByUserName("NonExistent User")
	if err != nil {
		fmt.Printf("Correctly failed to find non-existent user: %v\n", err)
	} else {
		fmt.Println("Error: Found a user that shouldn't exist!")
	}

	fmt.Println("Extending GenericDo example finished.")
}

/*
Note:
- This example demonstrates embedding the improved GenericDo[T, Self].
- Custom methods (like FindByUserName) are added to UserDo.
- Inherited chainable methods (Where, Order, etc.) now correctly return *UserDo.
- No need to override chainable methods in UserDo anymore.
- Assumes model.User has Name and Age fields, and the table is 'users'. Adjust as needed.
*/