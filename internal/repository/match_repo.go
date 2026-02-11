package repository

import (
	"context"

	"final-by-me/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MatchRepo struct {
	col *mongo.Collection
}

func NewMatchRepo(db *mongo.Database) *MatchRepo {
	return &MatchRepo{col: db.Collection("matches")}
}

func (r *MatchRepo) EnsureIndexes(ctx context.Context) error {
	// matchKey unique
	_, err := r.col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "matchKey", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}

func (r *MatchRepo) Create(ctx context.Context, m models.Match) (models.Match, error) {
	_, err := r.col.InsertOne(ctx, m)
	return m, err
}

func (r *MatchRepo) List(ctx context.Context) ([]models.Match, error) {
	cur, err := r.col.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "dateTime", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.Match
	for cur.Next(ctx) {
		var m models.Match
		if err := cur.Decode(&m); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, cur.Err()
}

func (r *MatchRepo) FindByKey(ctx context.Context, key string) (models.Match, bool, error) {
	var m models.Match
	err := r.col.FindOne(ctx, bson.M{"matchKey": key}).Decode(&m)
	if err == mongo.ErrNoDocuments {
		return models.Match{}, false, nil
	}
	if err != nil {
		return models.Match{}, false, err
	}
	return m, true, nil
}

func (r *MatchRepo) AddGoal(ctx context.Context, key string, g models.Goal) error {
	_, err := r.col.UpdateOne(ctx,
		bson.M{"matchKey": key, "status": models.Scheduled},
		bson.M{"$push": bson.M{"goals": g}},
	)
	return err
}

func (r *MatchRepo) AddCard(ctx context.Context, key string, c models.Card) error {
	_, err := r.col.UpdateOne(ctx,
		bson.M{"matchKey": key, "status": models.Scheduled},
		bson.M{"$push": bson.M{"cards": c}},
	)
	return err
}

func (r *MatchRepo) Finalize(ctx context.Context, key string, homeGoals, awayGoals int) error {
	_, err := r.col.UpdateOne(ctx,
		bson.M{"matchKey": key},
		bson.M{"$set": bson.M{
			"homeGoals": homeGoals,
			"awayGoals": awayGoals,
			"status":    models.Finished,
		}},
	)
	return err
}

func (r *MatchRepo) ListFinished(ctx context.Context) ([]models.Match, error) {
	cur, err := r.col.Find(ctx, bson.M{"status": models.Finished})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.Match
	for cur.Next(ctx) {
		var m models.Match
		if err := cur.Decode(&m); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, cur.Err()
}
