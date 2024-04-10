package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type JoinClause string
type OrderType string

const (
	INNER_JOIN JoinClause = "INNER JOIN"
	OUTER_JOIN JoinClause = "OUTER JOIN"
	LEFT_JOIN  JoinClause = "LEFT JOIN"
	RIGHT_JOIN JoinClause = "RIGHT JOIN"
)

const (
	ASCENDING  OrderType = "ASC"
	DESCENDING OrderType = "DESC"
)

type TableRelation struct {
	ReceiverTableName    string
	ReceiverIDColumnName string
	SourceTableName      string
	SourceIDColumnName   string
}

type JoinQueryExecutor struct {
	DriverName      string
	orderBy         string
	orderByType     OrderType
	limit           int
	useExplicitCast bool
	tableRelations  []TableRelation
}

func (e *JoinQueryExecutor) OrderBy(orderByField string, orderType OrderType) {
	e.orderBy = orderByField
	e.orderByType = orderType
}

func (e *JoinQueryExecutor) SetLimit(limit int) {
	e.limit = limit
}

func (e *JoinQueryExecutor) UseExplicitCast() {
	e.useExplicitCast = true
}

func (e *JoinQueryExecutor) AddJoinTable(joinTableName string, joinKeyName string, receiverTableName string, receiverForeignKeyName string) {
	e.tableRelations = append(e.tableRelations, TableRelation{
		SourceTableName:      joinTableName,
		SourceIDColumnName:   joinKeyName,
		ReceiverTableName:    receiverTableName,
		ReceiverIDColumnName: receiverForeignKeyName,
	})
}

// returnFields is a map with key is the table name and the value is the column name that you want to return
func (e *JoinQueryExecutor) Find(db *sql.DB, condition []QueryCondition, dest interface{}, tableName string, joinClause JoinClause, returnFields map[string][]string) error {
	// if reflect.TypeOf(dest).Kind() != reflect.Slice {
	// 	return fmt.Errorf("illegal arguments. dest must be a slice")
	// }

	join := ""

	for _, relation := range e.tableRelations {
		join += constructJoinClause(
			relation.SourceTableName,
			relation.SourceIDColumnName,
			relation.ReceiverTableName,
			relation.ReceiverIDColumnName,
			joinClause,
		) + " "
	}

	// construct where clause
	conditionString, valueArgs := constructConditionClause(condition, 0, e.useExplicitCast)

	fmt.Println(valueArgs...)

	whereClause := strings.Join(conditionString, " AND ")

	// construct select clause
	fields := make([]string, 0, len(returnFields))

	for k, v := range returnFields {
		for _, column := range v {
			nameSplit := strings.SplitN(column, ",", 2)
			if len(nameSplit) > 1 {
				s, err := constructDefaultValue(fmt.Sprintf("%s.%s", k, nameSplit[0]), nameSplit[0], nameSplit[1])

				if err != nil {
					return err
				}

				fields = append(fields, s)
				continue
			}
			fields = append(fields, fmt.Sprintf("%s.%s", k, column))
		}
	}

	returnFieldsString := strings.Join(fields, ",")

	// construct query
	query := fmt.Sprintf(`SELECT %s FROM %s %s WHERE %s`, returnFieldsString, tableName, join, whereClause)

	if e.orderBy != "" {
		query = fmt.Sprintf("%s ORDER BY %s %s", query, e.orderBy, e.orderByType)
	}

	if e.limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, e.limit)
	}

	fmt.Println(query)
	dbX := sqlx.NewDb(db, "postgres")

	err := dbX.Select(dest, query, valueArgs...)

	return err
}

func constructDefaultValue(returnFieldsName string, columnName string, dataType string) (string, error) {
	switch dataType {
	case "string":
		return fmt.Sprintf("COALESCE(%s, '') as %s", returnFieldsName, columnName), nil
	case "int":
		return fmt.Sprintf("COALESCE(%s, 0) as %s", returnFieldsName, columnName), nil
	case "bool":
		return fmt.Sprintf("COALESCE(%s, false) as %s", returnFieldsName, columnName), nil
	default:
		return "", fmt.Errorf("dataType isn't supported")

	}
}
func constructJoinClause(sourceTable string, sourceID string, receiverTable string, receiverID string, joinClause JoinClause) string {
	clause := fmt.Sprintf("%s %s ON %s.%s = %s.%s", joinClause, sourceTable, sourceTable, sourceID, receiverTable, receiverID)
	return clause
}

// func constructConditionClause(conditions []QueryCondition, offset int) ([]string, []any) {
// 	conditionStrings := make([]string, 0, len(conditions))
// 	valueArgs := make([]any, 0, len(conditions))

// 	for i, condition := range conditions {
// 		c := fmt.Sprintf("%s.%s%s$%d", condition.TableName, condition.ColumnName, condition.Operand, i+offset+1)
// 		conditionStrings = append(conditionStrings, c)
// 		valueArgs = append(valueArgs, condition.MatchValue)
// 	}

// 	return conditionStrings, valueArgs
// }
