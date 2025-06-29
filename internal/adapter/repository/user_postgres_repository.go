// Clean Architecture - Interface Adapter Layer
// UserRepository implementation for PostgreSQL
package repository

import (
	"clean-go-rest-api/internal/domain/entity"
	"database/sql"
	"log"

	"github.com/google/uuid"
)

type TxExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Rollback() error
	Commit() error
}

type DBExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Begin() (TxExecutor, error)
}

type PostgresUserRepository struct {
	db DBExecutor
}

func NewPostgresUserRepository(db DBExecutor) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Add(user entity.User) error {
	_, err := r.db.Exec(
		`INSERT INTO users (id, name, email) VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, email = EXCLUDED.email`,
		user.ID, user.Name, user.Email,
	)
	return err
}

func (r *PostgresUserRepository) Delete(user entity.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Advisory lock by user ID to prevent concurrent modifications
	_, err = tx.Exec("SELECT pg_advisory_xact_lock($1)", user.ID.ID())
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM users WHERE id = $1", user.ID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *PostgresUserRepository) Update(user entity.User) error {
	_, err := r.db.Exec(
		`UPDATE users (id, name, email) VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, email = EXCLUDED.email`,
		user.ID, user.Name, user.Email,
	)
	return err
}

func (r *PostgresUserRepository) GetById(id uuid.UUID) (entity.User, error) {
	var user entity.User
	err := r.db.QueryRow(
		"SELECT id, name, email FROM users WHERE id = $1", id,
	).Scan(&user.ID, &user.Name, &user.Email)
	if err == sql.ErrNoRows {
		return entity.User{}, nil
	}
	return user, err
}

func (r *PostgresUserRepository) Search(name string) ([]entity.User, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM users WHERE name ILIKE $1", "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *PostgresUserRepository) EmailExists(email string) bool {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		log.Println("Error checking if email exists:", err)
		return false
	}
	return exists
}
