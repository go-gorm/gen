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

{{if gt .TableCount 1}}
// TableCountOf{{.ModelStructName}} {{.ModelStructName}}'s table Count
func TableCountOf{{.ModelStructName}}() int {
	return {{.TableCount}}
}

// TableNameOf{{.ModelStructName}} {{.ModelStructName}}'s actual table name
func TableNameOf{{.ModelStructName}}(shardKey int64) string {
	return TableName{{.ModelStructName}} + strconv.FormatInt(shardKey%{{.TableCount}}, 10)
}
{{end}}
`

// ModelMethod model struct DIY method
const ModelMethod = `

{{if .Doc -}}// {{.DocComment -}}{{end}}
func ({{.GetBaseStructTmpl}}){{.MethodName}}({{.GetParamInTmpl}})({{.GetResultParamInTmpl}}){{.Body}}
`
