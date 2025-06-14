package storage

import (
	"database/sql"
	"time"
)

type LogEntry struct {
	MessageID string
	IPAddress string
	Country   string
	Timestamp time.Time
}

func InitLogs(db *sql.DB) error {
	schema := `
    CREATE TABLE IF NOT EXISTS read_logs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        message_id TEXT,
        ip_address TEXT,
        country TEXT,
        timestamp DATETIME
    );

    CREATE TABLE IF NOT EXISTS write_logs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        message_id TEXT,
        ip_address TEXT,
        country TEXT,
        timestamp DATETIME
    );
    `
	_, err := db.Exec(schema)
	return err
}

func LogRead(db *sql.DB, entry LogEntry) error {
	_, err := db.Exec(`
        INSERT INTO read_logs (message_id, ip_address, country, timestamp)
        VALUES (?, ?, ?, ?)
    `, entry.MessageID, entry.IPAddress, entry.Country, entry.Timestamp)
	return err
}

func LogWrite(db *sql.DB, entry LogEntry) error {
	_, err := db.Exec(`
        INSERT INTO write_logs (message_id, ip_address, country, timestamp)
        VALUES (?, ?, ?, ?)
    `, entry.MessageID, entry.IPAddress, entry.Country, entry.Timestamp)
	return err
}
