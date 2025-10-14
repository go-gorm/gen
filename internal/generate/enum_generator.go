package generate

import (
	"fmt"
	"strings"
	"text/template"

	"gorm.io/gen/internal/model"
	"gorm.io/gorm"
)

// EnumGenerator generates enum constants for PostgreSQL enum types
type EnumGenerator struct {
	db        *gorm.DB
	enumTypes []model.EnumType
}

// NewEnumGenerator creates a new EnumGenerator
func NewEnumGenerator(db *gorm.DB) (*EnumGenerator, error) {
	enumTypes, err := model.GetEnumTypes(db)
	if err != nil {
		return nil, err
	}

	return &EnumGenerator{
		db:        db,
		enumTypes: enumTypes,
	}, nil
}

// GenerateEnumCode generates Go code for enum constants
func (g *EnumGenerator) GenerateEnumCode(tableName string, columns []*model.Column) string {
	if g.db.Dialector.Name() != "postgres" || len(g.enumTypes) == 0 {
		return ""
	}

	var enumsForTable []model.EnumType
	for _, column := range columns {
		enumType := model.GetColumnEnumType(g.db, tableName, column.Name(), g.enumTypes)
		if enumType != nil {
			enumsForTable = append(enumsForTable, *enumType)
		}
	}

	if len(enumsForTable) == 0 {
		return ""
	}

	return g.renderEnumCode(tableName, enumsForTable)
}

// renderEnumCode renders the enum constants code using a template
func (g *EnumGenerator) renderEnumCode(tableName string, enums []model.EnumType) string {
	const enumTemplate = `
// Enum values for {{.TableName}} table
const (
{{- range .Enums}}
	// {{.Name}} enum values
	{{- range .Values}}
	{{$.TableName}}{{.Name}} = "{{.Value}}"
	{{- end}}
{{- end}}
)
`

	tmpl, err := template.New("enum").Parse(enumTemplate)
	if err != nil {
		return ""
	}

	data := struct {
		TableName string
		Enums     []model.EnumType
	}{
		TableName: formatTableName(tableName),
		Enums:     enums,
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return ""
	}

	return buf.String()
}

// formatTableName formats a table name for use in constant names
// Example: "user_profiles" -> "UserProfile"
func formatTableName(tableName string) string {
	// Remove schema prefix if present
	if idx := strings.LastIndex(tableName, "."); idx >= 0 {
		tableName = tableName[idx+1:]
	}

	// Split by underscores and capitalize each part
	parts := strings.Split(tableName, "_")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}

	// Remove trailing 's' for singular form
	result := strings.Join(parts, "")
	if len(result) > 0 && result[len(result)-1] == 's' {
		result = result[:len(result)-1]
	}

	return result
}