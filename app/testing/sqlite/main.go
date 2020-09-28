package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if _, err := db.Exec("DROP TABLE IF EXISTS test"); err != nil {
		panic(err)
	}
	if _, err := db.Exec("CREATE TABLE test(id int)"); err != nil {
		panic(err)
	}
	if _, err := db.Exec("INSERT INTO test(id)VALUES(1)"); err != nil {
		panic(err)
	}
	result := 0
	if err := db.QueryRow("SELECT id FROM test LIMIT 1").Scan(&result); err != nil {
		panic(err)
	}
	fmt.Println("testing pass :", result)
}
