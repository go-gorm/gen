package template

const DIYMethod = `

//{{.DocComment }}
func ({{.S}} {{.TargetStruct}}Do){{.MethodName}}({{.GetParamInTmpl}})({{.GetResultParamInTmpl}}){
	{{if .HasSqlData}}params :=make(map[string]interface{},0)

	{{end}}var generateSQL strings.Builder
	{{range $line:=.Sections.Tmpl}}{{$line}}
	{{end}}

	{{if .HasNeedNewResult}}result ={{if .ResultData.IsMap}}make{{else}}new{{end}}({{if ne .ResultData.Package ""}}{{.ResultData.Package}}.{{end}}{{.ResultData.Type}}){{end}}
	{{if or .ReturnRowsAffected .ReturnError}}var executeSQL *gorm.DB
	{{end}}{{if .HasSqlData}}if len(params)>0{
		{{if or .ReturnRowsAffected .ReturnError}}executeSQL{{else}}_{{end}}= {{.S}}.UnderlyingDB().{{.GormOption}}(generateSQL.String(){{if .HasSqlData}},params{{end}}){{if not .ResultData.IsNull}}.{{.GormRunMethodName}}({{if .HasGotPoint}}&{{end}}{{.ResultData.Name}}){{end}}
	}else{
		{{if or .ReturnRowsAffected .ReturnError}}executeSQL{{else}}_{{end}}= {{.S}}.UnderlyingDB().{{.GormOption}}(generateSQL.String()){{if not .ResultData.IsNull}}.{{.GormRunMethodName}}({{if .HasGotPoint}}&{{end}}{{.ResultData.Name}}){{end}}
	}{{else}}{{if or .ReturnRowsAffected .ReturnError}}executeSQL{{else}}_{{end}}= {{.S}}.UnderlyingDB().{{.GormOption}}(generateSQL.String()){{if not .ResultData.IsNull}}.{{.GormRunMethodName}}({{if .HasGotPoint}}&{{end}}{{.ResultData.Name}}){{end}}{{end}}
	{{if .ReturnRowsAffected}}rowsAffected = executeSQL.RowsAffected
	{{end}}{{if .ReturnError}}err = executeSQL.Error
	{{end}}return
}

`

