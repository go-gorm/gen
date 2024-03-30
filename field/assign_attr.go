package field

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
)

var testDB, _ = gorm.Open(tests.DummyDialector{}, nil)

type IValues interface {
	Values() interface{}
}

type attrs struct {
	expr
	value        interface{}
	db           *gorm.DB
	selectFields []IColumnName
	omitFields   []IColumnName
}

func (att *attrs) AssignExpr() expression {
	return att
}

func (att *attrs) BeCond() interface{} {
	return att.db.Statement.BuildCondition(att.Values())
}

func (att *attrs) Values() interface{} {
	if att == nil || att.value == nil {
		return nil
	}
	if len(att.selectFields) == 0 && len(att.omitFields) == 0 {
		return att.value
	}
	values := make(map[string]interface{})
	if value, ok := att.value.(map[string]interface{}); ok {
		values = value
	} else if value, ok := att.value.(*map[string]interface{}); ok {
		values = *value
	} else {
		reflectValue := reflect.Indirect(reflect.ValueOf(att.value))
		for reflectValue.Kind() == reflect.Ptr || reflectValue.Kind() == reflect.Interface {
			reflectValue = reflect.Indirect(reflectValue)
		}
		switch reflectValue.Kind() {
		case reflect.Struct:
			if err := att.db.Statement.Parse(att.value); err == nil {
				ignoreZero := len(att.selectFields) == 0
				for _, f := range att.db.Statement.Schema.Fields {
					if f.Readable {
						if v, isZero := f.ValueOf(att.db.Statement.Context, reflectValue); !isZero || !ignoreZero {
							values[f.DBName] = v
						}
					}
				}
			}
		}
	}
	if len(att.selectFields) > 0 {
		fm, all := toFieldMap(att.selectFields)
		if all {
			return values
		}
		tvs := make(map[string]interface{}, len(fm))
		for fn, vl := range values {
			if fm[fn] {
				tvs[fn] = vl
			}
		}
		return tvs
	}
	fm, all := toFieldMap(att.omitFields)
	if all {
		return map[string]interface{}{}
	}
	for fn := range fm {
		delete(values, fn)
	}
	return values
}

func toFieldMap(fields []IColumnName) (fieldsMap map[string]bool, all bool) {
	fieldsMap = make(map[string]bool, len(fields))
	for _, f := range fields {
		if strings.HasSuffix(string(f.ColumnName()), "*") {
			all = true
			return
		}
		fieldsMap[string(f.ColumnName())] = true
	}
	return
}

func (att *attrs) Select(fields ...IColumnName) *attrs {
	if att == nil || att.db == nil {
		return att
	}
	att.selectFields = fields
	return att
}

func (att *attrs) Omit(fields ...IColumnName) *attrs {
	if att == nil || att.db == nil {
		return att
	}
	att.omitFields = fields
	return att
}

func Attrs(attr interface{}) *attrs {
	res := &attrs{db: testDB.Debug()}
	if attr != nil {
		res.value = attr
	}
	return res
}
