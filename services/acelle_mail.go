package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Cliengo/acelle-mail/config"
	"github.com/Cliengo/acelle-mail/container/logger"
	"github.com/Cliengo/acelle-mail/model"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	conf "github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
)

type (
	AcelleMailService struct {
		logger logger.Logger
	}

	acelleApiToken struct {
		Token string `json:"token"`
		URL   string `json:"url"`
	}

	acelleNewAccount struct {
		ApiToken      string `json:"api_token"`
		Email         string `json:"email"`
		FirstName     string `json:"first_name"`
		LastName      string `json:"last_name"`
		TimeZone      string `json:"timezone"`
		Language      int    `json:"language_id"`
		Password      string `json:"password"`
		ColorScheme   string `json:"color_scheme"`
		TextDirection string `json:"text-direction"`
	}

	acelleMainListContact struct {
		Company   string `json:"company"`
		State     string `json:"state"`
		City      string `json:"city"`
		Zip       string `json:"zip"`
		Phone     string `json:"phone"`
		Address1  string `json:"address_1"`
		CountryID int    `json:"country_id"`
		Email     string `json:"email"`
	}

	acelleMainList struct {
		ApiToken                string                `json:"api_token"`
		Name                    string                `json:"name"`
		SubscribeConfirmation   int                   `json:"subscribe_confirmation"`
		SendWelcomeEmail        int                   `json:"send_welcome_email"`
		UnsubscribeNotification int                   `json:"unsubscribe_notification"`
		FromEmail               string                `json:"from_email"`
		FromName                string                `json:"from_name"`
		Contact                 acelleMainListContact `json:"contact"`
	}

	CustomFieldAC struct {
		ApiToken string `json:"api_token"`
		Type     string `json:"type"`
		Label    string `json:"label"`
		Tag      string `json:"tag"`
	}

	responseNewAccount struct {
		CustomerID string `json:"customer_uid"`
		ApiToken   string `json:"api_token"`
	}

	responseMainList struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		ListUID string `json:"list_uid"`
	}
)

const (
	PathGetLoginToken       = "%s/api/v1/login-token"
	PathPostNewCustomer     = "%s/api/v1/customers"
	PathPostAssignPlan      = "%s/api/v1/customers/%s/assign-plan/%s"
	PathPostCreateList      = "%s/api/v1/lists"
	PathPostAddField        = "%s/api/v1/lists/%s/add-field"
	PathGetPlans            = "/api/v1/plans"
	PathGetExistingCustomer = "/api/v1/customers/getCustomerByEmail"
)

var (
	acelleMailURI string
	customFields  []CustomFieldAC
)

func init() {
	acelleMailURI = conf.GetString(config.AcelleMailURI)
	customFields = []CustomFieldAC{
		{Type: "text", Label: "Status", Tag: "STATUS"},
		{Type: "text", Label: "Sub status", Tag: "SUB_STATUS"},
		{Type: "text", Label: "Rating", Tag: "RATING"},
		{Type: "text", Label: "UTMS Source", Tag: "UTMS_SOURCE"},
		{Type: "text", Label: "UTMS Medium", Tag: "UTMS_MEDIUM"},
		{Type: "text", Label: "UTMS Campaign", Tag: "UTMS_CAMPAIGN"},
		{Type: "text", Label: "Channel", Tag: "CHANNEL"},
		{Type: "text", Label: "Assigned To Name", Tag: "ASSIGNED_TO_NAME"},
		{Type: "text", Label: "Assigned To Email", Tag: "ASSIGNED_TO_EMAIL"},
		{Type: "text", Label: "Assigned To Url", Tag: "ASSIGNED_TO_URL"},
		{Type: "text", Label: "Country", Tag: "COUNTRY"},
		{Type: "text", Label: "City", Tag: "CITY"},
	}
}

func NewAcelleMailService(log logger.Logger) AcelleMailService {
	return AcelleMailService{
		logger: log,
	}
}

func (ams AcelleMailService) retrievePath(path string, parameters ...interface{}) string {
	if len(parameters) == 0 {
		return fmt.Sprintf(path, acelleMailURI)
	}
	return fmt.Sprintf(path, acelleMailURI, parameters)
}

func (ams AcelleMailService) postRequest(path string, body interface{}, result interface{}) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "fail parsing post data")
	}
	ams.logger.Info(path)
	ams.logger.Info(string(jsonData))
	resp, err := http.Post(path, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrap(err, "something bad happen sending post request")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		message := fmt.Sprintf("fail to retrieve information from acelle mail, status: %d, response: %s", resp.StatusCode, respBody)
		return errors.New(message)
	}

	if result == nil {
		return nil
	}

	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return errors.Wrap(err, "fail to decode result")
	}
	return nil
}

