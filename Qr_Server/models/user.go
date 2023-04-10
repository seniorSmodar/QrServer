package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id primitive.ObjectID `bson:"_id"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	IsAdmin bool `json:"admin,omitempty"`
	IsTerminal bool `json:"terminal,omitempty"`
}