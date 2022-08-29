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

{{if .TableName -}}const TableName{{.ModelStructName}} = "{{.TableName}}"{{- end}}

// {{.ModelStructName}} {{.StructComment}}
type {{.ModelStructName}} struct {
    {{range .Fields}}
	{{if .MultilineComment -}}
	/*
{{.ColumnComment}}
    */
	{{end -}}
    {{.Name}} {{.Type}} ` + "`{{.Tags}}` " +
	"{{if not .MultilineComment}}{{if .ColumnComment}}// {{.ColumnComment}}{{end}}{{end}}" +
	`{{end}}
}

{{if .TableName -}}
// TableName {{.ModelStructName}}'s table name
func (*{{.ModelStructName}}) TableName() string {
    return TableName{{.ModelStructName}}
}
{{- end}}
`

// ModelMethod model struct DIY method
const ModelMethod = `

{{if .Doc -}}// {{.DocComment -}}{{end}}
func ({{.GetBaseStructTmpl}}){{.MethodName}}({{.GetParamInTmpl}})({{.GetResultParamInTmpl}}){{.Body}}
`
