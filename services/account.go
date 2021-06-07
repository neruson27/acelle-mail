package services

import (
	"github.com/Cliengo/acelle-mail/container/logger"
	"github.com/Cliengo/acelle-mail/model"
	"github.com/Cliengo/acelle-mail/repository"
	"github.com/pkg/errors"
	"strings"
)

type AccountService struct {
	logger      logger.Logger
	repository  repository.IntegrationRepository
	integration AcelleMailService
}

func NewAccountService(logger logger.Logger, integrationRepository repository.IntegrationRepository, integration AcelleMailService) AccountService {
	return AccountService{
		logger:      logger,
		repository:  integrationRepository,
		integration: integration,
	}
}

func (ac AccountService) GetIntegrationInfo(companyID string) (model.Company, error) {
	return ac.repository.RetrieveIntegrationInfo(companyID)
}

func (ac AccountService) MakeCompanyIntegration(company model.Company) error {
	if company.ApiToken != "" {
		// Already have a integration
		return nil
	}

	apiToken, accountID, err := ac.integration.CreateAccount(company)
	if err != nil {
		return errors.Wrap(err, "account service")
	}

	if err = ac.repository.StoreNewAccountIntegration(company.ID, accountID, apiToken); err != nil {
		return errors.Wrap(err, "account service")
	}

	company.ApiToken = apiToken
	company.AccountID = accountID

	if err = ac.setCompanyPlan(company); err != nil {
		return errors.Wrap(err, "account service")
	}

	mainListUID, err := ac.integration.CreateMainList(company)
	if err != nil {
		return errors.Wrap(err, "account service")
	}

	company.MainListID = mainListUID
	if err = ac.repository.StoreAccountMainList(company.ID, mainListUID); err != nil {
		return errors.Wrap(err, "account service")
	}

	if err = ac.integration.CreateFieldsMainList(company); err != nil {
		return errors.Wrap(err, "account service")
	}

	return nil
}

func (ac AccountService) setCompanyPlan(company model.Company) error {
	plans, err := ac.integration.GetPlans()
	if err != nil {
		return errors.Wrap(err, "service")
	}
	plan := strings.ToLower(company.Plan)
	suffix := "_annual"
	if strings.HasSuffix(plan, suffix) {
		plan = plan[:len(plan)-len(suffix)]
	}

	var uidPlan string
	for _, amPlan := range plans {
		if plan == strings.ToLower(amPlan.Name) {
			uidPlan = amPlan.UID
			break
		}
	}

	if uidPlan == "" {
		return errors.New("not found any plan to set")
	}

	if err = ac.integration.AssignPlan(company, uidPlan); err != nil {
		return errors.New("service fail")
	}
	return nil
}
