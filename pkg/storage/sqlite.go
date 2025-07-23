package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sajjxd/pubsub-scan-processor/pkg/types"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(dbPath string) (*Repository, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := createTable(db); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &Repository{db: db}, nil
}

func createTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS scans (
        ip TEXT NOT NULL,
        port INTEGER NOT NULL,
        service TEXT NOT NULL,
        response TEXT,
        last_scanned DATETIME NOT NULL,
        PRIMARY KEY (ip, port, service)
    );`
	_, err := db.Exec(query)
	return err
}

func (r *Repository) UpsertRecord(record types.ScanRecord) error {
	query := `
    INSERT INTO scans (ip, port, service, response, last_scanned)
    VALUES (?, ?, ?, ?, ?)
    ON CONFLICT(ip, port, service) DO UPDATE SET
        response = excluded.response,
        last_scanned = excluded.last_scanned;`

	_, err := r.db.Exec(query, record.Ip, record.Port, record.Service, record.Response, record.LastScanned.UTC())
	if err != nil {
		return fmt.Errorf("failed to upsert record: %w", err)
	}

	log.Printf("Processed record for %s:%d (%s)", record.Ip, record.Port, record.Service)
	return nil
}

func (r *Repository) Close() {
	r.db.Close()
}
