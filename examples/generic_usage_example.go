package main

import (
	"context"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gorm.io/gen"
	"gorm.io/gen/examples/dal/model" // Assuming models are here, adjust if needed
	"gorm.io/gen/field"
)

// Example usage of GenericDo[T]

func main() {
	// Initialize a dummy GORM DB connection (replace with your actual DB setup)
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// AutoMigrate the models (optional, for demonstration)
	err = db.AutoMigrate(&model.User{}, &model.Person{})
	if err != nil {
		fmt.Println("AutoMigrate failed:", err)
		// Continue even if migration fails for demo purposes
	}

	// --- Using GenericDo for User model ---
	fmt.Println("--- User Operations ---")
	// 1. Create a GenericDo instance for the User model
	userDAO := gen.NewGenericDo[model.User](db)

	// 2. Example: Create a new user
	newUser := &model.User{Name: "Generic User", Age: 30} // Assuming User has Name and Age
	err = userDAO.Create(newUser)
	if err != nil {
		fmt.Println("Failed to create user:", err)
	} else {
		fmt.Printf("Created user: %+v\n", newUser)
	}

	// 3. Example: Find a user by ID (assuming ID field exists and is primary key)
	// We need the actual field definitions. Let's assume User has an 'ID' field.
	// To access fields like in generated code, you'd typically define them separately
	// or access them via reflection/struct tags if GenericDo provided helpers (it doesn't directly).
	// For simple queries, you can use strings or gorm clauses.
	foundUser, err := userDAO.Where(field.NewInt64("users", "id").Eq(newUser.ID)).First()
	if err != nil {
		fmt.Println("Failed to find user:", err)
	} else {
		fmt.Printf("Found user: %+v\n", foundUser)
	}

	// 4. Example: Find users with age > 25
	usersOver25, err := userDAO.Where(field.NewInt32("users", "age").Gt(25)).Find()
	if err != nil {
		fmt.Println("Failed to find users over 25:", err)
	} else {
		fmt.Printf("Users over 25: %d found\n", len(usersOver25))
		// for _, u := range usersOver25 {
		// 	fmt.Printf("  - %+v\n", u)
		// }
	}

	// --- Using GenericDo for Person model ---
	fmt.Println("\n--- Person Operations ---")
	// 1. Create a GenericDo instance for the Person model
	personDAO := gen.NewGenericDo[model.Person](db)

	// 2. Example: Create a new person
	newPerson := &model.Person{Name: "Generic Person", Age: 40} // Assuming Person has Name and Age
	err = personDAO.Create(newPerson)
	if err != nil {
		fmt.Println("Failed to create person:", err)
	} else {
		fmt.Printf("Created person: %+v\n", newPerson)
	}

	// 3. Example: Find the created person
	foundPerson, err := personDAO.Where(field.NewInt64("people", "id").Eq(newPerson.ID)).First()
	if err != nil {
		fmt.Println("Failed to find person:", err)
	} else {
		fmt.Printf("Found person: %+v\n", foundPerson)
	}

	// --- Chaining example ---
	fmt.Println("\n--- Chaining Example --- ")
	ctx := context.Background()
	userByName, err := userDAO.WithContext(ctx).Where(field.NewString("users", "name").Eq("Generic User")).Order(field.NewInt64("users", "id").Desc()).First()
	if err != nil {
		fmt.Println("Failed to find user by name with chaining:", err)
	} else {
		fmt.Printf("Found user via chaining: %+v\n", userByName)
	}

	fmt.Println("\nGenericDo example finished.")
}

/*
Note:
- Replace `gorm.io/gen/examples/dal/model` with the actual path to your models.
- This example assumes `User` has `ID`, `Name`, `Age` fields and `Person` has `ID`, `Name`, `Age` fields.
  Adjust field names (`field.NewXXX("table_name", "column_name")`) according to your actual model definitions and table names.
- For complex queries or type safety similar to generated code, you might still want to define field constants or helpers separately,
  as `GenericDo` itself doesn't automatically provide the `User.Name` style field access.
*/