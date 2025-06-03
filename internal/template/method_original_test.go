package template_test

// CRUDMethodTest CRUD method test
const CRUDMethodTest = `
func init() {
	InitializeDB()
	err := _gen_test_db.AutoMigrate(&{{.StructInfo.Package}}.{{.ModelStructName}}{})
	if err != nil {
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

	t.Run("FieldByName", func(t *testing.T) {
		_, ok := {{.QueryStructName}}.GetFieldByName("")
		if ok {
			t.Error("GetFieldByName(\"\") from {{.QueryStructName}} success")
		}
	})

	t.Run("Create", func(t *testing.T) {
		err := _do.Create(&{{.StructInfo.Package}}.{{.ModelStructName}}{})
		if err != nil {
			t.Error("create item in table <{{.TableName}}> fail:", err)
		}
	})

	t.Run("Save", func(t *testing.T) {
		err := _do.Save(&{{.StructInfo.Package}}.{{.ModelStructName}}{})
		if err != nil {
			t.Error("create item in table <{{.TableName}}> fail:", err)
		}
	})

	t.Run("CreateInBatches", func(t *testing.T) {
		err := _do.CreateInBatches([]*{{.StructInfo.Package}}.{{.ModelStructName}}{ {}, {} }, 10)
		if err != nil {
			t.Error("create item in table <{{.TableName}}> fail:", err)
		}
	})

	t.Run("Take", func(t *testing.T) {
		_, err := _do.Select({{.QueryStructName}}.ALL).Take()
		if err != nil {
			t.Error("Take() on table <{{.TableName}}> fail:", err)
		}
	})

	t.Run("First", func(t *testing.T) {
		_, err := _do.First()
		if err != nil {
			t.Error("First() on table <{{.TableName}}> fail:", err)
		}
	})

	t.Run("Last", func(t *testing.T) {
		_, err := _do.Last()
		if err != nil {
			t.Error("First() on table <{{.TableName}}> fail:", err)
		}
	})

	t.Run("FindInBatch", func(t *testing.T) {
		_, err := _do.Where(primaryKey.IsNotNull()).FindInBatch(10, func(tx gen.Dao, batch int) error { return nil })
		if err != nil {
			t.Error("FindInBatch() on table <{{.TableName}}> fail:", err)
		}
	})

	t.Run("FindInBatches", func(t *testing.T) {
		err := _do.Where(primaryKey.IsNotNull()).FindInBatches(&[]*{{.StructInfo.Package}}.{{.ModelStructName}}{}, 10, func(tx gen.Dao, batch int) error { return nil })
		if err != nil {
			t.Error("FindInBatches() on table <{{.TableName}}> fail:", err)
		}
	})

	t.Run("Find", func(t *testing.T) {
		_, err := _do.Select({{.QueryStructName}}.ALL).Where(primaryKey.IsNotNull()).Order(primaryKey.Desc()).Find()
		if err != nil {
			t.Error("Find() on table <{{.TableName}}> fail:", err)
		}
	})

	t.Run("Distinct", func(t *testing.T) {
		_, err := _do.Distinct(primaryKey).Take()
		if err != nil {
			t.Error("select Distinct() on table <{{.TableName}}> fail:", err)
		}
	})

	t.Run("Omit", func(t *testing.T) {
		_, err := _do.Select({{.QueryStructName}}.ALL).Omit(primaryKey).Take()
		if err != nil {
			t.Error("Omit() on table <{{.TableName}}> fail:", err)
		}
	})

	t.Run("Group", func(t *testing.T) {
		_, err := _do.Group(primaryKey).Find()
		if err != nil {
			t.Error("Group() on table <{{.TableName}}> fail:", err)
		}
	})

	t.Run("Scopes", func(t *testing.T) {
		_, err := _do.Scopes(func(dao gen.Dao) gen.Dao { return dao.Where(primaryKey.IsNotNull()) }).Find()
		if err != nil {
			t.Error("Scopes() on table <{{.TableName}}> fail:", err)
		}
	})

	t.Run("FindByPage", func(t *testing.T) {
		_, _, err := _do.FindByPage(0, 1)
		if err != nil {
			t.Error("FindByPage() on table <{{.TableName}}> fail:", err)
		}
	})

	t.Run("ScanByPage", func(t *testing.T) {
		_, err := _do.ScanByPage(&{{.StructInfo.Package}}.{{.ModelStructName}}{}, 0, 1)
		if err != nil {
			t.Error("ScanByPage() on table <{{.TableName}}> fail:", err)
		}
	})

	t.Run("FirstOrInit", func(t *testing.T) {
		_, err := _do.Attrs(primaryKey).Assign(primaryKey).FirstOrInit()
		if err != nil {
			t.Error("FirstOrInit() on table <{{.TableName}}> fail:", err)
		}
	})

	t.Run("FirstOrCreate", func(t *testing.T) {
		_, err := _do.Attrs(primaryKey).Assign(primaryKey).FirstOrCreate()
		if err != nil {
			t.Error("FirstOrCreate() on table <{{.TableName}}> fail:", err)
		}
	})

	t.Run("Join", func(t *testing.T) {
		var _a _another
		var _aPK = field.NewString(_a.TableName(), "id")

		t.Run("Basic Join", func(t *testing.T) {
			err := _do.Join(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
			if err != nil {
				t.Error("Join() on table <{{.TableName}}> fail:", err)
			}
		})

		t.Run("LeftJoin", func(t *testing.T) {
			err := _do.LeftJoin(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
			if err != nil {
				t.Error("LeftJoin() on table <{{.TableName}}> fail:", err)
			}
		})

		t.Run("NotOrClauses", func(t *testing.T) {
			_, err := _do.Not().Or().Clauses().Take()
			if err != nil {
				t.Error("Not/Or/Clauses on table <{{.TableName}}> fail:", err)
			}
		})


	})
}
`
