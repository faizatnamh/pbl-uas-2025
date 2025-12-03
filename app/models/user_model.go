package models

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	PasswordHash string `json:"-"`      
	RoleName     string `json:"role"`
	IsActive     bool   `json:"is_active"`
}
