package field

type RelationPath interface {
	Path() relationPath
}

type relationPath string

func (p relationPath) Path() relationPath { return p }

type Relation struct {
	varName string
}

func (r Relation) Path() relationPath {
	return relationPath(r.varName)
}
