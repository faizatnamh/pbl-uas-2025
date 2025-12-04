package repository

import (
	"database/sql"
)

type PermissionRepository struct {
	DB *sql.DB
}

func NewPermissionRepository(db *sql.DB) *PermissionRepository {
	return &PermissionRepository{DB: db}
}

// Check whether a user (by userID string) has a specific permission name
func (r *PermissionRepository) UserHasPermission(userID string, permission string) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM users u
		JOIN roles ro ON u.role_id = ro.id
		JOIN role_permissions rp ON rp.role_id = ro.id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE u.id = $1 AND p.name = $2
	`
	err := r.DB.QueryRow(query, userID, permission).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Optional: get all permission names for a user
func (r *PermissionRepository) GetPermissionsByUser(userID string) ([]string, error) {
	query := `
		SELECT p.name
		FROM users u
		JOIN roles ro ON u.role_id = ro.id
		JOIN role_permissions rp ON rp.role_id = ro.id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE u.id = $1
	`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, nil
}

// Get permission list by role name (ex: "Admin", "Mahasiswa", "Dosen Wali")
func (r *PermissionRepository) GetPermissionsByRole(roleName string) ([]string, error) {
	query := `
		SELECT p.name
		FROM roles ro
		JOIN role_permissions rp ON rp.role_id = ro.id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE ro.name = $1
	`

	rows, err := r.DB.Query(query, roleName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}

	return perms, nil
}
