package helper

import (
	model "gorm.io/gen/internal/model"
)

// GetDataType returns the corresponding Go data type for a given SQL data type string.
// It delegates the type mapping logic to the model.GetDataType function.
//
// Parameters:
//   - sqlDataType: the SQL data type as a string.
//
// Returns:
//   - The Go data type as a string.
func GetDataType(sqlDataType string) string {
	return model.GetDataType(sqlDataType)
}

// SetDataType registers a mapping between a database type and a function that determines the final type.
// The getTypeFunc parameter is a function that takes a detailType string and returns the corresponding finalType string.
// This function delegates the registration to the model.SetDataType function.
func SetDataType(dbType string, getTypeFunc func(detailType string) (finalType string)) {
	model.SetDataType(dbType, getTypeFunc)
}
