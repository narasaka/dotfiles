package models

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID           string    `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Name         string    `db:"name" json:"name"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(email, passwordHash, name string) (*User, error) {
	user := &User{}
	err := s.db.QueryRowx(
		`INSERT INTO users (email, password_hash, name) VALUES (?, ?, ?) RETURNING *`,
		email, passwordHash, name,
	).StructScan(user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (s *UserStore) GetByEmail(email string) (*User, error) {
	user := &User{}
	err := s.db.Get(user, `SELECT * FROM users WHERE email = ?`, email)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil
}

func (s *UserStore) GetByID(id string) (*User, error) {
	user := &User{}
	err := s.db.Get(user, `SELECT * FROM users WHERE id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

func (s *UserStore) Count() (int, error) {
	var count int
	err := s.db.Get(&count, `SELECT COUNT(*) FROM users`)
	return count, err
}
