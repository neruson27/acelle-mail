package services

import (
	"github.com/Cliengo/acelle-mail/container/logger"
	"github.com/Cliengo/acelle-mail/model"
	"github.com/Cliengo/acelle-mail/repository"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StreamsService struct {
	logger      logger.Logger
	repository  repository.IntegrationRepository
	integration AcelleMailService
}

func NewStreamService(logger logger.Logger, repository repository.IntegrationRepository, integration AcelleMailService) StreamsService {
	return StreamsService{
		logger:      logger,
		repository:  repository,
		integration: integration,
	}
}

func (ss StreamsService) parseToStruct(stream bson.M, result interface{}) error {
	bsonBytes, err := bson.Marshal(stream)
	if err != nil {
		return errors.Wrap(err, "fail to parse stream to result struct")
	}
	bson.Unmarshal(bsonBytes, result)
	return nil
}

func (ss StreamsService) ProcessUpdateCompany(stream bson.M) (primitive.ObjectID, error) {
	ss.logger.Info("Update company")
	ss.logger.Info(stream)
	return primitive.ObjectID{}, nil
}

func (ss StreamsService) ProcessContactEvent(stream bson.M) (primitive.ObjectID, error) {
	var result model.StreamContact
	if err := ss.parseToStruct(stream, &result); err != nil {
		return primitive.ObjectID{}, err
	}

	ss.logger.Infof("Operation: %s, Information: %s", result.OperationType, result.FullDocument)
	return primitive.ObjectID{}, nil
}
