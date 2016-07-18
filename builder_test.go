package query

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	var columns []string

	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	qBuilder.Select("col1", "col2", "col3")
	assert.Equal(t, []string{"col1", "col2", "col3"}, qBuilder.(*builder).columns)
	assert.Equal(t, "select", qBuilder.(*builder).statement)

	qBuilder.Select("col4", "col5", "col6")
	assert.Equal(t, []string{"col4", "col5", "col6"}, qBuilder.(*builder).columns)

	qBuilder.Select()
	assert.Equal(t, columns, qBuilder.(*builder).columns)
}

func TestCount(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	qBuilder.Count("col")
	assert.Equal(t, "col", qBuilder.(*builder).column)
	assert.Equal(t, "count", qBuilder.(*builder).statement)

	qBuilder.Count("col2")
	assert.Equal(t, "col2", qBuilder.(*builder).column)

	qBuilder.Count("")
	assert.Equal(t, "*", qBuilder.(*builder).column)

}

func TestDistinct(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	qBuilder.Distinct(true)
	assert.Equal(t, true, qBuilder.(*builder).distinct)
}

func TestSave(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	users := []*user{&user{}, &user{}}
	qBuilder.Save(users)
	assert.Equal(t, users, qBuilder.(*builder).v)
	assert.Equal(t, "save", qBuilder.(*builder).statement)
}

func TestDelete(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	qBuilder.Delete()
	assert.Equal(t, "delete", qBuilder.(*builder).statement)
}

func TestFrom(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	qBuilder.From("foo")
	assert.Equal(t, "foo", qBuilder.(*builder).from)

	qBuilder.From("foo2")
	assert.Equal(t, "foo2", qBuilder.(*builder).from)
}

func TestLimit(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	qBuilder.Limit(1)
	assert.Equal(t, int64(1), qBuilder.(*builder).limit)

	qBuilder.Limit(2)
	assert.Equal(t, int64(2), qBuilder.(*builder).limit)
}

func TestOffset(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	qBuilder.Offset(1)
	assert.Equal(t, int64(1), qBuilder.(*builder).offset)

	qBuilder.Offset(2)
	assert.Equal(t, int64(2), qBuilder.(*builder).offset)
}

