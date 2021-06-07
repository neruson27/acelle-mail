package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Company struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name"`
	Plan       string             `bson:"planType"`
	Email      string             `bson:"email"`
	AccountID  string             `bson:"acelleAccountId"`
	ApiToken   string             `bson:"acelleAccountApiToken"`
	MainListID string             `bson:"acelleAccountMainListId"`
}
