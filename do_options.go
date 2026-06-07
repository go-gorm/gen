package gen

import (
	"errors"

	"gorm.io/gorm/clause"
)

// DOOption gorm option interface
type DOOption interface {
	Apply(*DOConfig) error
	AfterInitialize(*DO) error
}

type ClauseChecker func(clause.Expression) error

var ErrClauseNotHandled = errors.New("clause not handled")

type DOConfig struct {
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

func WithClauseChecker(checker ClauseChecker) DOOption {
	return clauseCheckerOption{checker: checker}
}
