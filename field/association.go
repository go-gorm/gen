package field

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// RelationshipType table relationship
type RelationshipType schema.RelationshipType

const (
	// HasOne a has one association sets up a one-to-one connection with another model. Reference https://gorm.io/docs/has_one.html
	HasOne RelationshipType = RelationshipType(schema.HasOne) // HasOneRel has one relationship
	// HasMany a has many association sets up a one-to-many connection with another model. Reference https://gorm.io/docs/has_many.html
	HasMany RelationshipType = RelationshipType(schema.HasMany) // HasManyRel has many relationships
	// BelongsTo A belongs to association sets up a one-to-one connection with another model. Reference https://gorm.io/docs/belongs_to.html
	BelongsTo RelationshipType = RelationshipType(schema.BelongsTo) // BelongsToRel belongs to relationship
	// Many2Many Many to Many add a join table between two models. Reference https://gorm.io/docs/many2many.html
	Many2Many RelationshipType = RelationshipType(schema.Many2Many) // Many2ManyRel many to many relationship
)

type relationScope func(*gorm.DB) *gorm.DB

var (
	// RelationFieldUnscoped relation fild unscoped
	RelationFieldUnscoped relationScope = func(tx *gorm.DB) *gorm.DB {
		return tx.Unscoped()
	}
)

var ns = schema.NamingStrategy{}

// RelationField interface for relation field
type RelationField interface {
	Name() string
	Path() string
	// Field return expr for Select
	// Field() return "<self>" field name in struct
	// Field("RelateField") return "<self>.RelateField" for Select
	// Field("RelateField", "RelateRelateField") return "<self>.RelateField.RelateRelateField" for Select
	// ex:
	// 	Select(u.CreditCards.Field()) equals to GORM: Select("CreditCards")
	// 	Select(u.CreditCards.Field("Bank")) equals to GORM: Select("CreditCards.Bank")
	// 	Select(u.CreditCards.Field("Bank","Owner")) equals to GORM: Select("CreditCards.Bank.Owner")
	Field(fields ...string) Expr

	On(conds ...Expr) RelationField
	Select(conds ...Expr) RelationField
	Order(columns ...Expr) RelationField
	Clauses(hints ...clause.Expression) RelationField
	Scopes(funcs ...relationScope) RelationField
	Offset(offset int) RelationField
	Limit(limit int) RelationField

	GetConds() []Expr
	GetSelects() []Expr
	GetOrderCol() []Expr
	GetClauses() []clause.Expression
	GetScopes() []relationScope
	GetPage() (offset, limit int)
}

// Relation relation meta info
type Relation struct {
	relationship RelationshipType

	fieldName  string
	fieldType  string
	fieldPath  string
	fieldModel interface{} // store relaiton model

	childRelations []Relation

	conds         []Expr
	selects       []Expr
	order         []Expr
	clauses       []clause.Expression
	scopes        []relationScope
	limit, offset int
}

// Name relation field' name
func (r Relation) Name() string { return r.fieldName }

// Path relation field's path
func (r Relation) Path() string { return r.fieldPath }

// Type relation field's type
func (r Relation) Type() string { return r.fieldType }

// Model relation field's model
func (r Relation) Model() interface{} { return r.fieldModel }

// Relationship relationship between field and table struct
func (r Relation) Relationship() RelationshipType { return r.relationship }

// RelationshipName relationship's name
func (r Relation) RelationshipName() string { return ns.SchemaName(string(r.relationship)) }

// ChildRelations return child relations
func (r Relation) ChildRelations() []Relation { return r.childRelations }

// Field build field
func (r Relation) Field(fields ...string) Expr {
	if len(fields) > 0 {
		return NewString("", r.fieldName+"."+strings.Join(fields, ".")).appendBuildOpts(WithoutQuote)
	}
	return NewString("", r.fieldName).appendBuildOpts(WithoutQuote)
}

// AppendChildRelation append child relationship
func (r *Relation) AppendChildRelation(relations ...Relation) {
	r.childRelations = append(r.childRelations, wrapPath(r.fieldPath, relations)...)
}

// On relation condition
func (r Relation) On(conds ...Expr) RelationField {
	r.conds = append(r.conds, conds...)
	return &r
}

// Select relation select columns
func (r Relation) Select(columns ...Expr) RelationField {
	r.selects = append(r.selects, columns...)
	return &r
}

// Order relation order columns
func (r Relation) Order(columns ...Expr) RelationField {
	r.order = append(r.order, columns...)
	return &r
}

// Clauses set relation clauses
func (r Relation) Clauses(hints ...clause.Expression) RelationField {
	r.clauses = append(r.clauses, hints...)
	return &r
}

// Scopes set scopes func
func (r Relation) Scopes(funcs ...relationScope) RelationField {
	r.scopes = append(r.scopes, funcs...)
	return &r
}

// Offset set relation offset
func (r Relation) Offset(offset int) RelationField {
	r.offset = offset
	return &r
}

// Limit set relation limit
func (r Relation) Limit(limit int) RelationField {
	r.limit = limit
	return &r
}

// GetConds get query conditions
func (r *Relation) GetConds() []Expr { return r.conds }

// GetSelects get select columns
func (r *Relation) GetSelects() []Expr { return r.selects }

// GetOrderCol get order columns
func (r *Relation) GetOrderCol() []Expr { return r.order }

// GetClauses get clauses
func (r *Relation) GetClauses() []clause.Expression { return r.clauses }

// GetScopes get scope functions
func (r *Relation) GetScopes() []relationScope { return r.scopes } // nolint

// GetPage get offset and limit
func (r *Relation) GetPage() (offset, limit int) { return r.offset, r.limit }

// StructField return struct field code
func (r *Relation) StructField() (fieldStr string) {
	for _, relation := range r.childRelations {
		fieldStr += relation.fieldName + " struct {\nfield.RelationField\n" + relation.StructField() + "}\n"
	}
	return fieldStr
}

// StructFieldInit return field initialize code
func (r *Relation) StructFieldInit() string {
	initStr := fmt.Sprintf("RelationField: field.NewRelation(%q, %q),\n", r.fieldPath, r.fieldType)
	for _, relation := range r.childRelations {
		initStr += relation.fieldName + ": struct {\nfield.RelationField\n" + strings.TrimSpace(relation.StructField()) + "}"
		initStr += "{\n" + relation.StructFieldInit() + "},\n"
	}
	return initStr
}

func wrapPath(root string, rs []Relation) []Relation {
	result := make([]Relation, len(rs))
	for i, r := range rs {
		r.fieldPath = root + "." + r.fieldPath
		r.childRelations = wrapPath(root, r.childRelations)
		result[i] = r
	}
	return result
}

var defaultRelationshipPrefix = map[RelationshipType]string{
	// HasOne:    "",
	// BelongsTo: "",
	HasMany:   "[]",
	Many2Many: "[]",
}

// RelateConfig config for relationship
type RelateConfig struct {
	RelatePointer      bool
	RelateSlice        bool
	RelateSlicePointer bool

	JSONTag      string
	GORMTag      string
	NewTag       string
	OverwriteTag string
}

// RelateFieldPrefix return generated relation field's type
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
