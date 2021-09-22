package template

// ModelTemplate used as a variable because it cannot load template file after packed, params still can pass file
const ModelTemplate = NotEditMark + `
package {{.StructInfo.Package}}

import "time"

const TableName{{.StructName}} = "{{.TableName}}"

// {{.TableName}}
type {{.StructName}} struct {
    {{range .Members}}
	{{if .MultilineComment -}}
	/*
{{.ColumnComment}}
    */
	{{end -}}
    {{.Name}} {{.ModelType}} ` + "`json:\"{{.JSONTag}}\" gorm:\"{{.GORMTag}}\"{{.NewTag}}` " +
	"{{if not .MultilineComment}}{{if .ColumnComment}}// {{.ColumnComment}}{{end}}{{end}}" +
	`{{end}}
}

// TableName {{.StructName}}'s table name
func (*{{.StructName}}) TableName() string {
    return TableName{{.StructName}}
}
`
