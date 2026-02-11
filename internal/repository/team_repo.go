package repository

import (
	"context"

	"final-by-me/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TeamRepo struct {
	col *mongo.Collection
}

func NewTeamRepo(db *mongo.Database) *TeamRepo {
	return &TeamRepo{col: db.Collection("teams")}
}

func (r *TeamRepo) EnsureIndexes(ctx context.Context) error {
	// _id is unique by default (we use Code as _id)
	return nil
}

func (r *TeamRepo) SeedIfEmpty(ctx context.Context, teams []models.Team) error {
	count, err := r.col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	docs := make([]any, 0, len(teams))
	for _, t := range teams {
		docs = append(docs, t)
	}
	_, err = r.col.InsertMany(ctx, docs)
	return err
}

func (r *TeamRepo) List(ctx context.Context) ([]models.Team, error) {
	cur, err := r.col.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.Team
	for cur.Next(ctx) {
		var t models.Team
		if err := cur.Decode(&t); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, cur.Err()
}

func (r *TeamRepo) Exists(ctx context.Context, code string) (bool, error) {
	c, err := r.col.CountDocuments(ctx, bson.M{"_id": code})
	return c > 0, err
}
