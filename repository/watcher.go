package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StreamWatcherRepository interface {
	WatcherContactEvents(resumeToken *primitive.ObjectID) (*mongo.ChangeStream, error)
	WatcherUpdatesCompanyPlan(resumeToken *primitive.ObjectID) (*mongo.ChangeStream, error)
	StoreCheckPoint(key string, resumeToken primitive.ObjectID) error
	RetrieveLastResumeToken(key string) (*primitive.ObjectID, error)
}
