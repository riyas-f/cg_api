package database

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type A struct {
	Bid int `db:"b_id"`
}

const (
	DATABASE_HOST     = "localhost"
	DATABASE_PORT     = 3565
	DATABASE_USERNAME = "test"
	DATABASE_PASSWORD = "admin"
	DATABASE_NAME     = "test_db"
)

func TestBatchInsert(t *testing.T) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DATABASE_HOST, DATABASE_PORT, DATABASE_USERNAME, DATABASE_PASSWORD, DATABASE_NAME)

	db, err := sqlx.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	data := []*A{
		{Bid: 100},
		{Bid: 250},
		{Bid: 300},
		{Bid: 400},
		{Bid: 2},
		{Bid: 3},
		{Bid: 4},
	}

	querynator := Querynator{}

	filterExecutor := querynator.PrepareFilterOperation()
	filterExecutor.AddTableSource("B", "id", "b_id")
	err = filterExecutor.UseTransaction(db.DB)

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	filterExecutor.UseExplicitCast()
	err = filterExecutor.BatchInsert(data, db.DB, "A")

	if err != nil {
		if filterExecutor.Tx != nil {
			filterExecutor.Rollback()
		}
		t.Errorf(err.Error())
		return
	}

	filterExecutor.Commit()
}

func TestJoinTableData(t *testing.T) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DATABASE_HOST, DATABASE_PORT, DATABASE_USERNAME, DATABASE_PASSWORD, DATABASE_NAME)

	db, err := sqlx.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	querynator := Querynator{}
	joinExecutor := querynator.PrepareJoinOperation()

	joinExecutor.AddJoinTable("account", "id", "balance", "user_id")

	temp := []struct {
		Name     string `db:"name"`
		Email    string `db:"email"`
		Balance  string `db:"balance"`
		Currency string `db:"currency"`
	}{}

	err = joinExecutor.Find(db.DB,
		[]QueryCondition{
			{TableName: "balance", ColumnName: "currency", MatchValue: "USD", Operand: EQ},
			{TableName: "balance", ColumnName: "balance", MatchValue: "3000", Operand: GEQ},
		}, &temp, "balance", INNER_JOIN, map[string][]string{
			"account": {"name", "email"},
			"balance": {"balance", "currency"},
		})

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	fmt.Println(temp)
}

func TestJoinThreeTable(t *testing.T) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DATABASE_HOST, DATABASE_PORT, DATABASE_USERNAME, DATABASE_PASSWORD, DATABASE_NAME)

	db, err := sqlx.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	querynator := Querynator{}
	joinExecutor := querynator.PrepareJoinOperation()

	joinExecutor.AddJoinTable("account", "id", "balance", "user_id")
	joinExecutor.AddJoinTable("home", "user_id", "account", "id")

	temp := []struct {
		Name     string `db:"name"`
		Email    string `db:"email"`
		Balance  string `db:"balance"`
		Currency string `db:"currency"`
		Country  string `db:"country"`
	}{}

	err = joinExecutor.Find(db.DB,
		[]QueryCondition{
			{TableName: "balance", ColumnName: "currency", MatchValue: "IDR", Operand: EQ},
			{TableName: "balance", ColumnName: "balance", MatchValue: "3000", Operand: GEQ},
		}, &temp, "balance", INNER_JOIN, map[string][]string{
			"account": {"name", "email"},
			"balance": {"balance", "currency"},
			"home":    {"country"},
		})

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	fmt.Println(temp)
}
func TestBatchInsertHardCode(t *testing.T) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DATABASE_HOST, DATABASE_PORT, DATABASE_USERNAME, DATABASE_PASSWORD, DATABASE_NAME)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	// bIDs := []interface{}{1, 2, 3, 4, 5, 6, 7} // replace with your actual integer values

	_, err = db.Exec(`WITH data(b_id) AS (VALUES ($1),($2),($3),($4),($5),($6),($7)) INSERT INTO A (b_id) SELECT d.b_id FROM data d WHERE EXISTS ( SELECT 1 FROM B WHERE B.id = d.b_id::INTEGER)`, 1, 2, 3, 4, 5, 6, 7)

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// 	data := []*A{
	// 		{Bid: 100},
	// 		{Bid: 250},
	// 		{Bid: 300},
	// 		{Bid: 400},
	// 		{Bid: 2},
	// 		{Bid: 3},
	// 		{Bid: 4},
	// 	}

	// 	valueStrings := []string{}
	// 	valueArgs := []interface{}{}
	// 	for _, w := range data {
	// 		valueStrings = append(valueStrings, "(?)")

	// 		valueArgs = append(valueArgs, w.Bid)
	// 	}

	// 	smt := `WITH data(b_id) AS (
	// 	VALUES %s
	//   ) INSERT INTO A(b_id)
	//   SELECT d.b_id
	//   FROM data d
	//   WHERE EXISTS (
	//   	SELECT 1
	//   	FROM B
	//   	WHERE B.id = d.b_id`

	// 	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	// 	_, err = db.Exec(smt, valueArgs...)
	// 	if err != nil {
	// 		t.Errorf(err.Error())
	// 		return
	// 	}

}
