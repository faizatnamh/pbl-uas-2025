package repository

import (
	"database/sql"
	"pbluas/app/models"
)

type StudentRepository interface {
	GetAllStudents() ([]models.StudentDetail, error)
	GetStudentsByAdvisor(advisorID string) ([]models.StudentDetail, error)
	GetStudentByID(id string) (*models.StudentDetail, error)
}

type studentRepository struct {
	DB *sql.DB
}

func NewStudentRepository(db *sql.DB) StudentRepository {
	return &studentRepository{DB: db}
}

// ADMIN → GET ALL STUDENTS
func (r *studentRepository) GetAllStudents() ([]models.StudentDetail, error) {
	query := `
		SELECT 
			s.id,
			s.user_id,
			s.student_id,
			u.full_name,
			u.email,
			s.program_study,
			s.academic_year,
			s.advisor_id,
			lec_user.full_name AS advisor_name
		FROM students s
		JOIN users u ON u.id = s.user_id
		LEFT JOIN lecturers lec ON lec.id = s.advisor_id
		LEFT JOIN users lec_user ON lec_user.id = lec.user_id
	`

	rows, err := r.DB.Query(query)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []models.StudentDetail

	for rows.Next() {
		var s models.StudentDetail
		var advisorID sql.NullString
		var advisorName sql.NullString

		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.StudentID,
			&s.FullName,
			&s.Email,
			&s.ProgramStudy,
			&s.AcademicYear,
			&advisorID,
			&advisorName,
		)
		if err != nil { return nil, err }

		if advisorID.Valid { tmp := advisorID.String; s.AdvisorID = &tmp }
		if advisorName.Valid { tmp := advisorName.String; s.AdvisorName = &tmp }

		list = append(list, s)
	}

	return list, nil
}

// DOSEN → GET BY ADVISOR
func (r *studentRepository) GetStudentsByAdvisor(advisorID string) ([]models.StudentDetail, error) {
	query := `
		SELECT 
			s.id,
			s.user_id,
			s.student_id,
			u.full_name,
			u.email,
			s.program_study,
			s.academic_year,
			s.advisor_id,
			lec_user.full_name AS advisor_name
		FROM students s
		JOIN users u ON u.id = s.user_id
		LEFT JOIN lecturers lec ON lec.id = s.advisor_id
		LEFT JOIN users lec_user ON lec_user.id = lec.user_id
		WHERE s.advisor_id = $1
	`

	rows, err := r.DB.Query(query, advisorID)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []models.StudentDetail

	for rows.Next() {
		var s models.StudentDetail
		var advID sql.NullString
		var advName sql.NullString

		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.StudentID,
			&s.FullName,
			&s.Email,
			&s.ProgramStudy,
			&s.AcademicYear,
			&advID,
			&advName,
		)
		if err != nil { return nil, err }

		if advID.Valid { tmp := advID.String; s.AdvisorID = &tmp }
		if advName.Valid { tmp := advName.String; s.AdvisorName = &tmp }

		list = append(list, s)
	}

	return list, nil
}

// GET DETAIL BY ID

func (r *studentRepository) GetStudentByID(id string) (*models.StudentDetail, error) {
	query := `
		SELECT 
			s.id,
			s.user_id,
			s.student_id,
			u.full_name,
			u.email,
			s.program_study,
			s.academic_year,
			s.advisor_id,
			lec_user.full_name AS advisor_name
		FROM students s
		JOIN users u ON u.id = s.user_id
		LEFT JOIN lecturers lec ON lec.id = s.advisor_id
		LEFT JOIN users lec_user ON lec_user.id = lec.user_id
		WHERE s.id = $1
	`

	var s models.StudentDetail
	var advID sql.NullString
	var advName sql.NullString

	err := r.DB.QueryRow(query, id).Scan(
		&s.ID,
		&s.UserID,
		&s.StudentID,
		&s.FullName,
		&s.Email,
		&s.ProgramStudy,
		&s.AcademicYear,
		&advID,
		&advName,
	)
	if err != nil { return nil, err }

	if advID.Valid { tmp := advID.String; s.AdvisorID = &tmp }
	if advName.Valid { tmp := advName.String; s.AdvisorName = &tmp }

	return &s, nil
}
