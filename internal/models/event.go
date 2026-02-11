package models

import "time"

type EventLog struct {
	Type      string    `bson:"type" json:"type"` // match_created, goal_added, finalized...
	Message   string    `bson:"message" json:"message"`
	MatchKey  string    `bson:"matchKey,omitempty" json:"matchKey,omitempty"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
