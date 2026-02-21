package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dataSourceName string) error {
	var err error
	DB, err = sql.Open("sqlite", dataSourceName+"?_pragma=busy_timeout=10000&_pragma=journal_mode=WAL")
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	// Make sure schemas are run correctly
	schemaBytes, err := os.ReadFile("migrations/schema.sql")
	if err == nil {
		_, err = DB.Exec(string(schemaBytes))
		if err != nil {
			return fmt.Errorf("failed to execute schema: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