func (ams AcelleMailService) getRequest(relativeURL string, params map[string]string, results interface{}) error {
	u, err := url.Parse(relativeURL)
	if err != nil {
		return errors.Wrap(err, "fail to parse url")
	}

	if len(params) > 0 {
		queryString := u.Query()
		for key, value := range params {
			queryString.Set(key, value)
		}
		u.RawQuery = queryString.Encode()
	}
	path, err := url.Parse(acelleMailURI)
	if err != nil {
		return errors.Wrap(err, "fail to parse url")
	}
	uri := path.ResolveReference(u)
	resp, err := http.Get(uri.String())
	if err != nil {
		return errors.Wrap(err, "fail to retrieve information")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		message := fmt.Sprintf("fail to retrieve information from acelle mail, status: %d, response: %s", resp.StatusCode, respBody)
		return errors.New(message)
	}
	err = json.NewDecoder(resp.Body).Decode(results)
	if err != nil {
		return errors.Wrap(err, "fail to decode result")
	}
	return nil
}

func (ams AcelleMailService) RetrieveLoginPath(information model.Company) (model.ResponseIntegration, error) {
	body := map[string]string{
		"api_token": information.ApiToken,
	}

	var result acelleApiToken
	uri := ams.retrievePath(PathGetLoginToken)
	err := ams.postRequest(uri, body, &result)
	if err != nil {
		return model.ResponseIntegration{}, errors.Wrap(err, "retrieve login path")
	}

	return model.ResponseIntegration{URI: result.URL}, nil
}

func (ams AcelleMailService) CreateAccount(company model.Company) (string, string, error) {
	data := acelleNewAccount{
		ApiToken: conf.GetString(config.AcelleMailToken),
		//Email:         company.Email,
		Email:         "test+13@test.com",
		FirstName:     company.Name,
		LastName:      "-",
		TimeZone:      "America/Argentina/Buenos_Aires", //TODO desharcodear
		Language:      2,                                //TODO 1 English, 2 Spanish, 3 Portuguese 4 German
		Password:      uuid.New().String(),              //genero pass random total no lo va a usar nunca
		ColorScheme:   "white",                          // Cliengo colors scheme
		TextDirection: company.ID.String(),              // para dejar el companyId en algún lado
	}

	var result responseNewAccount
	uri := ams.retrievePath(PathPostNewCustomer)
	if err := ams.postRequest(uri, data, &result); err != nil {
		return "", "", errors.Wrap(err, "create new acelle mail account")
	}

	return result.ApiToken, result.CustomerID, nil
}

func (ams AcelleMailService) GetExistingEmail(company model.Company) (string, string, error) {

	//var result responseNewAccount

	return "", "", nil
}

func (ams AcelleMailService) CreateMainList(company model.Company) (string, error) {
	body := acelleMainList{
		ApiToken:                company.ApiToken,
		Name:                    "Clientes Cliengo",
		SubscribeConfirmation:   0,
		SendWelcomeEmail:        0,
		UnsubscribeNotification: 0,
		FromEmail:               "automation@cliengomail.com",
		FromName:                "Automation",
		Contact: acelleMainListContact{
			Company:   "Juan Perez",
			State:     "Buenos Aires",
			City:      "Buenos Aires",
			Zip:       "1414",
			Phone:     "1122334455",
			Address1:  "address_1",
			CountryID: 9,
			Email:     "automation@cliengomail.com",
		},
	}

	path := ams.retrievePath(PathPostCreateList)
	var response responseMainList
	if err := ams.postRequest(path, body, &response); err != nil {
		return "", errors.Wrap(err, "fail to create main list")
	}
	return response.ListUID, nil
}

func (ams AcelleMailService) CreateFieldsMainList(company model.Company) error {
	fields := customFields // se hace esta copia para no estar modificando el listado base
	path := fmt.Sprintf(PathPostAddField, acelleMailURI, company.MainListID)
	for index := range fields {
		fields[index].ApiToken = company.ApiToken
		if err := ams.postRequest(path, fields[index], nil); err != nil {
			return errors.Wrap(err, "fail to create field")
		}
	}
	return nil
}

func (ams AcelleMailService) AssignPlan(company model.Company, planUID string) error {
	body := map[string]string{
		"api_token": conf.GetString(config.AcelleMailToken),
		"gateway":   "direct",
	}
	path := fmt.Sprintf(PathPostAssignPlan, acelleMailURI, company.AccountID, planUID)
	if err := ams.postRequest(path, body, nil); err != nil {
		return errors.Wrap(err, "fail to add plan to account")
	}
	return nil
}

func (ams AcelleMailService) GetPlans() ([]model.AcelleMailPlan, error) {
	results := make([]model.AcelleMailPlan, 0)
	queryParams := map[string]string{
		"api_token": conf.GetString(config.AcelleMailToken),
	}
	if err := ams.getRequest(PathGetPlans, queryParams, &results); err != nil {
		return nil, errors.Wrap(err, "fail to retrieve acelle mail plans")
	}
	return results, nil
}
