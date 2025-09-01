package infrastructure

import (
	"database/sql"
	"log"
	"user_service/internal/domain"
)

type DB struct {
	DB *sql.DB
}

func NewDB(db *sql.DB) *DB {
	return &DB{DB: db}
}

func (db *DB) Register(input domain.AuthorizationInput) (*domain.User, error) {
	var user domain.User
	query := `INSERT INTO users(name, password_hash, role) VALUES($1, $2, $3) RETURNING id, name, password_hash, role, registered_at`
	err := db.DB.QueryRow(query, input.Name, input.Password, input.Role).Scan(&user.ID, &user.Name, &user.PasswordHash, &user.Role, &user.RegisteredAt)
	if err != nil {
		log.Printf("error to register user:%v", err)
		return nil, err
	}

	return &user, nil
}

func (db *DB) Login(input domain.AuthorizationInput) (*domain.User, error) {
	var user domain.User
	query := `SELECT id,name,password_Hash,role,registered_at FROM users WHERE name = $1`
	err := db.DB.QueryRow(query, input.Name).Scan(&user.ID, &user.Name, &user.PasswordHash, &user.Role, &user.RegisteredAt)
	if err != nil {
		log.Printf("error user not found:%v", err)
		return nil, err
	}

	return &user, nil

}
