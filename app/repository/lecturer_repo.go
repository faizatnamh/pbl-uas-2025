package repository

import (
	"database/sql"
	"pbluas/app/models"
)

type LecturerRepository interface {
	GetLecturerByUserID(userID string) (*models.Lecturer, error)
	GetLecturerByID(id string) (*models.Lecturer, error)
}

type lecturerRepository struct {
	DB *sql.DB
}

func NewLecturerRepository(db *sql.DB) LecturerRepository {
	return &lecturerRepository{DB: db}
}

func (r *lecturerRepository) GetLecturerByUserID(userID string) (*models.Lecturer, error) {
	var lec models.Lecturer

	query := `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE user_id = $1
	`

	err := r.DB.QueryRow(query, userID).Scan(
		&lec.ID,
		&lec.UserID,
		&lec.LecturerID,
		&lec.Department,
		&lec.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &lec, nil
}

func (r *lecturerRepository) GetLecturerByID(id string) (*models.Lecturer, error) {
	var lec models.Lecturer

	query := `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE id = $1
	`

	err := r.DB.QueryRow(query, id).Scan(
		&lec.ID,
		&lec.UserID,
		&lec.LecturerID,
		&lec.Department,
		&lec.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &lec, nil
}
