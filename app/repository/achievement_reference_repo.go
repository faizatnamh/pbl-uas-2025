package repository

import (
	"database/sql"
	"time"

	"pbluas/app/models"

	"github.com/google/uuid"
)

type AchievementReferenceRepository struct {
	DB *sql.DB
}

func NewAchievementReferenceRepository(db *sql.DB) *AchievementReferenceRepository {
	return &AchievementReferenceRepository{DB: db}
}

// ================= CREATE =================

func (r *AchievementReferenceRepository) Create(ref *models.AchievementReference) error {
	ref.ID = uuid.NewString()
	ref.Status = "draft"
	ref.CreatedAt = time.Now()
	ref.UpdatedAt = time.Now()

	query := `
		INSERT INTO achievement_references
		(id, student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6)
	`

	_, err := r.DB.Exec(
		query,
		ref.ID,
		ref.StudentID,
		ref.MongoID,
		ref.Status,
		ref.CreatedAt,
		ref.UpdatedAt,
	)
	return err
}

// ================= GET BY STUDENT (Mahasiswa) =================

func (r *AchievementReferenceRepository) GetByStudentID(studentID string) ([]models.AchievementReference, error) {
	query := `
		SELECT 
			id, student_id, mongo_achievement_id, status,
			submitted_at, verified_at, verified_by, rejection_note,
			created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1
		  AND status != 'deleted'
	`

	rows, err := r.DB.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AchievementReference
	for rows.Next() {
		var ref models.AchievementReference
		if err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
			&ref.VerifiedBy,
			&ref.RejectionNote,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, ref)
	}

	return list, nil
}

// ================= GET ALL (Admin) =================

func (r *AchievementReferenceRepository) GetAll() ([]models.AchievementReference, error) {
	query := `
		SELECT 
			id, student_id, mongo_achievement_id, status,
			submitted_at, verified_at, verified_by, rejection_note,
			created_at, updated_at
		FROM achievement_references
		WHERE status != 'deleted'
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AchievementReference
	for rows.Next() {
		var ref models.AchievementReference
		if err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
			&ref.VerifiedBy,
			&ref.RejectionNote,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, ref)
	}

	return list, nil
}

// ================= GET BY ADVISOR (Dosen Wali) =================

func (r *AchievementReferenceRepository) GetByAdvisorUserID(userID string) ([]models.AchievementReference, error) {
	query := `
		SELECT 
			ar.id, ar.student_id, ar.mongo_achievement_id, ar.status,
			ar.submitted_at, ar.verified_at, ar.verified_by, ar.rejection_note,
			ar.created_at, ar.updated_at
		FROM achievement_references ar
		JOIN students s ON s.id = ar.student_id
		JOIN lecturers l ON l.id = s.advisor_id
		WHERE l.user_id = $1
		  AND ar.status != 'deleted'
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AchievementReference
	for rows.Next() {
		var ref models.AchievementReference
		if err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
			&ref.VerifiedBy,
			&ref.RejectionNote,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, ref)
	}

	return list, nil
}

func (r *AchievementReferenceRepository) GetByMongoID(mongoID string) (*models.AchievementReference, error) {
	query := `
		SELECT 
			id, student_id, mongo_achievement_id, status,
			submitted_at, verified_at, verified_by, rejection_note,
			created_at, updated_at
		FROM achievement_references
		WHERE mongo_achievement_id = $1
		  AND status != 'deleted'
	`

	var ref models.AchievementReference
	err := r.DB.QueryRow(query, mongoID).Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoID,
		&ref.Status,
		&ref.SubmittedAt,
		&ref.VerifiedAt,
		&ref.VerifiedBy,
		&ref.RejectionNote,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &ref, nil
}

func (r *AchievementReferenceRepository) IsAdvisorOfStudent(userID, studentID string) (bool, error) {
	query := `
		SELECT COUNT(1)
		FROM students s
		JOIN lecturers l ON l.id = s.advisor_id
		WHERE s.id = $1 AND l.user_id = $2
	`

	var count int
	err := r.DB.QueryRow(query, studentID, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AchievementReferenceRepository) SoftDelete(id string) error {
	query := `
		UPDATE achievement_references
		SET status = 'deleted',
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.DB.Exec(query, id)
	return err
}

func (r *AchievementReferenceRepository) Submit(id string) error {
	query := `
		UPDATE achievement_references
		SET status = 'submitted',
		    submitted_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
		  AND status = 'draft'
	`

	res, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *AchievementReferenceRepository) Verify(id string, verifiedBy string) error {
	query := `
		UPDATE achievement_references
		SET status = 'verified',
		    verified_at = NOW(),
		    verified_by = $2,
		    updated_at = NOW()
		WHERE id = $1
		  AND status = 'submitted'
	`

	res, err := r.DB.Exec(query, id, verifiedBy)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *AchievementReferenceRepository) Reject(id string, rejectedBy string, note string) error {
	query := `
		UPDATE achievement_references
		SET status = 'rejected',
		    rejection_note = $2,
		    verified_at = NOW(),
		    verified_by = $3,
		    updated_at = NOW()
		WHERE id = $1
		  AND status = 'submitted'
	`

	res, err := r.DB.Exec(query, id, note, rejectedBy)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *AchievementReferenceRepository)GetByStudentIDForReport(studentID string) ([]models.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status
		FROM achievement_references
		WHERE student_id = $1
		  AND status != 'deleted'
	`

	rows, err := r.DB.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AchievementReference
	for rows.Next() {
		var ref models.AchievementReference
		if err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoID,
			&ref.Status,
		); err != nil {
			return nil, err
		}
		list = append(list, ref)
	}

	return list, nil
}

