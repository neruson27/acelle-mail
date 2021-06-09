package http

import (
	"encoding/json"
	"github.com/Cliengo/acelle-mail/container/logger"
	"github.com/Cliengo/acelle-mail/infrastructure/http/middlewares"
	"github.com/Cliengo/acelle-mail/services"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type IntegrationHandler struct {
	information services.AccountService
	integration services.AcelleMailService
	logger      logger.Logger
}

func NewIntegrationHandler(logger logger.Logger, information services.AccountService, integration services.AcelleMailService) IntegrationHandler {
	return IntegrationHandler{
		information: information,
		integration: integration,
		logger:      logger,
	}
}

func (ih IntegrationHandler) NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/", handler(ih.getInformation))
	return r
}

func (ih IntegrationHandler) getInformation(w http.ResponseWriter, r *http.Request) error {
	//ih.information.GetIntegrationInfo()
	companyID, exists := middlewares.GetCompanyID(r.Context())
	if !exists {
		ih.logger.Info("No se ha encontrado la companyID ")
		//TODO: Retornar un codigo 204 o 400
		return nil
	}

	information, err := ih.information.GetIntegrationInfo(companyID)
	if err != nil {
		ih.logger.Info("Se ha generado un error al obtener la informacion de la integracion ")
		ih.logger.Info(err)
		//TODO: Retornar un codigo 204 o 400
		return nil
	}
	if information.ApiToken == "" {
		if err := ih.information.MakeCompanyIntegration(information); err != nil {
			ih.logger.Info(err)
			return nil
		}
		information, _ = ih.information.GetIntegrationInfo(companyID)
	}

	resultInfo, err := ih.integration.LoginPath(information)
	if err != nil {
		//TODO: Manejar el error
		ih.logger.Info("error con acelle mail")
		ih.logger.Info(err)
		return nil
	}

	_ = json.NewEncoder(w).Encode(resultInfo)
	return nil
}
