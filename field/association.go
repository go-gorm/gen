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
	varName string
	varType string
	path    string

	relations []*Relation

	conds   []Expr
	order   []Expr
	clauses []clause.Expression
}

func (r Relation) Name() string { return r.varName }

func (r Relation) Path() string { return r.path }

func (r Relation) Type() string { return r.varType }

func (r Relation) Field(member ...string) Expr {
	if len(member) > 0 {
		return NewString("", r.varName+"."+strings.Join(member, ".")).appendBuildOpts(WithoutQuote)
	}
	return NewString("", r.varName).appendBuildOpts(WithoutQuote)
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
	for _, relation := range r.relations {
		memberStr += relation.varName + " struct {\nfield.RelationField\n" + relation.StructMember() + "}\n"
	}
	return memberStr
}

func (r *Relation) StructMemberInit() string {
	initStr := fmt.Sprintf("RelationField: field.NewRelation(%q, %q),\n", r.path, r.varType)
	for _, relation := range r.relations {
		initStr += relation.varName + ": struct {\nfield.RelationField\n" + strings.TrimSpace(relation.StructMember()) + "}"
		initStr += "{\n" + relation.StructMemberInit() + "},\n"
	}
	return initStr
}

func wrapPath(root string, rs []*Relation) []*Relation {
	for _, r := range rs {
		r.path = root + "." + r.path
		r.relations = wrapPath(root, r.relations)
	}
	return rs
}
