package template

// DIYMethod DIY method
const DIYMethod = `

// {{.DocComment }}
func ({{.S}} {{.TargetStruct}}Do){{.FuncSign}}{
	{{if .HasSQLData}}var params []interface{}

	{{end}}var generateSQL strings.Builder
	{{range $line:=.Section.Tmpls}}{{$line}}
	{{end}}

	{{if .HasNeedNewResult}}result ={{if .ResultData.IsMap}}make{{else}}new{{end}}({{if ne .ResultData.Package ""}}{{.ResultData.Package}}.{{end}}{{.ResultData.Type}}){{end}}
	{{if .ReturnSQLResult}}stmt := {{.S}}.UnderlyingDB().Statement
	result,{{if .ReturnError}}err{{else}}_{{end}} = stmt.ConnPool.ExecContext(stmt.Context,generateSQL.String(){{if .HasSQLData}},params...{{end}}) // ignore_security_alert
	{{else if .ReturnSQLRow}}row = {{.S}}.UnderlyingDB().Raw(generateSQL.String(){{if .HasSQLData}},params...{{end}}).Row() // ignore_security_alert
	{{else if .ReturnSQLRows}}rows,{{if .ReturnError}}err{{else}}_{{end}} = {{.S}}.UnderlyingDB().Raw(generateSQL.String(){{if .HasSQLData}},params...{{end}}).Rows() // ignore_security_alert
	{{else}}var executeSQL *gorm.DB
	executeSQL = {{.S}}.UnderlyingDB().{{.GormOption}}(generateSQL.String(){{if .HasSQLData}},params...{{end}}){{if not .ResultData.IsNull}}.{{.GormRunMethodName}}({{if .HasGotPoint}}&{{end}}{{.ResultData.Name}}){{end}}  // ignore_security_alert
	{{if .ReturnRowsAffected}}rowsAffected = executeSQL.RowsAffected
	{{end}}{{if .ReturnError}}err = executeSQL.Error
	{{end}}{{if .ReturnNothing}}_ = executeSQL
	{{end}}{{end}}
	return
}

`

