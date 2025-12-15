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
