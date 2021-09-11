package template

const DIYMethod = `
/*
{{.Doc -}}
*/
func ({{.S}} {{.TargetStruct}}Do){{.MethodName}}({{.GetParamInTmpl}})({{.GetResultParamInTmpl}}){
	{{if .HasSqlData}}params := map[string]interface{}{ {{range $index,$data:=.SqlData}}
		"{{$data}}":{{$data}},{{end}}
	}

	{{end}}var generateSQL string{{range $line:=.SqlTmplList}}{{$line}}
	{{end}}

	{{if .HasNeedNewResult}}result ={{if .ResultData.IsMap}}make{{else}}new{{end}}({{if ne .ResultData.Package ""}}{{.ResultData.Package}}.{{end}}{{.ResultData.Type}}){{end}}
	{{.ExecuteResult}} = {{.S}}.UnderlyingDB().{{.GormOption}}(generateSQL{{if .HasSqlData}},params{{end}}){{if not .ResultData.IsNull}}.Find({{if .HasGotPoint}}&{{end}}{{.ResultData.Name}}){{end}}.Error
	return
}

`

const CRUDMethod = `
func ({{.S}} {{.NewStructName}}Do) Debug() *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Debug().(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) WithContext(ctx context.Context) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.WithContext(ctx).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) Clauses(conds ...clause.Expression) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Clauses(conds...).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) Not(conds ...gen.Condition) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Not(conds...).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) Or(conds ...gen.Condition) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Or(conds...).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) Select(conds ...field.Expr) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Select(conds...).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) Where(conds ...gen.Condition) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Where(conds...).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) Order(conds ...field.Expr) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Order(conds...).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) Distinct(cols ...field.Expr) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Distinct(cols...).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) Omit(cols ...field.Expr) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Omit(cols...).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) Join(table schema.Tabler, on ...field.Expr) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Join(table, on...).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) LeftJoin(table schema.Tabler, on ...field.Expr) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.LeftJoin(table, on...).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) RightJoin(table schema.Tabler, on ...field.Expr) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.RightJoin(table, on...).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) Group(col field.Expr) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Group(col).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) Having(conds ...gen.Condition) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Having(conds...).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) Limit(limit int) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Limit(limit).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) Offset(offset int) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Offset(offset).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) Scopes(funcs ...func(gen.Dao) gen.Dao) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Scopes(funcs...).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) Unscoped() *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Unscoped().(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) Create(values ...*{{.StructInfo.Package}}.{{.StructInfo.Type}}) error {
	if len(values) == 0 {
		return nil
	}
	return {{.S}}.DO.Create(values)
}

func ({{.S}} {{.NewStructName}}Do) CreateInBatches(values []*{{.StructInfo.Package}}.{{.StructInfo.Type}}, batchSize int) error {
	return {{.S}}.DO.CreateInBatches(values, batchSize)
}

func ({{.S}} {{.NewStructName}}Do) Save(values ...*{{.StructInfo.Package}}.{{.StructInfo.Type}}) error {
	if len(values) == 0 {
		return nil
	}
	return {{.S}}.DO.Save(values)
}

func ({{.S}} {{.NewStructName}}Do) First() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error) {
	if result, err := {{.S}}.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*{{.StructInfo.Package}}.{{.StructInfo.Type}}), nil
	}
}

func ({{.S}} {{.NewStructName}}Do) Take() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error) {
	if result, err := {{.S}}.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*{{.StructInfo.Package}}.{{.StructInfo.Type}}), nil
	}
}

func ({{.S}} {{.NewStructName}}Do) Last() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error) {
	if result, err := {{.S}}.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*{{.StructInfo.Package}}.{{.StructInfo.Type}}), nil
	}
}

func ({{.S}} {{.NewStructName}}Do) Find() ([]*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error) {
	result, err := {{.S}}.DO.Find()
	return result.([]*{{.StructInfo.Package}}.{{.StructInfo.Type}}), err
}

func ({{.S}} {{.NewStructName}}Do) FindInBatches(result []*{{.StructInfo.Package}}.{{.StructInfo.Type}}, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return {{.S}}.DO.FindInBatches(&result, batchSize, fc)
}

func ({{.S}} {{.NewStructName}}Do) FindByPage(offset int, limit int) (result []*{{.StructInfo.Package}}.{{.StructInfo.Type}}, count int64, err error) {
	count, err = {{.S}}.Count()
	if err != nil {
		return
	}

	result, err = {{.S}}.Offset(offset).Limit(limit).Find()
	return
}

func ({{.S}} {{.NewStructName}}Do) Model(result *{{.StructInfo.Package}}.{{.StructInfo.Type}}) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Model(result).(*gen.DO)
	return &{{.S}}
}

func ({{.S}} {{.NewStructName}}Do) Begin(opts ...*sql.TxOptions) *{{.NewStructName}}Do {
	{{.S}}.DO = *{{.S}}.DO.Begin(opts...).(*gen.DO)
	return &{{.S}}
}

`
