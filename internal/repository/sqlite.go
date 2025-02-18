package repository

import (
    "database/sql"

    _ "github.com/mattn/go-sqlite3"
)

// GetLegacyData retrieves read‑only data from a legacy SQLite database.
func GetLegacyData(query string, args ...interface{}) (*sql.Rows, error) {
    db, err := sql.Open("sqlite3", "./legacy.db")
    if err != nil {
        return nil, err
    }
    return db.Query(query, args...)
}
