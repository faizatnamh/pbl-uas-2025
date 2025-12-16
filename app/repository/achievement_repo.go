package repository

import (
	"context"
	"time"

	"pbluas/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (r *AchievementRepository) FindByIDs(ctx context.Context, ids []string) ([]models.Achievement, error) {
	var objIDs []primitive.ObjectID
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		objIDs = append(objIDs, oid)
	}

	filter := bson.M{
		"_id": bson.M{"$in": objIDs},
	}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.Achievement
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *AchievementRepository) FindByID(ctx context.Context, id string) (*models.Achievement, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result models.Achievement
	err = r.Collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
