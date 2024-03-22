package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

type MatchOperand string

const (
	GT  MatchOperand = ">"
	GEQ MatchOperand = ">="
	LT  MatchOperand = "<"
	LEQ MatchOperand = "<="
	EQ  MatchOperand = "="
)

type QueryCondition struct {
	TableName  string
	ColumnName string
	MatchValue any
	Operand    any
}

type QueryOperation interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type Querynator struct {
	DriverName string
	// useTransaction bool
}

// func (q *Querynator) UseTransaction() *Querynator {
// 	return &Querynator{
// 		useTransaction: true,
// 	}
// }

func (q *Querynator) PrepareFilterOperation() *FilteredQueryExecutor {
	return &FilteredQueryExecutor{
		useTransaction: false,
		filterMap:      make(map[string]FilterCondition),
		DriverName:     q.DriverName,
	}
}

func (q *Querynator) PrepareJoinOperation() *JoinQueryExecutor {
	return &JoinQueryExecutor{
		DriverName: q.DriverName,
	}
}

// func (q *Querynator) InserMultiple(data []interface{}, db *sql.DB, tableName string, columnNames []string) error {
// 	sqlDbx := sqlx.NewDb(db, q.DriverName)

// 	v := ""

// 	for _, name := range columnNames {
// 		v += fmt.Sprintf(":%s, ", name)
// 	}

// 	// Test named query
// 	/* INSERT INTO person (first_name, last_name, email)
// 	   VALUES (:first_name, :last_name, :email) */

// 	v = v[:len(v)-len(", ")]            // values string
// 	k := strings.Join(columnNames, ",") // column string

// 	query := fmt.Sprintf("INSERT INTO (%s) VALUES(%s)", k, v)

// 	if q.useTransaction {
// 		tx, err := sqlDbx.Beginx()
// 		if err != nil {
// 			return err
// 		}

// 		_, err = tx.NamedExec(query, data)

// 		return err
// 	}

// 	_, err := sqlDbx.NamedExec(query, data)

// 	return err

// }

func (q *Querynator) Insert(v interface{}, db QueryOperation, tableName string, returnField string) (interface{}, error) {
	var query string
	var id interface{}
	var err error

	// Insert stuff here
	fields, values, _ := getField(v, true)

	elements := strings.Join(fields, ", ")

	p := ""
	for i := 0; i < len(values); i++ {
		p += fmt.Sprintf("$%d,", i+1)
	}

	if returnField != "" {
		query = fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s) RETURNING %s", tableName, elements, p[:len(p)-1], returnField)
		err = db.QueryRow(query, values...).Scan(&id)

	} else {
		query = fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", tableName, elements, p[:len(p)-1])
		_, err = db.Exec(query, values...)
	}

	// query := fmt.Sprintf(
	// 	`INSERT INTO %s (username, name, email, password, is_active)
	// 	VALUES($1, $2, $3, $4, $5)`, tableName,
	// )

	if err != nil {
		return -1, err
	}

	return id, nil
}

func (q *Querynator) Delete(v interface{}, db QueryOperation, tableName string) error {
	// Delete stuff with condition from v here
	keys, values, _ := getField(v, true)
	conditionFieldsString := transformNamesToUpdateQuery(keys, 1, " AND ")

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, conditionFieldsString)

	_, err := db.Exec(query, values...)

	return err
}