const CRUDMethod = `
func ({{.S}} {{.NewStructName}}Do) Debug() I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Debug())
}

func ({{.S}} {{.NewStructName}}Do) WithContext(ctx context.Context) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.WithContext(ctx))
}

func ({{.S}} {{.NewStructName}}Do) Clauses(conds ...clause.Expression) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Clauses(conds...))
}

func ({{.S}} {{.NewStructName}}Do) Not(conds ...gen.Condition) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Not(conds...))
}

func ({{.S}} {{.NewStructName}}Do) Or(conds ...gen.Condition) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Or(conds...))
}

func ({{.S}} {{.NewStructName}}Do) Select(conds ...field.Expr) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Select(conds...))
}

func ({{.S}} {{.NewStructName}}Do) Where(conds ...gen.Condition) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Where(conds...))
}

func ({{.S}} {{.NewStructName}}Do) Exists(subquery interface{UnderlyingDB() *gorm.DB}) I{{.StructName}}Do {
	return {{.S}}.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func ({{.S}} {{.NewStructName}}Do) Order(conds ...field.Expr) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Order(conds...))
}

func ({{.S}} {{.NewStructName}}Do) Distinct(cols ...field.Expr) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Distinct(cols...))
}

func ({{.S}} {{.NewStructName}}Do) Omit(cols ...field.Expr) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Omit(cols...))
}

func ({{.S}} {{.NewStructName}}Do) Join(table schema.Tabler, on ...field.Expr) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Join(table, on...))
}

func ({{.S}} {{.NewStructName}}Do) LeftJoin(table schema.Tabler, on ...field.Expr) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.LeftJoin(table, on...))
}

func ({{.S}} {{.NewStructName}}Do) RightJoin(table schema.Tabler, on ...field.Expr) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.RightJoin(table, on...))
}

func ({{.S}} {{.NewStructName}}Do) Group(cols ...field.Expr) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Group(cols...))
}

func ({{.S}} {{.NewStructName}}Do) Having(conds ...gen.Condition) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Having(conds...))
}

func ({{.S}} {{.NewStructName}}Do) Limit(limit int) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Limit(limit))
}

func ({{.S}} {{.NewStructName}}Do) Offset(offset int) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Offset(offset))
}

func ({{.S}} {{.NewStructName}}Do) Scopes(funcs ...func(gen.Dao) gen.Dao) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Scopes(funcs...))
}

func ({{.S}} {{.NewStructName}}Do) Unscoped() I{{.StructName}}Do {
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

func ({{.S}} {{.NewStructName}}Do) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*{{.StructInfo.Package}}.{{.StructInfo.Type}}, err error) {
	buf := make([]*{{.StructInfo.Package}}.{{.StructInfo.Type}}, 0, batchSize)
	err = {{.S}}.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func ({{.S}} {{.NewStructName}}Do) FindInBatches(result *[]*{{.StructInfo.Package}}.{{.StructInfo.Type}}, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return {{.S}}.DO.FindInBatches(result, batchSize, fc)
}

func ({{.S}} {{.NewStructName}}Do) Attrs(attrs ...field.AssignExpr) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Attrs(attrs...))
}

func ({{.S}} {{.NewStructName}}Do) Assign(attrs ...field.AssignExpr) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Assign(attrs...))
}

func ({{.S}} {{.NewStructName}}Do) Joins(field field.RelationField) I{{.StructName}}Do {
	return {{.S}}.withDO({{.S}}.DO.Joins(field))
}

func ({{.S}} {{.NewStructName}}Do) Preload(field field.RelationField) I{{.StructName}}Do {
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

	if limit <= 0 {
		return
	}

	result, err = {{.S}}.Offset(offset).Limit(limit).Find()
	return
}

func ({{.S}} {{.NewStructName}}Do) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = {{.S}}.Count()
	if err != nil {
		return
	}

	err = {{.S}}.Offset(offset).Limit(limit).Scan(result)
	return
}

func ({{.S}} {{.NewStructName}}Do) Scan(result interface{}) (err error) {
	return {{.S}}.DO.Scan(result)
}

func ({{.S}} *{{.NewStructName}}Do) withDO(do gen.Dao) (*{{.NewStructName}}Do) {
	{{.S}}.DO = *do.(*gen.DO)
	return {{.S}}
}

`

