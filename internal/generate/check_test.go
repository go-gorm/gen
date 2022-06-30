package generate

import "strings"

//struct test custom method,if this in struct_test.go file, Will be ignored running the test

// OnlyForTestUser ...
type OnlyForTestUser struct {
	ID   int32
	Name string
}

// IsEmpty is a custom method
func (u *OnlyForTestUser) IsEmpty() bool {
	if u == nil {
		return true
	}

	return u.ID == 0
}

// SetName set user name
func (u *OnlyForTestUser) SetName(name string) {
	u.Name = name
}

// GetName get to lower name
func (u *OnlyForTestUser) GetName() string {
	return strings.ToLower(u.Name)
}
