package field

import (
	"fmt"
	"strings"

	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type RelationshipType schema.RelationshipType

const (
	HasOne    RelationshipType = RelationshipType(schema.HasOne)    // HasOneRel has one relationship
	HasMany   RelationshipType = RelationshipType(schema.HasMany)   // HasManyRel has many relationships
	BelongsTo RelationshipType = RelationshipType(schema.BelongsTo) // BelongsToRel belongs to relationship
	Many2Many RelationshipType = RelationshipType(schema.Many2Many) // Many2ManyRel many to many relationship
)

type Relations struct {
	HasOne    []*Relation
	BelongsTo []*Relation
	HasMany   []*Relation
	Many2Many []*Relation
}

func (r *Relations) Accept(relations ...*Relation) {
	for _, relation := range relations {
		switch relation.Relationship() {
		case HasOne:
			r.HasOne = append(r.HasOne, relation)
		case HasMany:
			r.HasMany = append(r.HasMany, relation)
		case BelongsTo:
			r.BelongsTo = append(r.BelongsTo, relation)
		case Many2Many:
			r.Many2Many = append(r.Many2Many, relation)
		}
	}
}

func (r *Relations) SingleRelation() []*Relation {
	return append(append(append(append(make([]*Relation, 0, 4), r.HasOne...), r.BelongsTo...), r.HasMany...), r.Many2Many...)
}

type RelationField interface {
	Name() string
	Path() string
	Field(member ...string) Expr

	On(conds ...Expr) RelationField
	Order(columns ...Expr) RelationField
	Clauses(hints ...clause.Expression) RelationField

	GetConds() []Expr
	GetOrderCol() []Expr
	GetClauses() []clause.Expression
}

type Relation struct {
	relationship RelationshipType

	fieldName  string
	fieldType  string
	fieldPath  string
	fieldModel interface{} // store relaiton model

	childRelations []*Relation

	conds   []Expr
	order   []Expr
	clauses []clause.Expression
}

func (r Relation) Name() string { return r.fieldName }

func (r Relation) Path() string { return r.fieldPath }

func (r Relation) Type() string { return r.fieldType }

func (r Relation) Model() interface{} { return r.fieldModel }

func (r Relation) Relationship() RelationshipType { return r.relationship }

func (r Relation) Field(member ...string) Expr {
	if len(member) > 0 {
		return NewString("", r.fieldName+"."+strings.Join(member, ".")).appendBuildOpts(WithoutQuote)
	}
	return NewString("", r.fieldName).appendBuildOpts(WithoutQuote)
}

func (r *Relation) AppendChildRelation(relations ...*Relation) {
	r.childRelations = append(r.childRelations, wrapPath(r.fieldPath, relations)...)
}

func (r *Relation) On(conds ...Expr) RelationField {
	r.conds = append(r.conds, conds...)
	return r
}

func (r *Relation) Order(columns ...Expr) RelationField {
	r.order = append(r.order, columns...)
	return r
}

func (r *Relation) Clauses(hints ...clause.Expression) RelationField {
	r.clauses = append(r.clauses, hints...)
	return r
}

func (r *Relation) GetConds() []Expr { return r.conds }

func (r *Relation) GetOrderCol() []Expr { return r.order }

func (r *Relation) GetClauses() []clause.Expression { return r.clauses }

func (r *Relation) StructMember() string {
	var memberStr string
	for _, relation := range r.childRelations {
		memberStr += relation.fieldName + " struct {\nfield.RelationField\n" + relation.StructMember() + "}\n"
	}
	return memberStr
}

func (r *Relation) StructMemberInit() string {
	initStr := fmt.Sprintf("RelationField: field.NewRelation(%q, %q),\n", r.fieldPath, r.fieldType)
	for _, relation := range r.childRelations {
		initStr += relation.fieldName + ": struct {\nfield.RelationField\n" + strings.TrimSpace(relation.StructMember()) + "}"
		initStr += "{\n" + relation.StructMemberInit() + "},\n"
	}
	return initStr
}

func wrapPath(root string, rs []*Relation) []*Relation {
	for _, r := range rs {
		r.fieldPath = root + "." + r.fieldPath
		r.childRelations = wrapPath(root, r.childRelations)
	}
	return rs
}

var defaultRelationshipPrefix = map[RelationshipType]string{
	// HasOne:    "",
	// BelongsTo: "",
	HasMany:   "[]",
	Many2Many: "[]",
}

type RelateConfig struct {
	RelatePointer      bool
	RelateSlice        bool
	RelateSlicePointer bool

	JSONTag      string
	GORMTag      string
	NewTag       string
	OverwriteTag string
}

func (c *RelateConfig) RelateFieldPrefix(relationshipType RelationshipType) string {
	switch {
	case c.RelatePointer:
		return "*"
	case c.RelateSlice:
		return "[]"
	case c.RelateSlicePointer:
		return "[]*"
	default:
		return defaultRelationshipPrefix[relationshipType]
	}
}
