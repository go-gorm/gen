package model

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// EnumValue represents a single value in an enum type
type EnumValue struct {
	Name  string
	Value string
}

// EnumType represents a PostgreSQL enum type with its values
type EnumType struct {
	Schema string
	Name   string
	Values []EnumValue
}

// GetEnumTypes retrieves all enum types and their values from PostgreSQL
func GetEnumTypes(db *gorm.DB) ([]EnumType, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	// Only proceed for PostgreSQL
	if db.Dialector.Name() != "postgres" {
		return nil, nil
	}

	// Query to get all enum types and their values
	var results []struct {
		EnumSchema string `gorm:"column:enum_schema"`
		EnumName   string `gorm:"column:enum_name"`
		EnumValue  string `gorm:"column:enum_value"`
		SortOrder  int    `gorm:"column:enum_sortorder"`
	}

	err := db.Raw(`
		SELECT n.nspname AS enum_schema, 
		       t.typname AS enum_name, 
		       e.enumlabel AS enum_value,
		       e.enumsortorder AS enum_sortorder
		FROM pg_type t 
		JOIN pg_enum e ON t.oid = e.enumtypid 
		JOIN pg_catalog.pg_namespace n ON n.oid = t.typnamespace 
		ORDER BY enum_schema, enum_name, e.enumsortorder
	`).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Group results by enum type
	enumTypes := make(map[string]EnumType)
	for _, r := range results {
		key := fmt.Sprintf("%s.%s", r.EnumSchema, r.EnumName)
		enum, exists := enumTypes[key]
		if !exists {
			enum = EnumType{
				Schema: r.EnumSchema,
				Name:   r.EnumName,
				Values: []EnumValue{},
			}
		}

		// Create a Go-friendly name for the enum value
		valueName := formatEnumValueName(r.EnumName, r.EnumValue)
		enum.Values = append(enum.Values, EnumValue{Name: valueName, Value: r.EnumValue})
		enumTypes[key] = enum
	}

	// Convert map to slice
	result := make([]EnumType, 0, len(enumTypes))
	for _, enum := range enumTypes {
		result = append(result, enum)
	}

	return result, nil
}

// GetColumnEnumType gets the enum type for a specific column if it exists
func GetColumnEnumType(db *gorm.DB, tableName, columnName string, enumTypes []EnumType) *EnumType {
	if db == nil || db.Dialector.Name() != "postgres" {
		return nil
	}

	// Query to get the enum type for this column
	var result struct {
		UdtSchema string `gorm:"column:udt_schema"`
		UdtName   string `gorm:"column:udt_name"`
	}

	err := db.Raw(`
		SELECT udt_schema, udt_name 
		FROM information_schema.columns 
		WHERE table_name = ? AND column_name = ? AND data_type = 'USER-DEFINED'
	`, tableName, columnName).Scan(&result).Error

	if err != nil || result.UdtName == "" {
		return nil
	}

	// Find matching enum type
	for _, enum := range enumTypes {
		if enum.Schema == result.UdtSchema && enum.Name == result.UdtName {
			return &enum
		}
	}

	return nil
}

// formatEnumValueName converts enum values to Go constant names
// Example: 'active_status' -> 'ActiveStatus'
func formatEnumValueName(enumName, enumValue string) string {
	// Remove any special characters and replace with underscores
	clean := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, enumValue)

	// Split by underscores and capitalize each part
	parts := strings.Split(clean, "_")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}

	return strings.Join(parts, "")
}