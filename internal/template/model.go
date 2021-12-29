package template

// Model used as a variable because it cannot load template file after packed, params still can pass file
const Model = NotEditMark + `
package {{.StructInfo.Package}}

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	{{range .ImportPkgPaths}}{{.}} ` + "\n" + `{{end}}
)

const TableName{{.StructName}} = "{{.TableName}}"

// {{.StructName}} mapped from table <{{.TableName}}>
type {{.StructName}} struct {
    {{range .Fields}}
	{{if .MultilineComment -}}
	/*
{{.ColumnComment}}
    */
	{{end -}}
    {{.Name}} {{.Type}} ` + "`{{if .OverwriteTag}}{{.OverwriteTag}}{{else}}gorm:\"{{.GORMTag}}\" json:\"{{.JSONTag}}\"{{.NewTag}}{{end}}` " +
	"{{if not .MultilineComment}}{{if .ColumnComment}}// {{.ColumnComment}}{{end}}{{end}}" +
	`{{end}}
}

// TableName {{.StructName}}'s table name
func (*{{.StructName}}) TableName() string {
    return TableName{{.StructName}}
}
`
