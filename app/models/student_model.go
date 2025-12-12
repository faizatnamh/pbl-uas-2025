package models

import "time"

type Student struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	StudentID    string    `json:"student_id" db:"student_id"`
	ProgramStudy string    `json:"program_study" db:"program_study"`
	AcademicYear string    `json:"academic_year" db:"academic_year"`
	AdvisorID    *string   `json:"advisor_id" db:"advisor_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// MODEL DETAIL 
type StudentDetail struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	StudentID    string  `json:"student_id"`
	FullName     string  `json:"full_name"`
	Email        string  `json:"email"`
	ProgramStudy string  `json:"program_study"`
	AcademicYear string  `json:"academic_year"`
	AdvisorID    *string `json:"advisor_id"`  
	AdvisorName  *string `json:"advisor_name"`
}
