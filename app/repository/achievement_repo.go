package repository

import (
	"context"
	"time"

	"pbluas/app/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementRepository struct {
	Collection *mongo.Collection
}

func NewAchievementRepository(db *mongo.Database) *AchievementRepository {
	return &AchievementRepository{
		Collection: db.Collection("achievements"),
	}
}

func (r *AchievementRepository) Create(ctx context.Context, a *models.Achievement) error {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	res, err := r.Collection.InsertOne(ctx, a)
	if err != nil {
		return err
	}

	a.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}
