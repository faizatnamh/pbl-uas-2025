package models

import "time"

type AchievementReference struct {
    ID            string     `db:"id"`
    StudentID     string     `db:"student_id"`
    MongoID       string     `db:"mongo_achievement_id"`
    Status        string     `db:"status"` // draft, submitted, verified, rejected, deleted    
    SubmittedAt   *time.Time `db:"submitted_at"`
    VerifiedAt    *time.Time `db:"verified_at"`
    VerifiedBy    *string    `db:"verified_by"`
    RejectionNote *string    `db:"rejection_note"`
    CreatedAt     time.Time  `db:"created_at"`
    UpdatedAt     time.Time  `db:"updated_at"`

}