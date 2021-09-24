package field

import (
	"fmt"
	"strings"

	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

var s = clause.Associations

type RelationshipType schema.RelationshipType

const (
	HasOne    RelationshipType = "has_one"      // HasOneRel has one relationship
	HasMany   RelationshipType = "has_many"     // HasManyRel has many relationship
	BelongsTo RelationshipType = "belongs_to"   // BelongsToRel belongs to relationship
	Many2Many RelationshipType = "many_to_many" // Many2ManyRel many to many relationship
)

type Relations struct {
	HasOne    []*Relation
	BelongsTo []*Relation
	HasMany   []*Relation
	Many2Many []*Relation
}

type RelationPath interface {
	Path() relationPath
}

type relationPath string

func (p relationPath) Path() relationPath { return p }

type Relation struct {
	varName string
	varType string
	path    string

	relations []*Relation
}

func (r Relation) Name() string { return r.varName }

func (r Relation) Path() relationPath { return relationPath(r.path) }

func (r Relation) Type() string { return r.varType }

func (r *Relation) StructMember() string {
	var memberStr string
	for _, relation := range r.relations {
		memberStr += relation.varName + " struct {\nfield.Relation\n" + relation.StructMember() + "}\n"
	}
	return memberStr
}

func (r *Relation) StructMemberInit() string {
	initStr := fmt.Sprintf("Relation: *field.NewRelation(%q, %q),\n", r.path, r.varType)
	for _, relation := range r.relations {
		initStr += relation.varName + ": struct {\nfield.Relation\n" + strings.TrimPrefix(strings.TrimSpace(relation.StructMember()), relation.varName) + "}"
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
