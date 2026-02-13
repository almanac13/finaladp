package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	Name         string             `bson:"name" json:"name"`
	PasswordHash string             `bson:"passwordHash" json:"-"`
	Role         string             `bson:"role" json:"role"`

	FavoriteLeague   string `bson:"favoriteLeague" json:"favoriteLeague"`
	FavoriteTeamCode string `bson:"favoriteTeamCode" json:"favoriteTeamCode"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
