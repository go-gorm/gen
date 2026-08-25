package gen

import (
	"errors"

	"gorm.io/gorm/clause"
)

// DOOption gorm option interface
type DOOption interface {
	// Apply updates the configuration before a DO is initialized.
	Apply(*DOConfig) error
	// AfterInitialize defines an optional post-initialization hook for a DO.
	AfterInitialize(*DO) error
}

// ClauseChecker validates a clause expression before it reaches the built-in
// checker. Return ErrClauseNotHandled to delegate validation to CheckClause.
type ClauseChecker func(clause.Expression) error

// ErrClauseNotHandled tells the clause validation pipeline to fall back to the
// built-in CheckClause implementation.
var ErrClauseNotHandled = errors.New("clause not handled")

// DOConfig contains runtime options shared by generated data objects.
type DOConfig struct {
	// ClauseChecker optionally validates clauses before the built-in checker.
	ClauseChecker ClauseChecker
}

// Apply update config to new config
func (c *DOConfig) Apply(config *DOConfig) error {
	if config != c {
		*config = *c
	}
	return nil
}

// AfterInitialize initialize plugins after db connected
func (c *DOConfig) AfterInitialize(db *DO) error {
	return nil
}

type clauseCheckerOption struct {
	checker ClauseChecker
}

func (o clauseCheckerOption) Apply(cfg *DOConfig) error {
	cfg.ClauseChecker = o.checker
	return nil
}

func (clauseCheckerOption) AfterInitialize(*DO) error { return nil }

// WithClauseChecker installs checker for generated data objects.
// A nil result accepts a clause, ErrClauseNotHandled delegates to CheckClause,
// and any other error rejects the clause.
func WithClauseChecker(checker ClauseChecker) DOOption {
	return clauseCheckerOption{checker: checker}
}
