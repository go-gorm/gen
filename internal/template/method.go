package template

const DIYMethod = `

//{{.DocComment }}
func ({{.S}} {{.TargetStruct}}Do){{.MethodName}}({{.GetParamInTmpl}})({{.GetResultParamInTmpl}}){
	{{if .HasSqlData}}params := map[string]interface{}{ {{range $index,$data:=.SqlData}}
		"{{$data}}":{{$data}},{{end}}
	}

	{{end}}var generateSQL string{{range $line:=.SqlTmplList}}{{$line}}
	{{end}}

	{{if .HasNeedNewResult}}result ={{if .ResultData.IsMap}}make{{else}}new{{end}}({{if ne .ResultData.Package ""}}{{.ResultData.Package}}.{{end}}{{.ResultData.Type}}){{end}}
	{{.ExecuteResult}} = {{.S}}.UnderlyingDB().{{.GormOption}}(generateSQL{{if .HasSqlData}},params{{end}}){{if not .ResultData.IsNull}}.{{.GormRunMethodName}}({{if .HasGotPoint}}&{{end}}{{.ResultData.Name}}){{end}}.Error
	return
}

`

const CRUDMethod = `
func ({{.S}} {{.NewStructName}}Do) Debug() *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Debug())
}

func ({{.S}} {{.NewStructName}}Do) WithContext(ctx context.Context) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.WithContext(ctx))
}

func ({{.S}} {{.NewStructName}}Do) Clauses(conds ...clause.Expression) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Clauses(conds...))
}

func ({{.S}} {{.NewStructName}}Do) Not(conds ...gen.Condition) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Not(conds...))
}

func ({{.S}} {{.NewStructName}}Do) Or(conds ...gen.Condition) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Or(conds...))
}

func ({{.S}} {{.NewStructName}}Do) Select(conds ...field.Expr) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Select(conds...))
}

func ({{.S}} {{.NewStructName}}Do) Where(conds ...gen.Condition) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Where(conds...))
}

func ({{.S}} {{.NewStructName}}Do) Order(conds ...field.Expr) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Order(conds...))
}

func ({{.S}} {{.NewStructName}}Do) Distinct(cols ...field.Expr) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Distinct(cols...))
}

func ({{.S}} {{.NewStructName}}Do) Omit(cols ...field.Expr) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Omit(cols...))
}

func ({{.S}} {{.NewStructName}}Do) Join(table schema.Tabler, on ...field.Expr) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Join(table, on...))
}

func ({{.S}} {{.NewStructName}}Do) LeftJoin(table schema.Tabler, on ...field.Expr) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.LeftJoin(table, on...))
}

func ({{.S}} {{.NewStructName}}Do) RightJoin(table schema.Tabler, on ...field.Expr) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.RightJoin(table, on...))
}

func ({{.S}} {{.NewStructName}}Do) Group(col field.Expr) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Group(col))
}

func ({{.S}} {{.NewStructName}}Do) Having(conds ...gen.Condition) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Having(conds...))
}

func ({{.S}} {{.NewStructName}}Do) Limit(limit int) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Limit(limit))
}

func ({{.S}} {{.NewStructName}}Do) Offset(offset int) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Offset(offset))
}

func ({{.S}} {{.NewStructName}}Do) Scopes(funcs ...func(gen.Dao) gen.Dao) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Scopes(funcs...))
}

func ({{.S}} {{.NewStructName}}Do) Unscoped() *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Unscoped())
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

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
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

func ({{.S}} {{.NewStructName}}Do) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) ([]*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error) {
	result, err := {{.S}}.DO.FindInBatch(batchSize, fc)
	return result.([]*{{.StructInfo.Package}}.{{.StructInfo.Type}}), err
}

func ({{.S}} {{.NewStructName}}Do) FindInBatches(result []*{{.StructInfo.Package}}.{{.StructInfo.Type}}, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return {{.S}}.DO.FindInBatches(&result, batchSize, fc)
}

func ({{.S}} {{.NewStructName}}Do) Attrs(attrs ...field.Expr) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Attrs(attrs...))
}

func ({{.S}} {{.NewStructName}}Do) Assign(attrs ...field.Expr) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Assign(attrs...))
}

func ({{.S}} {{.NewStructName}}Do) Joins(field field.RelationField) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Joins(field))
}

func ({{.S}} {{.NewStructName}}Do) Preload(field field.RelationField) *{{.NewStructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Preload(field))
}

func ({{.S}} {{.NewStructName}}Do) FirstOrInit() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error) {
	if result, err := {{.S}}.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*{{.StructInfo.Package}}.{{.StructInfo.Type}}), nil
	}
}

func ({{.S}} {{.NewStructName}}Do) FirstOrCreate() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error) {
	if result, err := {{.S}}.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*{{.StructInfo.Package}}.{{.StructInfo.Type}}), nil
	}
}

func ({{.S}} {{.NewStructName}}Do) FindByPage(offset int, limit int) (result []*{{.StructInfo.Package}}.{{.StructInfo.Type}}, count int64, err error) {
	count, err = {{.S}}.Count()
	if err != nil {
		return
	}

	result, err = {{.S}}.Offset(offset).Limit(limit).Find()
	return
}

func ({{.S}} *{{.NewStructName}}Do) withDO(do gen.Dao) (*{{.NewStructName}}Do) {
	{{.S}}.DO = *do.(*gen.DO)
	return {{.S}}
}

`
