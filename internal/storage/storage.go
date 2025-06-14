package storage

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	ID            string
	EncryptedText []byte
	NonceMsg      []byte
	EncryptedKey  []byte
	NonceKey      []byte
	ReadOnce      bool
	CreatedAt     time.Time
	ExpiresAt     time.Time
}

type DBHandle struct {
	Conn *sql.DB
}

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	schema := `
    CREATE TABLE IF NOT EXISTS messages (
        id TEXT PRIMARY KEY,
        encrypted_text BLOB NOT NULL,
        nonce_msg BLOB NOT NULL,
        encrypted_key BLOB NOT NULL,
        nonce_key BLOB NOT NULL,
        read_once BOOLEAN NOT NULL,
        created_at DATETIME NOT NULL,
        expires_at DATETIME NOT NULL
    );
    `
	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SaveMessage(db *sql.DB, msg Message) error {
	_, err := db.Exec(`
        INSERT INTO messages (id, encrypted_text, nonce_msg, encrypted_key, nonce_key, read_once, created_at, expires_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `, msg.ID, msg.EncryptedText, msg.NonceMsg, msg.EncryptedKey, msg.NonceKey, msg.ReadOnce, msg.CreatedAt, msg.ExpiresAt)
	return err
}

func LoadMessage(db *sql.DB, id string) (*Message, error) {
	row := db.QueryRow(`
        SELECT encrypted_text, nonce_msg, encrypted_key, nonce_key, read_once, created_at, expires_at
        FROM messages
        WHERE id = ?
    `, id)

	var msg Message
	msg.ID = id
	err := row.Scan(&msg.EncryptedText, &msg.NonceMsg, &msg.EncryptedKey, &msg.NonceKey, &msg.ReadOnce, &msg.CreatedAt, &msg.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func DeleteMessage(db *sql.DB, id string) error {
	_, err := db.Exec(`DELETE FROM messages WHERE id = ?`, id)
	return err
}
