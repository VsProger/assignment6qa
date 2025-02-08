package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/models"

	"golang.org/x/crypto/bcrypt"
)

func (s *Sqlite) GetUserByEmail(email string) (*models.User, error) {
	op := "sqlite.GetUserByEmail"
	var u models.User
	stmt := `SELECT id, name, email, created FROM users WHERE id=?`
	err := s.db.QueryRow(stmt, email).Scan(&u.ID, &u.Name, &u.Email, &u.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &u, nil
}

func (s *Sqlite) UpdateUserByID(string) (*models.User, error) { return nil, nil }

func (s *Sqlite) CreateUser(u models.User) error {
	op := "sqlite.CreateUser"
	stmt := `INSERT INTO users (name, email,hashed_password, created) VALUES(?, ?, ?, CURRENT_TIMESTAMP)`
	_, err := s.db.Exec(stmt, u.Name, u.Email, string(u.HashedPassword))
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return models.ErrDuplicateEmail
		}
		if err.Error() == "UNIQUE constraint failed: users.name" {
			return models.ErrDuplicateName
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Sqlite) GetUserByID(id int) (*models.User, error) {
	op := "sqlite.GetUserByID"
	var u models.User
	stmt := `SELECT id, name, email, created FROM users WHERE id=?`
	err := s.db.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &u, nil
}

func (s *Sqlite) Authenticate(email, password string) (int, error) {
	op := "sqlite.Authenticate"
	var id int
	var hashed_password []byte
	stmt := `SELECT id, hashed_password FROM users WHERE email=?`
	err := s.db.QueryRow(stmt, email).Scan(&id, &hashed_password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecord
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	err = bcrypt.CompareHashAndPassword(hashed_password, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
