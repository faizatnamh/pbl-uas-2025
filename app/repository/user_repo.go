package repository

import (
	"database/sql"
	"errors"
	"pbluas/app/models"
	"fmt"
)

type UserRepository interface {
    FindByUsername(username string) (*models.User, error)
    FindByUserID(id string) (*models.User, error)
    // ADMIN CRUD
    GetAllUsers() ([]models.User, error)
    CreateUser(user *models.User) error
    UpdateUser(user *models.User) error
    DeleteUser(id string) error
    UpdateUserRole(userID string, roleID string) error
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

// FindByUserID retrieves a user by their ID
func (r *userRepository) FindByUserID(id string) (*models.User, error) {
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
		WHERE users.id = $1
		LIMIT 1
	`

	err := r.DB.QueryRow(query, id).Scan(
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

	return &user, nil
}

// ===== CREATE USER =====
func (r *userRepository) CreateUser(user *models.User) error {
    query := `
        INSERT INTO users (username, email, password_hash, full_name, role_id, is_active)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
    _, err := r.DB.Exec(query,
        user.Username,
        user.Email,
        user.PasswordHash,
        user.FullName,
        user.RoleID,
        user.IsActive,
    )
    return err
}

// ===== GET ALL USERS =====
func (r *userRepository) GetAllUsers() ([]models.User, error) {
    query := `
        SELECT u.id, u.username, u.email, u.full_name, u.role_id, r.name, u.is_active
        FROM users u
        JOIN roles r ON r.id = u.role_id
    `
    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    users := []models.User{}
    for rows.Next() {
        var u models.User
        rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.RoleName, &u.IsActive)
        users = append(users, u)
    }

    return users, nil
}

// ===== UPDATE USER =====
func (r *userRepository) UpdateUser(user *models.User) error {
    query := `
        UPDATE users 
        SET email=$1, full_name=$2, role_id=$3, is_active=$4, updated_at=NOW()
        WHERE id=$5
    `
    _, err := r.DB.Exec(query,
        user.Email,
        user.FullName,
        user.RoleID,
        user.IsActive,
        user.ID,
    )
    return err
}

// ===== DELETE USER =====
func (r *userRepository) DeleteUser(id string) error {
    query := `DELETE FROM users WHERE id=$1`
    _, err := r.DB.Exec(query, id)
    return err
}

// ===== UPDATE ONLY ROLE =====
func (r *userRepository) UpdateUserRole(userID string, roleID string) error {
    query := `
        UPDATE users SET role_id=$1, updated_at=NOW()
        WHERE id=$2
    `
    _, err := r.DB.Exec(query, roleID, userID)
    return err
}

