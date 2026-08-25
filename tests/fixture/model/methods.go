package model

import "gorm.io/gen"

// UserMethods defines the executable, dialect-neutral DIY SQL fixture.
type UserMethods interface {
	// FindByNameBranch chooses a name filter when supplied and an active filter otherwise.
	//
	// SELECT * FROM @@table
	// {{where}}
	// {{if name != ""}}
	// name = @name
	// {{else}}
	// active = @active
	// {{end}}
	// {{end}}
	FindByNameBranch(name string, active bool) ([]*gen.T, error)

	// FindByNames builds an OR list. An empty list intentionally has no WHERE clause.
	//
	// SELECT * FROM @@table
	// {{where}}
	// {{for _, name := range names}}
	// OR name = @name
	// {{end}}
	// {{end}}
	FindByNames(names []string) ([]*gen.T, error)

	// FindByAttributes applies every supplied string attribute as an AND predicate.
	//
	// SELECT * FROM @@table
	// {{where}}
	// {{for column, value := range attributes}}
	// AND @@column = @value
	// {{end}}
	// {{end}}
	FindByAttributes(attributes map[string]string) ([]*gen.T, error)

	// UpdateOptional updates requested values and remains valid when both are omitted.
	//
	// UPDATE @@table
	// {{set}}
	// {{if name != ""}}
	// name = @name,
	// {{else}}
	// name = name,
	// {{end}}
	// {{if age > 0}}
	// age = @age,
	// {{end}}
	// {{end}}
	// WHERE id = @id
	UpdateOptional(id uint, name string, age int) (gen.RowsAffected, error)

	// FindWithTrim trims the trailing OR and falls back to an always-false predicate.
	//
	// SELECT * FROM @@table WHERE
	// {{trim}}
	// {{for _, name := range names}}
	// name = @name OR
	// {{end}}
	// id < 0
	// {{end}}
	FindWithTrim(names []string) ([]*gen.T, error)

	// InsertUser omits MySQL-specific upsert syntax so it can run on SQLite.
	//
	// INSERT INTO @@table (name, age, active, role, company_id)
	// VALUES (@name, @age, @active, @role, @companyID)
	InsertUser(name string, age int, active bool, role string, companyID uint) (gen.RowsAffected, error)
}