const CRUDMethod_TEST = `
func init() {
	InitializeDB()
	err := db.AutoMigrate(&{{.StructInfo.Package}}.{{.StructName}}{})
	if err != nil{
		fmt.Printf("Error: AutoMigrate(&{{.StructInfo.Package}}.{{.StructName}}{}) fail: %s", err)
	}
}

func Test_{{.NewStructName}}Query(t *testing.T) {
	{{.NewStructName}} := new{{.StructName}}(db)
	{{.NewStructName}} = *{{.NewStructName}}.As({{.NewStructName}}.TableName())
	do := {{.NewStructName}}.WithContext(context.Background()).Debug()

	primaryKey := field.NewString({{.NewStructName}}.TableName(), clause.PrimaryKey)
	_, err := do.Unscoped().Where(primaryKey.IsNotNull()).Delete()
	if err != nil {
		t.Error("clean table <{{.TableName}}> fail:", err)
		return
	}
	
	_, ok := {{.NewStructName}}.GetFieldByName("")
	if ok {
		t.Error("GetFieldByName(\"\") from {{.NewStructName}} success")
	}

	err = do.Create(&{{.StructInfo.Package}}.{{.StructName}}{})
	if err != nil {
		t.Error("create item in table <{{.TableName}}> fail:", err)
	}

	err = do.Save(&{{.StructInfo.Package}}.{{.StructName}}{})
	if err != nil {
		t.Error("create item in table <{{.TableName}}> fail:", err)
	}

	err = do.CreateInBatches([]*{{.StructInfo.Package}}.{{.StructName}}{ {}, {} }, 10)
	if err != nil {
		t.Error("create item in table <{{.TableName}}> fail:", err)
	}

	_, err = do.Select({{.NewStructName}}.ALL).Take()
	if err != nil {
		t.Error("Take() on table <{{.TableName}}> fail:", err)
	}

	_, err = do.First()
	if err != nil {
		t.Error("First() on table <{{.TableName}}> fail:", err)
	}

	_, err = do.Last()
	if err != nil {
		t.Error("First() on table <{{.TableName}}> fail:", err)
	}

	_, err = do.Where(primaryKey.IsNotNull()).FindInBatch(10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatch() on table <{{.TableName}}> fail:", err)
	}

	err = do.Where(primaryKey.IsNotNull()).FindInBatches(&[]*{{.StructInfo.Package}}.{{.StructName}}{}, 10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatches() on table <{{.TableName}}> fail:", err)
	}

	_, err = do.Select({{.NewStructName}}.ALL).Where(primaryKey.IsNotNull()).Order(primaryKey.Desc()).Find()
	if err != nil {
		t.Error("Find() on table <{{.TableName}}> fail:", err)
	}

	_, err = do.Distinct(primaryKey).Take()
	if err != nil {
		t.Error("select Distinct() on table <{{.TableName}}> fail:", err)
	}

	_, err = do.Select({{.NewStructName}}.ALL).Omit(primaryKey).Take()
	if err != nil {
		t.Error("Omit() on table <{{.TableName}}> fail:", err)
	}

	_, err = do.Group(primaryKey).Find()
	if err != nil {
		t.Error("Group() on table <{{.TableName}}> fail:", err)
	}

	_, err = do.Scopes(func(dao gen.Dao) gen.Dao { return dao.Where(primaryKey.IsNotNull()) }).Find()
	if err != nil {
		t.Error("Scopes() on table <{{.TableName}}> fail:", err)
	}

	_, _, err = do.FindByPage(0, 1)
	if err != nil {
		t.Error("FindByPage() on table <{{.TableName}}> fail:", err)
	}
	
	_, err = do.ScanByPage(&{{.StructInfo.Package}}.{{.StructName}}{}, 0, 1)
	if err != nil {
		t.Error("ScanByPage() on table <{{.TableName}}> fail:", err)
	}

	_, err = do.Attrs(primaryKey).Assign(primaryKey).FirstOrInit()
	if err != nil {
		t.Error("FirstOrInit() on table <{{.TableName}}> fail:", err)
	}

	_, err = do.Attrs(primaryKey).Assign(primaryKey).FirstOrCreate()
	if err != nil {
		t.Error("FirstOrCreate() on table <{{.TableName}}> fail:", err)
	}
	
	var _a _another
	var _aPK = field.NewString(_a.TableName(), clause.PrimaryKey)

	err = do.Join(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("Join() on table <{{.TableName}}> fail:", err)
	}

	err = do.LeftJoin(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("LeftJoin() on table <{{.TableName}}> fail:", err)
	}
	
	_, err = do.Not().Or().Clauses().Take()
	if err != nil {
		t.Error("Not/Or/Clauses on table <{{.TableName}}> fail:", err)
	}
}
`

const DIYMethod_TEST_Basic = `
type Input struct {
	Args []interface{}
}

type Expectation struct {
	Ret []interface{}
}

type TestCase struct {
	Input
	Expectation
}

`

const DIYMethod_TEST = `

var {{.OriginStruct.Type}}{{.MethodName}}TestCase = []TestCase{}

func Test_{{.TargetStruct}}_{{.MethodName}}(t *testing.T) {
	{{.TargetStruct}} := new{{.OriginStruct.Type}}(db)
	do := {{.TargetStruct}}.WithContext(context.Background()).Debug()

	for i, tt := range {{.OriginStruct.Type}}{{.MethodName}}TestCase {
		t.Run("{{.MethodName}}_"+strconv.Itoa(i), func(t *testing.T) {
			{{.GetTestResultParamInTmpl}} := do.{{.MethodName}}({{.GetTestParamInTmpl}})
			{{.GetAssertInTmpl}}
		})
	}
}

`
