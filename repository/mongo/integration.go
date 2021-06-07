package mongo

import (
	"context"
	"github.com/Cliengo/acelle-mail/model"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IntegrationRepository struct {
	collection *mongo.Collection
}

func NewMongoIntegrationRepository(client Client) IntegrationRepository {
	return IntegrationRepository{
		collection: client.database.Collection(companyCollection),
	}
}

func (ir IntegrationRepository) RetrieveIntegrationInfo(companyID string) (model.Company, error) {
	if companyID == "" {
		return model.Company{}, errors.New("not valid identifier")
	}
	objID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return model.Company{}, errors.Wrap(err, "fail parsing object id from hex")
	}

	filter := bson.M{"_id": objID}
	var result model.Company
	if ir.collection.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		return model.Company{}, err
	}
	return result, nil
}

func (ir IntegrationRepository) StoreNewAccountIntegration(companyID primitive.ObjectID, customerID string, apiToken string) error {
	filter := bson.M{"_id": companyID}
	set := bson.M{"$set": bson.M{"acelleAccountId": customerID, "acelleAccountApiToken": apiToken}}
	if _, err := ir.collection.UpdateOne(context.TODO(), filter, set); err != nil {
		return errors.Wrap(err, "fail to update collection")
	}
	return nil
}

func (ir IntegrationRepository) StoreAccountMainList(companyID primitive.ObjectID, listUID string) error {
	filter := bson.M{"_id": companyID}
	set := bson.M{"$set": bson.M{"acelleAccountMainListId": listUID}}
	if _, err := ir.collection.UpdateOne(context.TODO(), filter, set); err != nil {
		return errors.Wrap(err, "fail to update collection")
	}
	return nil
}
