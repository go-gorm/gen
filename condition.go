package gen

import (
	"fmt"

	"gorm.io/datatypes"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
)

func Cond(exprs ...clause.Expression) []Condition {
	return exprToCondition(exprs...)
}

var _ Condition = &cond{}

type cond struct{ value interface{} }

func (c *cond) BeCond() interface{} { return c.value }

func exprToCondition(exprs ...clause.Expression) []Condition {
	conds := make([]Condition, len(exprs))
	for i, e := range exprs {
		if c, ok := e.(Condition); ok {
			conds[i] = c
		} else {
			conds[i] = &cond{value: e}
		}
	}
	return conds
}

func condToExpression(conds []Condition) ([]clause.Expression, error) {
	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		switch cond := cond.BeCond().(type) {
		case subQuery:
			exprs = append(exprs, cond.underlyingDO().buildCondition()...)
		case field.Expr:
			if expr, ok := cond.RawExpr().(clause.Expression); ok {
				exprs = append(exprs, expr)
			}
		case *datatypes.JSONQueryExpression:
			exprs = append(exprs, cond)
		default:
			return nil, fmt.Errorf("unsupport Condition %T", cond)
		}
	}
	return exprs, nil
}
