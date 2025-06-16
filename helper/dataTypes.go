package helper

import (
	model "gorm.io/gen/internal/model"
)

// Gets the code Type for the sqlDataType.
func GetDataType(sqlDataType string) string {
	return model.GetDataType(sqlDataType)
}

// Add or override values in the override map.
func SetDataType(dbType string, getTypeFunc func(detailType string) (finalType string)) {
	model.SetDataType(dbType, getTypeFunc)
}
