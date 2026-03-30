package db

import (
	"fmt"
	"log"
	"sort"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"github.com/kubeploy/kubeploy/internal/db/migrations"
)

func Open(dbPath string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("file:%s?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on", dbPath)
	db, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(1)

	if err := runMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return db, nil
}

func runMigrations(db *sqlx.DB) error {
	entries, err := migrations.FS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var sqlFiles []string
	for _, e := range entries {
		if !e.IsDir() && len(e.Name()) > 4 && e.Name()[len(e.Name())-4:] == ".sql" {
			sqlFiles = append(sqlFiles, e.Name())
		}
	}
	sort.Strings(sqlFiles)

	for _, name := range sqlFiles {
		data, err := migrations.FS.ReadFile(name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		log.Printf("Running migration: %s", name)
		if _, err := db.Exec(string(data)); err != nil {
			return fmt.Errorf("execute migration %s: %w", name, err)
		}
	}

	return nil
}
