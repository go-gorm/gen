package gen

import (
	"fmt"

	"gorm.io/datatypes"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
)

// Cond convert expression array to condition array
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
		case *datatypes.JSONQueryExpression, *datatypes.JSONOverlapsExpression:
			conds = append(conds, &condContainer{value: e})
		default:
			conds = append(conds, &condContainer{err: fmt.Errorf("unsupported Expression %T to converted to Condition", e)})
		}
	}
	return conds
}

func condToExpression(conds []Condition) ([]clause.Expression, error) {
	if len(conds) == 0 {
		return nil, nil
	}
	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		if cond == nil {
			continue
		}
		if err := cond.CondError(); err != nil {
			return nil, err
		}

		switch cond.(type) {
		case *condContainer, field.Expr, SubQuery:
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
