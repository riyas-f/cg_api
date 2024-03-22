package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type FilterCondition struct {
	ColumnSource   string
	ColumnReceiver string
}

type FilteredQueryExecutor struct {
	useTransaction bool
	explicitCast   bool
	filterMap      map[string]FilterCondition
	DriverName     string
	Tx             *sql.Tx
}

func (f *FilteredQueryExecutor) AddTableSource(tableSource string, columnSource string, columnReceiver string) {
	f.filterMap[tableSource] = FilterCondition{
		ColumnSource:   columnSource,
		ColumnReceiver: columnReceiver,
	}
}

func (f *FilteredQueryExecutor) Insert() error {
	return nil
}

// Add explicit cast on the values placeholder on the query
// to make sure the data is being correctly interpreted
// ie. integer being interpreted as INTEGER not text
//
// string type value don't need explicit cast
func (f *FilteredQueryExecutor) UseExplicitCast() {
	f.explicitCast = true
}

// This will make subsequent query use transaction
// the sql.Tx struct used on the transaction will be returned
// after any operation that use transaction
// such as Insert, Update , and Delete
func (f *FilteredQueryExecutor) UseTransaction(db *sql.DB) error {
	f.useTransaction = true
	tx, err := db.Begin()

	if err != nil {
		return err
	}

	f.Tx = tx

	return nil
}

func (f *FilteredQueryExecutor) Rollback() error {
	if f.Tx == nil {
		return fmt.Errorf("illegal operations. trying to rollback from non transactional query")
	}
	return f.Tx.Rollback()
}

func (f *FilteredQueryExecutor) Commit() error {
	if f.Tx == nil {
		return fmt.Errorf("illegal operations. trying to commit from non transactional query")
	}
	return f.Tx.Commit()
}

// Insert multiple data in one query. Make sure data is a slice containing pointer of a struct of the same type
// no check is done, so multiple types of data in a single insert
// may lead to a panic
//
// When using transaction, *sql.Tx
// struct will be returned when successfully begin a transaction, else nil will be returned
func (f *FilteredQueryExecutor) BatchInsert(data interface{}, db *sql.DB, tableName string) error {
	var valueArgs []interface{}

	var nCol int
	// check data is an interface

	v := ""
	r := ""
	k := ""

	val := reflect.ValueOf(data)

	if val.Kind() != reflect.Slice {
		return fmt.Errorf("illegal arguments. data must be a slice")
	}

	for i := 0; i < val.Len(); i++ {
		obj := val.Index(i).Interface()

		column, value, _ := getField(obj, false)

		if k == "" {
			k = strings.Join(column, ",")
		}

		if r == "" {
			nCol = len(column)
			for _, s := range column {
				r += fmt.Sprintf("d.%s,", s)
			}
		}

		if valueArgs == nil {
			valueArgs = make([]interface{}, 0, val.Len()*len(value))
		}

		valueArgs = append(valueArgs, value...)
	}

	r = r[:len(r)-len(",")] // select string

	for i := 0; i < len(valueArgs); i++ {
		v += "("
		for j := 0; j < nCol; j++ {
			v += fmt.Sprintf("$%d", i+1)
			// add cast on the query for datatype other than string
			if f.explicitCast {
				v += addCast(reflect.TypeOf(valueArgs[i]).Kind())
			}
			v += ","
		}
		v = v[:len(v)-1]
		v += "),"
	}

	v = v[:len(v)-1]

	whereExistBlock := makeWhereExistClause(f.filterMap)

	query := fmt.Sprintf(`WITH data(%s) AS(
		VALUES %s
	) INSERT INTO %s (%s) 
	SELECT %s
	FROM data d
	%s`, k, v, tableName, k, r, whereExistBlock)

	fmt.Println(query)

	if f.useTransaction {
		_, err := f.Tx.Exec(query, valueArgs...)

		return err
	}

	_, err := db.Exec(query, valueArgs...)

	return err
}

func (f *FilteredQueryExecutor) Update() error {
	return nil
}

func (f *FilteredQueryExecutor) Delete() error {
	return nil
}

func makeWhereExistClause(condition map[string]FilterCondition) string {
	t := ""
	w := ""

	for key, element := range condition {
		t += fmt.Sprintf("%s,", key)
		w += fmt.Sprintf("%s.%s = d.%s AND ", key, element.ColumnSource, element.ColumnReceiver)
	}

	// remove separator
	t = t[:len(t)-len(",")]
	w = w[:len(w)-len(" AND ")]

	block := fmt.Sprintf(
		`WHERE EXISTS (
				SELECT 1
				FROM %s
				WHERE %s
			)`, t, w,
	)

	return block
}

func addCast(dataType reflect.Kind) string {
	// https://zontroy.com/postgresql-to-go-type-mapping

	switch dataType {
	case reflect.Int:
		return "::INTEGER"

	case reflect.Int64:
		return "::BIGINT"

	case reflect.Float32:
		return "::DOUBLE PRECISION"

	case reflect.Float64:
		return "::NUMERIC"
	}

	return ""
}
