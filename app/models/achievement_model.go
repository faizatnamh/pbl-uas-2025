package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	StudentID       string              `bson:"studentId" json:"studentId"`
	AchievementType string              `bson:"achievementType" json:"achievementType"`
	Title           string              `bson:"title" json:"title"`
	Description     string              `bson:"description" json:"description"`
	Details         AchievementDetails  `bson:"details" json:"details"`
	Tags            []string            `bson:"tags" json:"tags"`
	Attachments []AchievementAttachment `bson:"attachments,omitempty" json:"attachments,omitempty"`
	CreatedAt       time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time           `bson:"updatedAt" json:"updatedAt"`
}

type AchievementDetails struct {
	CompetitionName  string `bson:"competitionName" json:"competitionName"`
	CompetitionLevel string `bson:"competitionLevel" json:"competitionLevel"`
	Rank *float64 `bson:"rank,omitempty" json:"rank,omitempty"`
	MedalType        string `bson:"medalType" json:"medalType"`
	EventDate        string `bson:"eventDate" json:"eventDate"`
	Location         string `bson:"location" json:"location"`
	Organizer        string `bson:"organizer" json:"organizer"`
}
// ===== REQUEST BODY (CREATE ACHIEVEMENT) =====
type AchievementCreateRequest struct {
	StudentID       string             `json:"studentId,omitempty"` 
	AchievementType string             `json:"achievementType"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	Details         AchievementDetails `json:"details"` // ✅ FIX
	Tags            []string           `json:"tags"`
}

type AchievementAttachment struct {
	FileName   string    `bson:"fileName" json:"fileName"`
	FileURL    string    `bson:"fileUrl" json:"fileUrl"`
	FileType   string    `bson:"fileType" json:"fileType"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploadedAt"`
}
