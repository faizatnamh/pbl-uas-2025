package repository

import (
	"database/sql"
	"errors"
	"pbluas/app/models"
)

type UserRepository interface {
	FindByUsername(username string) (*models.User, error)
}

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	user := models.User{}

	query := `
        SELECT 
            users.id,
            users.username,
            users.email,
            users.full_name,
            users.password_hash,
            users.is_active,
            roles.name
        FROM users
        JOIN roles ON roles.id = users.role_id
        WHERE users.username = $1
        LIMIT 1
    `
	err := r.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.PasswordHash,
		&user.IsActive,
		&user.RoleName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
