package goquery

import (
	"fmt"
	"reflect"
)

func ExampleParseQuery() {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	sql := qBuilder.
		Select("id", "col1", "email", "col2", "col3").
		Where("id = ?").
		AndWhere("col2 = ?").
		OrWhere("col3 = ?").
		Having("COUNT(col1) > ?").
		AndHaving("COUNT(col2) > ?").
		OrHaving("COUNT(col3) > ?").
		OrderBy("id", "DESC").
		GroupBy("id").
		AddGroupBy("email").
		Limit(10).
		Offset(5).
		GetQuery().
		GetSQL()

	fmt.Println(sql)
	// Output: SELECT id, email FROM users WHERE id = ? AND col2 = ? OR col3 = ? ORDER BY id DESC GROUP BY id, email HAVING COUNT(col1) > ? AND COUNT(col2) > ? OR COUNT(col3) > ? LIMIT 10 OFFSET 5
}

func ExampleParseEmptySelect() {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	sql := qBuilder.
		Select().
		GetQuery().
		GetSQL()

	fmt.Println(sql)
	// Output: SELECT id, email FROM users
}

func ExampleParseFrom() {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	sql := qBuilder.
		Select().
		From("foo").
		GetQuery().
		GetSQL()

	fmt.Println(sql)
	// Output: SELECT id, email FROM foo
}

func ExampleParseCount() {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	sql := qBuilder.
		Count("id").
		GetQuery().
		GetSQL()

	fmt.Println(sql)
	// Output: SELECT COUNT(id) FROM users
}

func ExampleParseOrder() {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	sql := qBuilder.
		Select("id").
		OrderBy("id", "DESC").
		GetQuery().
		GetSQL()

	fmt.Println(sql)
	// Output: SELECT id FROM users ORDER BY id DESC
}

func ExampleParseAddOrder() {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	sql := qBuilder.
		Select("id").
		AddOrderBy("email", "ASC").
		GetQuery().
		GetSQL()

	fmt.Println(sql)
	// Output: SELECT id FROM users ORDER BY email ASC
}

func ExampleParseDistinct() {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	sql := qBuilder.
		Select("id").
		Distinct(true).
		GetQuery().
		GetSQL()

	fmt.Println(sql)
	// Output: SELECT DISTINCT id FROM users
}

func ExampleParseDistinctCount() {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	sql := qBuilder.
		Count("id").
		Distinct(true).
		GetQuery().
		GetSQL()

	fmt.Println(sql)
	// Output: SELECT COUNT(DISTINCT id) FROM users
}

func ExampleParseDelete() {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	sql := qBuilder.
		Delete().
		Where("id = ?").
		GetQuery().
		GetSQL()

	fmt.Println(sql)
	// Output: DELETE FROM users WHERE id = ?
}
