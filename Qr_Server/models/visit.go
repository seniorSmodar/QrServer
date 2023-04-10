package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Visit struct {
	Id primitive.ObjectID `bson:"_id"`
	Username string `json:"username" validate:"required"`
	VisitDate time.Time `json:"visit_date"`
}