// CRUDMethod CRUD method
const CRUDMethod = `
func ({{.S}} {{.QueryStructName}}Do) Debug() {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Debug())
}

func ({{.S}} {{.QueryStructName}}Do) WithContext(ctx context.Context) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.WithContext(ctx))
}

func ({{.S}} {{.QueryStructName}}Do) ReadDB() {{.ReturnObject}} {
	return {{.S}}.Clauses(dbresolver.Read)
}

func ({{.S}} {{.QueryStructName}}Do) WriteDB() {{.ReturnObject}} {
	return {{.S}}.Clauses(dbresolver.Write)
}

func ({{.S}} {{.QueryStructName}}Do) Session(config *gorm.Session) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Session(config))
}

func ({{.S}} {{.QueryStructName}}Do) Clauses(conds ...clause.Expression) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Clauses(conds...))
}

func ({{.S}} {{.QueryStructName}}Do) Returning(value interface{}, columns ...string) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Returning(value, columns...))
}

func ({{.S}} {{.QueryStructName}}Do) Not(conds ...gen.Condition) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Not(conds...))
}

func ({{.S}} {{.QueryStructName}}Do) Or(conds ...gen.Condition) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Or(conds...))
}

func ({{.S}} {{.QueryStructName}}Do) Select(conds ...field.Expr) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Select(conds...))
}

func ({{.S}} {{.QueryStructName}}Do) Where(conds ...gen.Condition) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Where(conds...))
}

func ({{.S}} {{.QueryStructName}}Do) Order(conds ...field.Expr) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Order(conds...))
}

func ({{.S}} {{.QueryStructName}}Do) Distinct(cols ...field.Expr) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Distinct(cols...))
}

func ({{.S}} {{.QueryStructName}}Do) Omit(cols ...field.Expr) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Omit(cols...))
}

func ({{.S}} {{.QueryStructName}}Do) Join(table schema.Tabler, on ...field.Expr) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Join(table, on...))
}

func ({{.S}} {{.QueryStructName}}Do) LeftJoin(table schema.Tabler, on ...field.Expr) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.LeftJoin(table, on...))
}

func ({{.S}} {{.QueryStructName}}Do) RightJoin(table schema.Tabler, on ...field.Expr) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.RightJoin(table, on...))
}

func ({{.S}} {{.QueryStructName}}Do) Group(cols ...field.Expr) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Group(cols...))
}

func ({{.S}} {{.QueryStructName}}Do) Having(conds ...gen.Condition) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Having(conds...))
}

func ({{.S}} {{.QueryStructName}}Do) Limit(limit int) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Limit(limit))
}

func ({{.S}} {{.QueryStructName}}Do) Offset(offset int) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Offset(offset))
}

func ({{.S}} {{.QueryStructName}}Do) Scopes(funcs ...func(gen.Dao) gen.Dao) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Scopes(funcs...))
}

func ({{.S}} {{.QueryStructName}}Do) Unscoped() {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Unscoped())
}

func ({{.S}} {{.QueryStructName}}Do) Create(values ...*{{.StructInfo.Package}}.{{.StructInfo.Type}}) error {
	if len(values) == 0 {
		return nil
	}
	return {{.S}}.DO.Create(values)
}

func ({{.S}} {{.QueryStructName}}Do) CreateInBatches(values []*{{.StructInfo.Package}}.{{.StructInfo.Type}}, batchSize int) error {
	return {{.S}}.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func ({{.S}} {{.QueryStructName}}Do) Save(values ...*{{.StructInfo.Package}}.{{.StructInfo.Type}}) error {
	if len(values) == 0 {
		return nil
	}
	return {{.S}}.DO.Save(values)
}

func ({{.S}} {{.QueryStructName}}Do) First() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error) {
	if result, err := {{.S}}.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*{{.StructInfo.Package}}.{{.StructInfo.Type}}), nil
	}
}

func ({{.S}} {{.QueryStructName}}Do) Take() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error) {
	if result, err := {{.S}}.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*{{.StructInfo.Package}}.{{.StructInfo.Type}}), nil
	}
}

func ({{.S}} {{.QueryStructName}}Do) Last() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error) {
	if result, err := {{.S}}.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*{{.StructInfo.Package}}.{{.StructInfo.Type}}), nil
	}
}

func ({{.S}} {{.QueryStructName}}Do) Find() ([]*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error) {
	result, err := {{.S}}.DO.Find()
	return result.([]*{{.StructInfo.Package}}.{{.StructInfo.Type}}), err
}

func ({{.S}} {{.QueryStructName}}Do) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*{{.StructInfo.Package}}.{{.StructInfo.Type}}, err error) {
	buf := make([]*{{.StructInfo.Package}}.{{.StructInfo.Type}}, 0, batchSize)
	err = {{.S}}.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func ({{.S}} {{.QueryStructName}}Do) FindInBatches(result *[]*{{.StructInfo.Package}}.{{.StructInfo.Type}}, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return {{.S}}.DO.FindInBatches(result, batchSize, fc)
}

func ({{.S}} {{.QueryStructName}}Do) Attrs(attrs ...field.AssignExpr) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Attrs(attrs...))
}

func ({{.S}} {{.QueryStructName}}Do) Assign(attrs ...field.AssignExpr) {{.ReturnObject}} {
	return {{.S}}.withDO({{.S}}.DO.Assign(attrs...))
}

func ({{.S}} {{.QueryStructName}}Do) Joins(fields ...field.RelationField) {{.ReturnObject}} {
	for _, _f := range fields {
        {{.S}} = *{{.S}}.withDO({{.S}}.DO.Joins(_f))
    }
	return &{{.S}}
}

func ({{.S}} {{.QueryStructName}}Do) Preload(fields ...field.RelationField) {{.ReturnObject}} {
    for _, _f := range fields {
        {{.S}} = *{{.S}}.withDO({{.S}}.DO.Preload(_f))
    }
	return &{{.S}}
}

func ({{.S}} {{.QueryStructName}}Do) FirstOrInit() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error) {
	if result, err := {{.S}}.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*{{.StructInfo.Package}}.{{.StructInfo.Type}}), nil
	}
}

func ({{.S}} {{.QueryStructName}}Do) FirstOrCreate() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error) {
	if result, err := {{.S}}.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*{{.StructInfo.Package}}.{{.StructInfo.Type}}), nil
	}
}

func ({{.S}} {{.QueryStructName}}Do) FindByPage(offset int, limit int) (result []*{{.StructInfo.Package}}.{{.StructInfo.Type}}, count int64, err error) {
	result, err = {{.S}}.Offset(offset).Limit(limit).Find()
	if err != nil{
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size+offset)
		return
	}

	count, err = {{.S}}.Offset(-1).Limit(-1).Count()
	return
}

func ({{.S}} {{.QueryStructName}}Do) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = {{.S}}.Count()
	if err != nil {
		return
	}

	err = {{.S}}.Offset(offset).Limit(limit).Scan(result)
	return
}

func ({{.S}} {{.QueryStructName}}Do) Scan(result interface{}) (err error) {
	return {{.S}}.DO.Scan(result)
}

func ({{.S}} {{.QueryStructName}}Do) Delete(models ...*{{.StructInfo.Package}}.{{.StructInfo.Type}}) (result gen.ResultInfo, err error) {
	return {{.S}}.DO.Delete(models)
}

func ({{.S}} *{{.QueryStructName}}Do) withDO(do gen.Dao) (*{{.QueryStructName}}Do) {
	{{.S}}.DO = *do.(*gen.DO)
	return {{.S}}
}

`

