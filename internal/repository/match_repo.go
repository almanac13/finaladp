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

// Universal event insert (goal increments score)
func (r *MatchRepo) AddEvent(ctx context.Context, key string, match models.Match, e models.MatchEvent) error {
	update := bson.M{
		"$push": bson.M{"events": e},
	}

	// If goal -> auto score update
	if e.Type == "goal" {
		if e.TeamCode == match.HomeCode {
			update["$inc"] = bson.M{"homeGoals": 1}
		} else if e.TeamCode == match.AwayCode {
			update["$inc"] = bson.M{"awayGoals": 1}
		}
	}

	_, err := r.col.UpdateOne(ctx,
		bson.M{"matchKey": key, "status": bson.M{"$in": []models.MatchStatus{models.Scheduled, models.Live}}},
		update,
	)
	return err
}
func (r *MatchRepo) SetStatus(ctx context.Context, key string, status models.MatchStatus) error {
	_, err := r.col.UpdateOne(ctx,
		bson.M{"matchKey": key},
		bson.M{"$set": bson.M{"status": status}},
	)
	return err
}

func (r *MatchRepo) Finalize(ctx context.Context, key string) error {
	_, err := r.col.UpdateOne(ctx,
		bson.M{"matchKey": key},
		bson.M{"$set": bson.M{"status": models.Finished}},
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
