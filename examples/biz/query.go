package biz

import (
	"context"
	"fmt"

	"gorm.io/gen/examples/dal/query"
)

var q = query.Q

func Query(ctx context.Context) {
	t := q.Mytable
	do := t.WithContext(context.Background())

	data, err := do.Take()
	catchError("Take", err)
	fmt.Printf("got %+v\n", data)

	dataArray, err := do.Find()
	catchError("Find", err)
	fmt.Printf("got %+v\n", dataArray)

	data, err = do.Where(t.ID.Eq(1)).Take()
	catchError("Take", err)
	fmt.Printf("got %+v\n", data)

	dataArray, err = do.Where(t.Age.Gt(18)).Order(t.Username).Find()
	catchError("Find", err)
	fmt.Printf("got %+v\n", dataArray)

	dataArray, err = do.Select(t.ID, t.Username).Order(t.Age.Desc()).Find()
	catchError("Find", err)
	fmt.Printf("got %+v\n", dataArray)

	info, err := do.Where(t.ID.Eq(1)).UpdateSimple(t.Age.Add(1))
	catchError("Update", err)
	fmt.Printf("got %+v\n", info)
}

func catchError(detail string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", detail, err)
	}
}