// CRUDMethodTest CRUD method test
const CRUDMethodTest = `
func init() {
	InitializeDB()
	err := _gen_test_db.AutoMigrate(&{{.StructInfo.Package}}.{{.ModelStructName}}{})
	if err != nil{
		fmt.Printf("Error: AutoMigrate(&{{.StructInfo.Package}}.{{.ModelStructName}}{}) fail: %s", err)
	}
}

func Test_{{.QueryStructName}}Query(t *testing.T) {
	{{.QueryStructName}} := new{{.ModelStructName}}(_gen_test_db)
	{{.QueryStructName}} = *{{.QueryStructName}}.As({{.QueryStructName}}.TableName())
	_do := {{.QueryStructName}}.WithContext(context.Background()).Debug()

	primaryKey := field.NewString({{.QueryStructName}}.TableName(), clause.PrimaryKey)
	_, err := _do.Unscoped().Where(primaryKey.IsNotNull()).Delete()
	if err != nil {
		t.Error("clean table <{{.TableName}}> fail:", err)
		return
	}
	
	_, ok := {{.QueryStructName}}.GetFieldByName("")
	if ok {
		t.Error("GetFieldByName(\"\") from {{.QueryStructName}} success")
	}

	err = _do.Create(&{{.StructInfo.Package}}.{{.ModelStructName}}{})
	if err != nil {
		t.Error("create item in table <{{.TableName}}> fail:", err)
	}

	err = _do.Save(&{{.StructInfo.Package}}.{{.ModelStructName}}{})
	if err != nil {
		t.Error("create item in table <{{.TableName}}> fail:", err)
	}

	err = _do.CreateInBatches([]*{{.StructInfo.Package}}.{{.ModelStructName}}{ {}, {} }, 10)
	if err != nil {
		t.Error("create item in table <{{.TableName}}> fail:", err)
	}

	_, err = _do.Select({{.QueryStructName}}.ALL).Take()
	if err != nil {
		t.Error("Take() on table <{{.TableName}}> fail:", err)
	}

	_, err = _do.First()
	if err != nil {
		t.Error("First() on table <{{.TableName}}> fail:", err)
	}

	_, err = _do.Last()
	if err != nil {
		t.Error("First() on table <{{.TableName}}> fail:", err)
	}

	_, err = _do.Where(primaryKey.IsNotNull()).FindInBatch(10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatch() on table <{{.TableName}}> fail:", err)
	}

	err = _do.Where(primaryKey.IsNotNull()).FindInBatches(&[]*{{.StructInfo.Package}}.{{.ModelStructName}}{}, 10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatches() on table <{{.TableName}}> fail:", err)
	}

	_, err = _do.Select({{.QueryStructName}}.ALL).Where(primaryKey.IsNotNull()).Order(primaryKey.Desc()).Find()
	if err != nil {
		t.Error("Find() on table <{{.TableName}}> fail:", err)
	}

	_, err = _do.Distinct(primaryKey).Take()
	if err != nil {
		t.Error("select Distinct() on table <{{.TableName}}> fail:", err)
	}

	_, err = _do.Select({{.QueryStructName}}.ALL).Omit(primaryKey).Take()
	if err != nil {
		t.Error("Omit() on table <{{.TableName}}> fail:", err)
	}

	_, err = _do.Group(primaryKey).Find()
	if err != nil {
		t.Error("Group() on table <{{.TableName}}> fail:", err)
	}

	_, err = _do.Scopes(func(dao gen.Dao) gen.Dao { return dao.Where(primaryKey.IsNotNull()) }).Find()
	if err != nil {
		t.Error("Scopes() on table <{{.TableName}}> fail:", err)
	}

	_, _, err = _do.FindByPage(0, 1)
	if err != nil {
		t.Error("FindByPage() on table <{{.TableName}}> fail:", err)
	}
	
	_, err = _do.ScanByPage(&{{.StructInfo.Package}}.{{.ModelStructName}}{}, 0, 1)
	if err != nil {
		t.Error("ScanByPage() on table <{{.TableName}}> fail:", err)
	}

	_, err = _do.Attrs(primaryKey).Assign(primaryKey).FirstOrInit()
	if err != nil {
		t.Error("FirstOrInit() on table <{{.TableName}}> fail:", err)
	}

	_, err = _do.Attrs(primaryKey).Assign(primaryKey).FirstOrCreate()
	if err != nil {
		t.Error("FirstOrCreate() on table <{{.TableName}}> fail:", err)
	}
	
	var _a _another
	var _aPK = field.NewString(_a.TableName(), "id")

	err = _do.Join(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("Join() on table <{{.TableName}}> fail:", err)
	}

	err = _do.LeftJoin(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("LeftJoin() on table <{{.TableName}}> fail:", err)
	}
	
	_, err = _do.Not().Or().Clauses().Take()
	if err != nil {
		t.Error("Not/Or/Clauses on table <{{.TableName}}> fail:", err)
	}
}
`

// DIYMethodTestBasic DIY method test basic
const DIYMethodTestBasic = `
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

// DIYMethodTest DIY method test
const DIYMethodTest = `

var {{.OriginStruct.Type}}{{.MethodName}}TestCase = []TestCase{}

func Test_{{.TargetStruct}}_{{.MethodName}}(t *testing.T) {
	{{.TargetStruct}} := new{{.OriginStruct.Type}}(_gen_test_db)
	do := {{.TargetStruct}}.WithContext(context.Background()).Debug()

	for i, tt := range {{.OriginStruct.Type}}{{.MethodName}}TestCase {
		t.Run("{{.MethodName}}_"+strconv.Itoa(i), func(t *testing.T) {
			{{if .GetTestResultParamInTmpl}}{{.GetTestResultParamInTmpl}} := do.{{.MethodName}}({{.GetTestParamInTmpl}})
			{{.GetAssertInTmpl}}{{else}}do.{{.MethodName}}({{.GetTestParamInTmpl}}){{end}}
		})
	}
}

`
