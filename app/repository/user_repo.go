package repository

import (
	"database/sql"
	"errors"
	"pbluas/app/models"
	"fmt"
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
			users.password_hash,
			users.full_name,
			users.is_active,
			users.role_id,
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
		&user.PasswordHash,
		&user.FullName,
		&user.IsActive,
		&user.RoleID, 
		&user.RoleName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	
fmt.Println("DEBUG: SCANNED USER = ")
fmt.Println(" ID        =", user.ID)
fmt.Println(" USERNAME  =", user.Username)
fmt.Println(" EMAIL     =", user.Email)
fmt.Println(" FULLNAME  =", user.FullName)
fmt.Println(" HASH      =", user.PasswordHash)
fmt.Println(" ACTIVE    =", user.IsActive)
fmt.Println(" ROLE      =", user.RoleName)
return &user, nil
}
