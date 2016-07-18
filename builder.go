package query

import (
	"reflect"
	"sync"
)

type (
	Builder interface {
		Select(columns ...string) Builder
		Count(column string) Builder
		Save(entities interface{}) Builder
		Delete() Builder
		Distinct(bool) Builder
		From(from string) Builder
		Where(where string) Builder
		AndWhere(where string) Builder
		OrWhere(where string) Builder
		Having(having string) Builder
		AndHaving(having string) Builder
		OrHaving(having string) Builder
		OrderBy(column, order string) Builder
		AddOrderBy(column, order string) Builder
		GroupBy(columns ...string) Builder
		AddGroupBy(columns ...string) Builder
		Limit(i int64) Builder
		Offset(i int64) Builder
		SetParameters(parameters ...interface{}) Builder
		AddParameters(parameters ...interface{}) Builder
		Reset() Builder
		GetQuery() Query
	}

	builder struct {
		t          reflect.Type
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

		orderMu sync.RWMutex
	}

	factory func(t reflect.Type) Builder
)

var New factory

func (b *builder) Select(columns ...string) Builder {
	b.statement = "select"
	b.columns = columns

	return b
}

func (b *builder) Count(column string) Builder {
	if column == "" {
		column = "*"
	}
	b.statement = "count"
	b.column = column

	return b
}

func (b *builder) Distinct(distinct bool) Builder {
	b.distinct = distinct

	return b
}

func (b *builder) Save(entities interface{}) Builder {
	b.statement = "save"
	b.v = entities

	return b
}

func (b *builder) Delete() Builder {
	b.statement = "delete"

	return b
}

func (b *builder) From(from string) Builder {
	b.from = from

	return b
}

func (b *builder) Where(where string) Builder {
	b.andWhere = b.andWhere[:0]
	b.orWhere = b.orWhere[:0]
	b.where = where

	return b
}

func (b *builder) AndWhere(where string) Builder {
	b.andWhere = append(b.andWhere, where)

	return b
}

func (b *builder) OrWhere(where string) Builder {
	b.orWhere = append(b.orWhere, where)

	return b
}

func (b *builder) Having(having string) Builder {
	b.andHaving = b.andHaving[:0]
	b.orHaving = b.orHaving[:0]
	b.having = having

	return b
}

func (b *builder) AndHaving(having string) Builder {
	b.andHaving = append(b.andHaving, having)

	return b
}

func (b *builder) OrHaving(having string) Builder {
	b.orHaving = append(b.orHaving, having)

	return b
}

func (b *builder) OrderBy(column, order string) Builder {
	b.orderMu.Lock()
	defer b.orderMu.Unlock()
	b.order = make(map[string]string)
	b.order[column] = order

	return b
}

func (b *builder) AddOrderBy(column, order string) Builder {
	b.orderMu.Lock()
	defer b.orderMu.Unlock()
	b.order[column] = order

	return b
}

func (b *builder) GroupBy(columns ...string) Builder {
	b.groupby = columns

	return b
}

func (b *builder) AddGroupBy(columns ...string) Builder {
	b.groupby = append(b.groupby, columns...)

	return b
}

func (b *builder) Limit(i int64) Builder {
	b.limit = i

	return b
}

func (b *builder) Offset(i int64) Builder {
	b.offset = i

	return b
}

func (b *builder) SetParameters(parameters ...interface{}) Builder {
	b.parameters = parameters

	return b
}

func (b *builder) AddParameters(parameters ...interface{}) Builder {
	for _, param := range parameters {
		b.parameters = append(b.parameters, param)
	}

	return b
}

func (b *builder) Reset() Builder {
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
	b.orderMu.Lock()
	defer b.orderMu.Unlock()

	b.v = v
	b.statement = statement
	b.columns = columns
	b.column = column
	b.distinct = distinct
	b.from = from
	b.where = where
	b.orWhere = orWhere
	b.andWhere = andWhere
	b.having = having
	b.orHaving = orHaving
	b.andHaving = andHaving
	b.order = order
	b.groupby = groupby
	b.limit = limit
	b.offset = offset
	b.parameters = parameters

	return b
}

func (b *builder) GetQuery() Query {
	return &query{b}
}

func newBuilder(t reflect.Type) Builder {
	return &builder{
		t:     t,
		order: make(map[string]string),
	}
}

func init() {
	New = newBuilder
}
