package repository

import (
	"github.com/Cliengo/acelle-mail/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IntegrationRepository interface {
	RetrieveIntegrationInfo(companyID string) (model.Company, error)
	StoreNewAccountIntegration(companyID primitive.ObjectID, customerID string, apiToken string) error
	StoreAccountMainList(companyID primitive.ObjectID, listUID string) error
	//StoreTokenIntegration(companyID string, token string) error
	//GetApiKey(companyID string) (string, error)
}