func (q *Querynator) Update(v interface{}, conditionNames []string, conditionValues []any, db QueryOperation, tableName string) error {
	// Update stuff from v with condition here
	keys, values, _ := getField(v, true)

	updateFieldsString := transformNamesToUpdateQuery(keys, 1, ",")
	conditionFieldsString := transformNamesToUpdateQuery(conditionNames, len(keys)+1, " AND ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, updateFieldsString, conditionFieldsString)

	values = append(values, conditionValues...)

	_, err := db.Exec(query, values...)

	return err
}

func (q *Querynator) IsExists(v interface{}, db *sql.DB, tableName string) (bool, error) {
	// Check if a record exist with any of the field in V
	//https://stackoverflow.com/questions/32554400/efficiently-determine-if-any-rows-satisfy-a-predicate-in-postgres?rq=3
	var exists bool

	keys, values, _ := getField(v, true)
	conditionString := transformNamesToUpdateQuery(keys, 1, " AND ")

	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s)", tableName, conditionString)

	err := db.QueryRow(query, values...).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (q *Querynator) FindOne(v interface{}, dest interface{}, db *sql.DB, tableName string, returnFieldsName ...string) error {
	dbSqlx := sqlx.NewDb(db, "postgres")

	keys, values, _ := getField(v, true)
	conditionString := transformNamesToUpdateQuery(keys, 1, " AND ")
	returnFieldsString := strings.Join(returnFieldsName, ",")

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s LIMIT 1",
		returnFieldsString, tableName, conditionString)

	err := dbSqlx.Get(dest, query, values...)

	return err

}

func (q *Querynator) Find(v interface{}, dest interface{}, limit int, db *sql.DB, tableName string, returnFieldsName ...string) error {
	// Do some query here

	dbSqlx := sqlx.NewDb(db, "postgres")

	keys, values, _ := getField(v, true)
	conditionString := transformNamesToUpdateQuery(keys, 1, " AND ")
	returnFieldsString := strings.Join(returnFieldsName, ",")

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		returnFieldsString, tableName, conditionString)

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	err := dbSqlx.Select(dest, query, values...)

	return err
}

func (q *Querynator) FindWithCondition(conditions []QueryCondition, dest interface{}, limit int, db *sql.DB, tableName string, returnFieldsName ...string) error {
	dbSqlx := sqlx.NewDb(db, "postgres")

	conditionStrings, conditionValues := constructConditionClause(conditions, 0)

	returnFieldsString := strings.Join(returnFieldsName, ",")

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		returnFieldsString, tableName, strings.Join(conditionStrings, " AND "))

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	err := dbSqlx.Select(dest, query, conditionValues...)

	return err
}

func getField(v interface{}, ignoreEmpty bool) ([]string, []any, int) {
	s := reflect.ValueOf(v).Elem()
	typeOfS := s.Type()

	names := make([]string, 0, 8)
	values := make([]any, 0, 8)

	emptyField := 0

	for i := 0; i < typeOfS.NumField(); i++ {
		field := typeOfS.Field(i)
		columnTag := field.Tag.Get("db")

		// Gatekeep conditional
		if columnTag == "-" || columnTag == "" {
			continue
		}

		k := strings.SplitAfter(columnTag, ",")[0]
		v := s.Field(i).Interface()

		if ignoreEmpty {
			// Check if a field is empty/has value of "zero"
			if v != reflect.Zero(s.Field(i).Type()).Interface() {
				names = append(names, k)
				values = append(values, v)
			}
			continue
		}

		names = append(names, k)
		values = append(values, v)

		emptyField++
	}

	return names, values, emptyField
}

func transformNamesToUpdateQuery(names []string, start int, sep string) string {
	q := ""
	c := start

	for _, k := range names {
		q += fmt.Sprintf("%s=$%d%s", k, c, sep)
		c++
	}

	return q[:len(q)-len(sep)]
}

func constructConditionClause(conditions []QueryCondition, offset int) ([]string, []any) {
	conditionStrings := make([]string, 0, len(conditions))
	valueArgs := make([]any, 0, len(conditions))

	for i, condition := range conditions {
		c := fmt.Sprintf("%s.%s%s$%d", condition.TableName, condition.ColumnName, condition.Operand, i+offset+1)
		conditionStrings = append(conditionStrings, c)
		valueArgs = append(valueArgs, condition.MatchValue)
	}

	return conditionStrings, valueArgs
}
