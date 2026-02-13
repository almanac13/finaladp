package models

type Team struct {
	Code   string `bson:"_id" json:"code"` // use Code as _id (ARS, LIV...)
	Name   string `bson:"name" json:"name"`
	League string `bson:"league" json:"league"` // EPL, LaLiga, SerieA, KPL
}
