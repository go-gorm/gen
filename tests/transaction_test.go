package tests_test

import (
	"fmt"
	"testing"

	"gorm.io/gen/tests/.expect/dal_test_relation/model"
	"gorm.io/gen/tests/.expect/dal_test_relation/query"
)

func TestQuery_Transaction_Relation(t *testing.T) {
	useOnce.Do(CRUDInit)

	t.Run("transaction has many", func(t *testing.T) {
		if err := query.Q.Transaction(func(tx *query.Query) error {
			c := tx.Customer
			customer := &model.Customer{
				Bank: model.Bank{
					Name:    "bank1",
					Address: "bank-address1",
					Scale:   1,
				},
				CreditCards: []model.CreditCard{
					{Number: "num1"},
					{Number: "num2"},
				},
			}
			if err := c.WithContext(ctx).Create(customer); err != nil {
				return fmt.Errorf("create model fail: %s", err)
			}

			got, err := c.WithContext(ctx).Where(c.ID.Eq(customer.ID)).
				Preload(c.CreditCards).
				Preload(c.Bank).
				First()
			if err != nil {
				return fmt.Errorf("find model fail: %s", err)
			}
			if len(got.CreditCards) != 2 {
				return fmt.Errorf("replace model fail, expect %d, got %d", 1, len(got.CreditCards))
			}

			if err := c.CreditCards.WithContext(ctx).Model(customer).Replace(&model.CreditCard{
				Number: "num_replace",
			}); err != nil {
				return fmt.Errorf("replace model fail: %s", err)
			}

			got, err = c.WithContext(ctx).Where(c.ID.Eq(customer.ID)).
				Preload(c.CreditCards).
				Preload(c.Bank).
				First()
			if err != nil {
				return fmt.Errorf("find model fail: %s", err)
			}
			if len(got.CreditCards) != 1 {
				return fmt.Errorf("replace model fail, expect %d, got %d", 1, len(got.CreditCards))
			}
			if got.CreditCards[0].Number != "num_replace" {
				return fmt.Errorf("replace model fail, expect %q, got %q", "num_replace", got.CreditCards[0].Number)
			}

			return nil
		}); err != nil {
			t.Errorf("transaction execute fail: %s", err)
		}
	})

	t.Run("transaction has one", func(t *testing.T) {
		if err := query.Q.Transaction(func(tx *query.Query) error {
			c := tx.Customer
			customer := &model.Customer{
				Bank: model.Bank{
					Name:    "bank1",
					Address: "bank-address1",
					Scale:   1,
				},
				CreditCards: []model.CreditCard{
					{Number: "num1"},
					{Number: "num2"},
				},
			}

			if err := c.WithContext(ctx).Create(customer); err != nil {
				return fmt.Errorf("create model fail: %s", err)
			}
			if err := c.Bank.WithContext(ctx).Model(customer).Replace(&model.Bank{
				Name:    "bank-replace",
				Address: "bank-replace-address",
				Scale:   2,
			}); err != nil {
				return fmt.Errorf("replace model fail: %s", err)
			}

			got, err := c.WithContext(ctx).Where(c.ID.Eq(customer.ID)).
				Preload(c.CreditCards).
				Preload(c.Bank).
				First()
			if err != nil {
				return fmt.Errorf("find model fail: %s", err)
			}
			if got.Bank.Name != "bank-replace" {
				return fmt.Errorf("replace model fail, expect %q, got %q", "bank-replace", got.Bank.Name)
			}

			return nil
		}); err != nil {
			t.Errorf("transaction execute fail: %s", err)
		}
	})
}
