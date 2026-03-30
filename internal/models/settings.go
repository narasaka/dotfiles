package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Setting struct {
	Key   string `db:"key" json:"key"`
	Value string `db:"value" json:"value"`
}

type SettingsStore struct {
	db *sqlx.DB
}

func NewSettingsStore(db *sqlx.DB) *SettingsStore {
	return &SettingsStore{db: db}
}

func (s *SettingsStore) GetAll() ([]Setting, error) {
	var settings []Setting
	err := s.db.Select(&settings, `SELECT key, value FROM settings ORDER BY key`)
	if err != nil {
		return nil, fmt.Errorf("get settings: %w", err)
	}
	return settings, nil
}

func (s *SettingsStore) Get(key string) (string, error) {
	var value string
	err := s.db.Get(&value, `SELECT value FROM settings WHERE key = ?`, key)
	if err != nil {
		return "", fmt.Errorf("get setting %s: %w", key, err)
	}
	return value, nil
}

func (s *SettingsStore) Set(key, value string) error {
	_, err := s.db.Exec(
		`INSERT INTO settings (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(key) DO UPDATE SET value=excluded.value, updated_at=CURRENT_TIMESTAMP`,
		key, value,
	)
	return err
}

func (s *SettingsStore) SetMultiple(settings map[string]string) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for key, value := range settings {
		_, err := tx.Exec(
			`INSERT INTO settings (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP)
			ON CONFLICT(key) DO UPDATE SET value=excluded.value, updated_at=CURRENT_TIMESTAMP`,
			key, value,
		)
		if err != nil {
			return fmt.Errorf("set setting %s: %w", key, err)
		}
	}

	return tx.Commit()
}
