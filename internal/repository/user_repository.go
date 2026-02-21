package repository

import (
	"database/sql"
	"event-registration-system/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(u *models.User) (int64, error) {
	query := `INSERT INTO users (name, email) VALUES (?, ?)`
	res, err := r.db.Exec(query, u.Name, u.Email)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	query := `SELECT id, name, email FROM users WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var u models.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email); err != nil {
		return nil, err
	}
	return &u, nil
}
