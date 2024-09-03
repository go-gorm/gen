package helper

import (
	model "gorm.io/gen/internal/model"
)

// Gets the code Type for the sqlDataType.
func GetDataType(sqlDataType string) string {
	return model.GetDataType(sqlDataType)
}
