package database

import (
	"database/sql"
	"fmt"

	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

func Open(path string) (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(FULL)", path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		fmt.Println("failed to open database")
		return nil, err
	}
	fmt.Println("database opened successfully")

	if err := db.Ping(); err != nil {
		fmt.Println("failed to ping database")
		return nil, err
	}
	fmt.Println("database pinged successfully")

	if _, err := db.Exec(schema); err != nil {
		fmt.Println("failed to execute schema")
		return nil, err
	}
	fmt.Println("schema executed successfully")
	return db, nil
}
