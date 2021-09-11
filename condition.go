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

var _ Condition = &condContainer{}

type condContainer struct {
	value interface{}
	err   error
}

func (c *condContainer) BeCond() interface{} { return c.value }
func (c *condContainer) CondError() error    { return c.err }

func exprToCondition(exprs ...clause.Expression) []Condition {
	conds := make([]Condition, 0, len(exprs))
	for _, e := range exprs {
		switch e := e.(type) {
		case *datatypes.JSONQueryExpression:
			conds = append(conds, &condContainer{value: e})
		default:
			conds = append(conds, &condContainer{err: fmt.Errorf("unsupported Expression %T to converted to Condition", e)})
		}
	}
	return conds
}

func condToExpression(conds []Condition) ([]clause.Expression, error) {
	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		if err := cond.CondError(); err != nil {
			return nil, err
		}

		switch cond.(type) {
		case *condContainer, field.Expr, subQuery:
		default:
			return nil, fmt.Errorf("unsupported condition: %+v", cond)
		}

		switch e := cond.BeCond().(type) {
		case []clause.Expression:
			exprs = append(exprs, e...)
		case clause.Expression:
			exprs = append(exprs, e)
		}
	}
	return exprs, nil
}
