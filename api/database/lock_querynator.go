package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

type LockQuerynator struct {
	Tx             *sql.Tx
	useJoin        bool
	joinClause     JoinClause
	joinQuerynator *JoinQueryExecutor
	condition      []QueryCondition
	querynator     *Querynator
}

func (q *LockQuerynator) UseJoin(joinClause JoinClause) *LockQuerynator {
	q.useJoin = true
	q.joinClause = joinClause
	q.joinQuerynator = q.querynator.PrepareJoinOperation()
	return q
}

func (q *LockQuerynator) AddJoinTable(joinTableName string, joinKeyName string, receiverTableName string, receiverForeignKeyName string) {
	q.joinQuerynator.AddJoinTable(joinTableName, joinKeyName, receiverTableName, receiverForeignKeyName)
}

func (q *LockQuerynator) SetLock(selectTableName string, lockTableName string, dest interface{}, returnFields map[string][]string) error {
	var query string
	var err error

	conditionString, conditionValue := constructConditionClause(q.condition, 0, false)

	if q.useJoin {
		query, err = q.joinQuerynator.createQuery(conditionString, selectTableName, q.joinClause, returnFields)

		if err != nil {
			return err
		}

		query += " FOR UPDATE OF " + lockTableName

	} else {
		query = fmt.Sprintf("SELECT * FROM %s", selectTableName)

		query += " FOR UPDATE OF " + lockTableName

		if conditionString != nil {
			query += fmt.Sprintf("WHERE %s", strings.Join(conditionString, " AND "))
		}

	}

	// Follow the default way on creating a mapper function
	tx := &sqlx.Tx{Tx: q.Tx, Mapper: reflectx.NewMapperFunc("db", sqlx.NameMapper)}

	fmt.Println(query)

	if dest == nil {
		fmt.Println("dest is nil")
	}

	return tx.Select(dest, query, conditionValue...)
}

func (q *LockQuerynator) Update(v interface{}, tableName string, conditionNames []string, conditionValues []any) (sql.Result, error) {
	keys, values, _ := getField(v, true)

	updateFieldsString := transformNamesToUpdateQuery(keys, 1, ",")
	conditionFieldsString := transformNamesToUpdateQuery(conditionNames, len(keys)+1, " AND ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		tableName,
		updateFieldsString,
		conditionFieldsString,
	)

	values = append(values, conditionValues...)

	return q.Tx.Exec(query, values...)
}
