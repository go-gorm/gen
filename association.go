package gen

import "gorm.io/gorm"

type Association struct{ gorm.Association }

func (d *DO) Association(column string) *Association {
	return &Association{*d.db.Association(column)}
}

func (association *Association) Find(out interface{}) error {
	return association.Association.Find(out)
}

func (association *Association) Append(values ...interface{}) error {
	return association.Association.Append(values...)
}

func (association *Association) Replace(values ...interface{}) error {
	return association.Association.Replace(values...)
}

func (association *Association) Delete(values ...interface{}) error {
	return association.Association.Delete(values...)
}

func (association *Association) Clear() error {
	return association.Association.Clear()
}

func (association *Association) Count() int64 {
	return association.Association.Count()
}
