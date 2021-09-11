package gen

import "errors"

var (
	// ErrInvalidExpression invalid Expression
	ErrInvalidExpression = errors.New("invalid expression")

	// ErrEmptyCondition empty condition
	ErrEmptyCondition = errors.New("empty condition")
)
