// Clean Architecture - Interface Adapter Layer
// UserRepository implementation for PostgreSQL
package repository

import (
	"clean-go-rest-api/internal/domain/entity"
	"database/sql"
	"log"

	"github.com/google/uuid"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Add(user entity.User) error {
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

	_, err = tx.Exec("INSERT INTO users (id, name, email) VALUES ($1, $2, $3)", user.ID, user.Name, user.Email)
	if err != nil {
		return err
	}
	return tx.Commit()
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

	_, err = tx.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", user.Name, user.Email, user.ID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *PostgresUserRepository) GetById(id uuid.UUID) (entity.User, error) {
	var user entity.User
	err := r.db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email)
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
