package goquery

import (
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/gedex/inflector"
)

type (
	Query interface {
		GetResults(*sql.DB) (interface{}, error)
		GetResult(*sql.DB) (interface{}, error)
		GetCount(*sql.DB) (int64, error)
		Execute(*sql.DB) (interface{}, error)
		GetSQL() string
	}

	query struct {
		builder *builder
	}
)

func (q *query) GetResults(db *sql.DB) (interface{}, error) {
	slice := reflect.New(reflect.SliceOf(q.builder.t)).Elem()
	ptr := reflect.New(q.builder.t)
	entity := ptr.Elem()
	fieldInfo, queryStr := prepareSelect(entity, q.builder)

	stmt, err := db.Prepare(queryStr)
	if err != nil {
		return slice.Interface(), err
	}
	defer stmt.Close()

	rows, err := stmt.Query(q.builder.parameters...)
	if err != nil {
		return slice.Interface(), err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(fieldInfo...); err != nil {
			return slice.Interface(), err
		}
		slice.Set(reflect.Append(slice, entity))
	}

	if err = rows.Err(); err != nil {
		return slice.Interface(), err
	}

	return slice.Interface(), nil
}

func (q *query) GetResult(db *sql.DB) (interface{}, error) {
	ptr := reflect.New(q.builder.t)
	entity := ptr.Elem()
	fieldInfo, queryStr := prepareSelect(entity, q.builder)
	stmt, err := db.Prepare(queryStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if err := stmt.QueryRow(q.builder.parameters...).Scan(fieldInfo...); err != nil {
		return nil, err
	}

	return entity.Interface(), nil
}

func (q *query) GetCount(db *sql.DB) (int64, error) {
	var count int64

	stmt, err := db.Prepare(q.GetSQL())
	if err != nil {
		return count, err
	}
	defer stmt.Close()

	if err := stmt.QueryRow(q.builder.parameters...).Scan(count); err != nil {
		return count, err
	}

	return count, nil
}

func (q *query) Execute(db *sql.DB) (interface{}, error) {
	switch q.builder.statement {
	case "save":
		return save(q, db)
	case "delete":
		return remove(q, db)
	default:
		return nil, errors.New("query: Invalid execute statement!")
	}
}

func (q *query) GetSQL() string {
	var (
		cols      = make(map[string]string)
		colsCount int
		queryStr  string
	)

	table := getTable(q.builder.t)

	switch q.builder.statement {
	case "select":
		for _, col := range q.builder.columns {
			cols[col] = col
			colsCount++
		}
		for i := 0; i < q.builder.t.NumField(); i++ {
			col := getColumn(q.builder.t.Field(i))
			if _, ok := cols[col]; col != "" && (ok || colsCount == 0) {
				if queryStr == "" {
					queryStr += "SELECT "
					if q.builder.distinct {
						queryStr += "DISTINCT "
					}
					queryStr += col
				} else {
					queryStr += ", " + col
				}
			}
		}
		break
	case "count":
		queryStr = "SELECT COUNT("
		if q.builder.distinct {
			queryStr += "DISTINCT "
		}
		queryStr += q.builder.column + ")"
		break
	case "insert":
		var values string
		for i := 1; i < q.builder.t.NumField(); i++ {
			col := getColumn(q.builder.t.Field(i))
			if _, ok := cols[col]; col != "" && (ok || colsCount == 0) {
				if values == "" {
					queryStr += "INSERT INTO " + table + " (" + col
					values = "?"
				} else {
					queryStr += ", " + col
					values += ", ?"
				}
			}
		}
		queryStr += ") VALUES (" + values + ")"
		break
	case "update":
		for i := 1; i < q.builder.t.NumField(); i++ {
			col := getColumn(q.builder.t.Field(i))
			if _, ok := cols[col]; col != "" && (ok || colsCount == 0) {
				if queryStr == "" {
					queryStr += "UPDATE  " + table + " SET " + col + "=?"
				} else {
					queryStr += ", " + col + "=?"
				}
			}
		}
		queryStr += " WHERE id=?"
		if q.builder.where != "" {
			queryStr += " AND "
		}
		break
	case "delete":
		queryStr = "DELETE FROM " + table
		break
	}

	return finishSQL(q.builder, queryStr, table)
}

func finishSQL(builder *builder, queryStr, table string) string {
	var from string
	switch builder.statement {
	case "select", "count":
		if builder.from != "" {
			from = builder.from
		} else {
			from = table
		}
		queryStr += " FROM " + from
	}

	setWhere := false
	if builder.where != "" {
		setWhere = true
		queryStr += " WHERE " + builder.where
	}
	for i := 0; i < len(builder.andWhere); i++ {
		if i == 0 && !setWhere {
			setWhere = true
			queryStr += " WHERE " + builder.andWhere[i]
		} else {
			queryStr += " AND " + builder.andWhere[i]
		}
	}
	for i := 0; i < len(builder.orWhere); i++ {
		if i == 0 && !setWhere {
			setWhere = true
			queryStr += " WHERE " + builder.orWhere[i]
		} else {
			queryStr += " OR " + builder.orWhere[i]
		}
	}

	orderby := ""
	for col, order := range builder.order {
		if orderby == "" {
			orderby += " ORDER BY " + col + " " + order
		} else {
			orderby += ", " + col + " " + order
		}
	}
	queryStr += orderby

	groupby := ""
	for _, col := range builder.groupby {
		if groupby == "" {
			groupby += " GROUP BY " + col
		} else {
			groupby += ", " + col
		}
	}
	queryStr += groupby

	setHaving := false
	if builder.having != "" {
		setHaving = true
		queryStr += " HAVING " + builder.having
	}
	for i := 0; i < len(builder.andHaving); i++ {
		if i == 0 && !setHaving {
			setHaving = true
			queryStr += " HAVING " + builder.andHaving[i]
		} else {
			queryStr += " AND " + builder.andHaving[i]
		}
	}
	for i := 0; i < len(builder.orHaving); i++ {
		if i == 0 && !setHaving {
			setHaving = true
			queryStr += " HAVING " + builder.orHaving[i]
		} else {
			queryStr += " OR " + builder.orHaving[i]
		}
	}

	if builder.limit > 0 {
		queryStr += " LIMIT " + strconv.FormatInt(builder.limit, 10)
	}

	if builder.offset > 0 {
		queryStr += " OFFSET " + strconv.FormatInt(builder.offset, 10)
	}

	return queryStr
}

func getTable(t reflect.Type) string {
	return inflector.Pluralize(strings.ToLower(t.Name()))
}

func getField(f reflect.Value) interface{} {
	return f.Addr().Interface()
}

func getColumn(f reflect.StructField) string {
	return f.Tag.Get("column")
}

func getColumns(t reflect.Type) []string {
	cols := []string{}
	for i := 0; i < t.NumField(); i++ {
		col := getColumn(t.Field(i))
		if col != "" {
			cols = append(cols, col)
		}
	}

	return cols
}

func getFields(s reflect.Value) []interface{} {
	fields := []interface{}{}
	t := s.Type()
	for i := 0; i < s.NumField(); i++ {
		col := getColumn(t.Field(i))
		if col != "" {
			fields = append(fields, getField(s.Field(i)))
		}
	}

	return fields
}

func prepareSelect(s reflect.Value, builder *builder) ([]interface{}, string) {
	var (
		cols      = make(map[string]string)
		fieldInfo []interface{}
		queryStr  string
		colsCount int
	)

	for _, col := range builder.columns {
		cols[col] = col
		colsCount++
	}

	t := s.Type()
	for i := 0; i < s.NumField(); i++ {
		col := getColumn(t.Field(i))
		if _, ok := cols[col]; col != "" && (ok || colsCount == 0) {
			fieldInfo = append(fieldInfo, getField(s.Field(i)))
			if queryStr == "" {
				queryStr += "SELECT "
				if builder.distinct {
					queryStr += "DISTINCT "
				}
				queryStr += col
			} else {
				queryStr += ", " + col
			}
		}
	}

	return fieldInfo, finishSQL(builder, queryStr, getTable(t))
}

func save(q *query, db *sql.DB) (interface{}, error) {
	slice := reflect.ValueOf(q.builder.v)
	if slice.Kind() == reflect.Ptr {
		slice = slice.Elem()
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	q.builder.statement = "update"
	stmtU, err := tx.Prepare(q.GetSQL())
	if err != nil {
		return nil, err
	}
	defer stmtU.Close()

	q.builder.statement = "insert"
	stmtA, err := tx.Prepare(q.GetSQL())
	if err != nil {
		return nil, err
	}
	defer stmtA.Close()

	for i := 0; i < slice.Len(); i++ {
		var fieldInfo []interface{}
		s := slice.Index(i)
		for j := 1; j < s.NumField(); j++ {
			fieldInfo = append(fieldInfo, s.Field(j).Interface())
		}
		isNew := s.Field(0).Int() == int64(0)
		if !isNew {
			fieldInfo = append(fieldInfo, s.Field(0).Interface())
		}

		if isNew {
			res, err := stmtA.Exec(fieldInfo...)
			if err != nil {
				return slice.Interface(), err
			}

			id, err := res.LastInsertId()
			if err != nil {
				return slice.Interface(), err
			}
			s.FieldByName("Id").SetInt(id)
		} else {
			_, err := stmtU.Exec(fieldInfo...)
			if err != nil {
				return slice.Interface(), err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return slice.Interface(), err
	}

	return slice.Interface(), nil
}

func remove(q *query, db *sql.DB) (interface{}, error) {
	_, err := db.Exec(q.GetSQL(), q.builder.parameters...)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
