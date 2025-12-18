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

func (r *AchievementRepository) UpdateByID(
	ctx context.Context,
	id string,
	req models.AchievementCreateRequest,
	points int,
) error {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"achievementType": req.AchievementType,
			"title":           req.Title,
			"description":     req.Description,
			"details":         req.Details,
			"tags":            req.Tags,
			"points": points,
			"updatedAt":       time.Now(),
		},
	}

	_, err = r.Collection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		update,
	)

	return err
}

func (r *AchievementRepository) AddAttachment(
	ctx context.Context,
	id string,
	attachment models.AchievementAttachment,
) error {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$push": bson.M{
			"attachments": attachment,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}

	_, err = r.Collection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		update,
	)

	return err
}