func TestWhere(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	var (
		orWhere  []string
		andWhere []string
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	where := "col = ?"
	qBuilder.Where(where)
	assert.Equal(t, where, qBuilder.(*builder).where)
	assert.Equal(t, andWhere, qBuilder.(*builder).andWhere)
	assert.Equal(t, orWhere, qBuilder.(*builder).orWhere)

	where = "col2 = ?"
	qBuilder.Where(where)
	assert.Equal(t, where, qBuilder.(*builder).where)
}

func TestAndWhere(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	andWhere := "col = ?"
	qBuilder.AndWhere(andWhere)
	assert.Equal(t, andWhere, qBuilder.(*builder).andWhere[0])

	andWhere = "col2 = ?"
	qBuilder.AndWhere(andWhere)
	assert.Equal(t, andWhere, qBuilder.(*builder).andWhere[1])
}

func TestOrWhere(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	orWhere := "col = ?"
	qBuilder.OrWhere(orWhere)
	assert.Equal(t, orWhere, qBuilder.(*builder).orWhere[0])

	orWhere = "col2 = ?"
	qBuilder.OrWhere(orWhere)
	assert.Equal(t, orWhere, qBuilder.(*builder).orWhere[1])
}

func TestHaving(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	var (
		orHaving  []string
		andHaving []string
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	having := "COUNT(col) > ?"
	qBuilder.Having(having)
	assert.Equal(t, having, qBuilder.(*builder).having)
	assert.Equal(t, andHaving, qBuilder.(*builder).andHaving)
	assert.Equal(t, orHaving, qBuilder.(*builder).orHaving)

	having = "COUNT(col2) > ?"
	qBuilder.Having(having)
	assert.Equal(t, having, qBuilder.(*builder).having)
}

func TestAndHaving(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	having := "COUNT(col) > ?"
	qBuilder.AndHaving(having)
	assert.Equal(t, having, qBuilder.(*builder).andHaving[0])

	having = "COUNT(col2) > ?"
	qBuilder.AndHaving(having)
	assert.Equal(t, having, qBuilder.(*builder).andHaving[1])
}

func TestOrHaving(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	having := "COUNT(col) > ?"
	qBuilder.OrHaving(having)
	assert.Equal(t, having, qBuilder.(*builder).orHaving[0])

	having = "COUNT(col2) > ?"
	qBuilder.OrHaving(having)
	assert.Equal(t, having, qBuilder.(*builder).orHaving[1])
}

func TestOrderBy(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	qBuilder.OrderBy("col", "ASC")

	qBuilder.(*builder).orderMu.Lock()
	defer qBuilder.(*builder).orderMu.Unlock()
	order := qBuilder.(*builder).order["col"]

	assert.Equal(t, "ASC", order)
}

func TestAddOrderBy(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	qBuilder.AddOrderBy("col", "DESC")

	qBuilder.(*builder).orderMu.Lock()
	defer qBuilder.(*builder).orderMu.Unlock()
	order := qBuilder.(*builder).order["col"]

	assert.Equal(t, "DESC", order)
}

func TestGroupBy(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	qBuilder.GroupBy("col", "col2")
	assert.Equal(t, []string{"col", "col2"}, qBuilder.(*builder).groupby)
}

func TestAddGroupBy(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	qBuilder.AddGroupBy("col", "col2")
	assert.Equal(t, []string{"col", "col2"}, qBuilder.(*builder).groupby)
}

func TestSetParameters(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	var params []interface{}
	params = append(params, 1)
	params = append(params, "string")
	qBuilder.SetParameters(params...)

	assert.Equal(t, params, qBuilder.(*builder).parameters)
}

func TestAddParameters(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)
	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	var params []interface{}
	params = append(params, 1)
	params = append(params, "string")
	qBuilder.AddParameters(params...)

	assert.Equal(t, params, qBuilder.(*builder).parameters)
}

func TestReset(t *testing.T) {
	type (
		user struct {
			Id    int64  `json:"id" column:"id"`
			Email string `json:"email" column:"email"`
		}
	)

	var (
		v          interface{}
		statement  string
		columns    []string
		column     string
		distinct   bool
		from       string
		where      string
		orWhere    []string
		andWhere   []string
		having     string
		orHaving   []string
		andHaving  []string
		order      map[string]string
		groupby    []string
		limit      int64
		offset     int64
		parameters []interface{}
	)

	reflectT := reflect.TypeOf(user{})
	qBuilder := New(reflectT)

	qBuilder.
		Select("col1", "col2", "col3").
		Count("col4").
		Save([]*user{&user{}, &user{}}).
		Distinct(true).
		From("foo").
		Where("col1 = ?").
		AndWhere("col2 = ?").
		OrWhere("col3 = ?").
		Having("HAVING COUNT(col1) > ?").
		AndHaving("HAVING COUNT(col2) > ?").
		OrHaving("HAVING COUNT(col3) > ?").
		OrderBy("col1", "ASC").
		AddOrderBy("col2", "DESC").
		GroupBy("col1", "col2").
		AddGroupBy("col3", "col4").
		Limit(60).
		Offset(2).
		SetParameters(1, 2, 3).
		AddParameters(6, 4, 5)

	qBuilder.Reset()

	assert.Equal(t, v, qBuilder.(*builder).v)
	assert.Equal(t, statement, qBuilder.(*builder).statement)
	assert.Equal(t, columns, qBuilder.(*builder).columns)
	assert.Equal(t, column, qBuilder.(*builder).column)
	assert.Equal(t, distinct, qBuilder.(*builder).distinct)
	assert.Equal(t, from, qBuilder.(*builder).from)
	assert.Equal(t, where, qBuilder.(*builder).where)
	assert.Equal(t, orWhere, qBuilder.(*builder).orWhere)
	assert.Equal(t, andWhere, qBuilder.(*builder).andWhere)
	assert.Equal(t, having, qBuilder.(*builder).having)
	assert.Equal(t, orHaving, qBuilder.(*builder).orHaving)
	assert.Equal(t, andHaving, qBuilder.(*builder).andHaving)
	assert.Equal(t, order, qBuilder.(*builder).order)
	assert.Equal(t, groupby, qBuilder.(*builder).groupby)
	assert.Equal(t, limit, qBuilder.(*builder).limit)
	assert.Equal(t, offset, qBuilder.(*builder).offset)
	assert.Equal(t, parameters, qBuilder.(*builder).parameters)
}
