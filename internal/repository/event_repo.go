package repository

import (
	"context"

	"final-by-me/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type EventRepo struct {
	col *mongo.Collection
}

func NewEventRepo(db *mongo.Database) *EventRepo {
	return &EventRepo{col: db.Collection("events")}
}

func (r *EventRepo) Insert(ctx context.Context, e models.EventLog) error {
	_, err := r.col.InsertOne(ctx, e)
	return err
}
