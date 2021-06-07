package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WatcherRepository struct {
	companyColl *mongo.Collection
	contactColl *mongo.Collection
}

func NewMongoWatcherRepository(client Client) WatcherRepository {
	return WatcherRepository{
		companyColl: client.database.Collection(companyCollection),
		contactColl: client.database.Collection(contactCollection),
	}
}

func (wr WatcherRepository) WatcherContactEvents(resumeToken *primitive.ObjectID) (*mongo.ChangeStream, error) {
	pipeline := bson.D{{
		"$match", bson.D{
			{"operationType", bson.D{{"$in", bson.A{"insert", "update"}}}},
		},
	}}
	opts := []*options.ChangeStreamOptions{
		options.ChangeStream().SetFullDocument(options.UpdateLookup),
	}
	if resumeToken != nil {
		opts = append(opts, options.ChangeStream().SetResumeAfter(resumeToken))
	}
	return wr.contactColl.Watch(context.TODO(), mongo.Pipeline{pipeline}, opts...)
}

func (wr WatcherRepository) WatcherUpdatesCompanyPlan(resumeToken *primitive.ObjectID) (*mongo.ChangeStream, error) {
	pipeline := bson.D{{
		"$match", bson.D{
			{"operationType", "update"},
			{"updateDescription.updatedFields.planType", bson.D{{"$exists", true}}},
		},
	}}
	opts := []*options.ChangeStreamOptions{
		options.ChangeStream().SetFullDocument(options.UpdateLookup),
	}
	if resumeToken != nil {
		opts = append(opts, options.ChangeStream().SetResumeAfter(resumeToken))
	}
	return wr.companyColl.Watch(context.TODO(), mongo.Pipeline{pipeline}, opts...)
}

func (wr WatcherRepository) StoreCheckPoint(key string, resumeToken primitive.ObjectID) error {
	//TODO: Make implementation
	return nil
}

func (wr WatcherRepository) RetrieveLastResumeToken(key string) (*primitive.ObjectID, error) {
	return nil, nil
}
