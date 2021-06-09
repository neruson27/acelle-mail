package services

import (
	"github.com/Cliengo/acelle-mail/container/logger"
	"github.com/Cliengo/acelle-mail/model"
	"github.com/Cliengo/acelle-mail/repository"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type StreamsService struct {
	logger      logger.Logger
	repository  repository.IntegrationRepository
	integration AcelleMailService
	account     AccountService
}

func NewStreamService(logger logger.Logger, repository repository.IntegrationRepository, integration AcelleMailService, account AccountService) StreamsService {
	return StreamsService{
		logger:      logger,
		repository:  repository,
		integration: integration,
		account:     account,
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

func (ss StreamsService) ProcessUpdateCompany(stream bson.M) error {
	var result model.StreamCompany
	if err := ss.parseToStruct(stream, &result); err != nil {
		return errors.Wrap(err, "fail to parse to struct stream")
	}
	if result.FullDocument.ApiToken == "" {
		// The company don't have integration with acelle mail
		if err := ss.account.MakeCompanyIntegration(result.FullDocument); err != nil {
			return errors.Wrap(err, "company stream processor")
		}
	} else {
		// Updates the plan in acelle mail
		if err := ss.account.setCompanyPlan(result.FullDocument); err != nil {
			return errors.Wrap(err, "company stream processor")
		}
	}

	return nil
}

func (ss StreamsService) ProcessContactEvent(stream bson.M) error {
	var result model.StreamContact
	if err := ss.parseToStruct(stream, &result); err != nil {
		return err
	}
	company, err := ss.account.GetIntegrationInfo(result.FullDocument.Company.ID)
	if err != nil {
		ss.logger.Info(err)
		return errors.Wrap(err, "contact stream processor")
	}

	if company.ApiToken == "" {
		if err := ss.account.MakeCompanyIntegration(company); err != nil {
			return errors.Wrap(err, "contact stream processor")
		}
	}

	if err = ss.integration.SendSubscriber(company, result.FullDocument); err != nil {
		return errors.Wrap(err, "contact stream processor")
	}

	return nil
}
