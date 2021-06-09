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

	loginResponseAM struct {
		Token string `json:"token"`
		URL   string `json:"url"`
	}

	customerAM struct {
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

	customerResponseAM struct {
		CustomerID string `json:"customer_uid"`
		ApiToken   string `json:"api_token"`
	}

	contactAM struct {
		Company   string `json:"company"`
		State     string `json:"state"`
		City      string `json:"city"`
		Zip       string `json:"zip"`
		Phone     string `json:"phone"`
		Address1  string `json:"address_1"`
		CountryID int    `json:"country_id"`
		Email     string `json:"email"`
	}

	listAM struct {
		ApiToken                string    `json:"api_token"`
		Name                    string    `json:"name"`
		SubscribeConfirmation   int       `json:"subscribe_confirmation"`
		SendWelcomeEmail        int       `json:"send_welcome_email"`
		UnsubscribeNotification int       `json:"unsubscribe_notification"`
		FromEmail               string    `json:"from_email"`
		FromName                string    `json:"from_name"`
		Contact                 contactAM `json:"contact"`
	}

	listResponseAM struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		ListUID string `json:"list_uid"`
	}

	CustomFieldAC struct {
		ApiToken string `json:"api_token"`
		Type     string `json:"type"`
		Label    string `json:"label"`
		Tag      string `json:"tag"`
	}
)

const (
	PathGetLoginToken       = "/api/v1/login-token"
	PathPostNewCustomer     = "/api/v1/customers"
	PathPostAssignPlan      = "/api/v1/customers/%s/assign-plan/%s"
	PathPostCreateList      = "/api/v1/lists"
	PathPostAddField        = "/api/v1/lists/%s/add-field"
	PathGetPlans            = "/api/v1/plans"
	PathGetExistingCustomer = "/api/v1/customers/getCustomerByEmail"
	PathPostSubscriber      = "/api/v1/subscribers"
)

var (
	customFields []CustomFieldAC
)

func init() {
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

func NewAcelleMailService(logger logger.Logger) AcelleMailService {
	return AcelleMailService{
		logger: logger,
	}
}

func (ams AcelleMailService) getURL(relativeURL string, params map[string]string) (string, error) {
	u, err := url.Parse(relativeURL)
	if err != nil {
		return "", errors.Wrap(err, "get-url parse fail relative")
	}

	if len(params) > 0 {
		queryString := u.Query()
		for key, value := range params {
			queryString.Set(key, value)
		}
		u.RawQuery = queryString.Encode()
	}
	path, err := url.Parse(conf.GetString(config.AcelleMailURI))
	if err != nil {
		return "", errors.Wrap(err, "get-url parse fail url")
	}
	uri := path.ResolveReference(u)
	return uri.String(), nil
}

func (ams AcelleMailService) httpPost(uri string, body interface{}, result interface{}) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "request body")
	}

	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrap(err, "fail post request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		message := fmt.Sprintf("fail request post, status: %d, response: %s", resp.StatusCode, respBody)
		return errors.New(message)
	}

	if result == nil {
		return nil
	}

	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return errors.Wrap(err, "fail decoding post request response")
	}
	return nil
}

func (ams AcelleMailService) httpGet(relativeURL string, params map[string]string, results interface{}) error {
	uri, err := ams.getURL(relativeURL, params)
	if err != nil {
		return err
	}
	resp, err := http.Get(uri)
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

func (ams AcelleMailService) LoginPath(information model.Company) (model.ResponseIntegration, error) {

	uri, err := ams.getURL(PathGetLoginToken, nil)
	if err != nil {
		return model.ResponseIntegration{}, errors.Wrap(err, "acelle-mail get login token")
	}

	var requestResponse loginResponseAM
	requestBody := map[string]string{
		"api_token": information.ApiToken,
	}
	if err = ams.httpPost(uri, requestBody, &requestResponse); err != nil {
		return model.ResponseIntegration{}, errors.Wrap(err, "acelle-mail get login token")
	}

	return model.ResponseIntegration{URI: requestResponse.URL}, nil
}

func (ams AcelleMailService) NewCustomer(company model.Company) (string, string, error) {

	uri, err := ams.getURL(PathPostNewCustomer, nil)
	if err != nil {
		return "", "", errors.Wrap(err, "acelle-mail post new customer")
	}

	var requestResponse customerResponseAM
	bodyRequest := customerAM{
		ApiToken: conf.GetString(config.AcelleMailToken),
		Email:    company.Email,
		//Email:         "test+13@test.com", //TODO: THIS IS FOR TESTS
		FirstName:     company.Name,
		LastName:      "-",
		TimeZone:      "America/Argentina/Buenos_Aires", //TODO desharcodear
		Language:      2,                                //TODO 1 English, 2 Spanish, 3 Portuguese 4 German
		Password:      uuid.New().String(),              //genero pass random total no lo va a usar nunca
		ColorScheme:   "white",                          // Cliengo colors scheme
		TextDirection: company.ID.String(),              // para dejar el companyId en algún lado
	}

	if err := ams.httpPost(uri, bodyRequest, &requestResponse); err != nil {
		return "", "", errors.Wrap(err, "acelle-mail post new customer")
	}

	return requestResponse.ApiToken, requestResponse.CustomerID, nil
}

func (ams AcelleMailService) NewMainList(company model.Company) (string, error) {

	uri, err := ams.getURL(PathPostCreateList, nil)
	if err != nil {
		return "", err
	}

	body := listAM{
		ApiToken:                company.ApiToken,
		Name:                    "Clientes Cliengo",
		SubscribeConfirmation:   0,
		SendWelcomeEmail:        0,
		UnsubscribeNotification: 0,
		FromEmail:               "automation@cliengomail.com",
		FromName:                "Automation",
		Contact: contactAM{
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

	var response listResponseAM
	if err := ams.httpPost(uri, body, &response); err != nil {
		return "", errors.Wrap(err, "acelle-mail post main list")
	}
	return response.ListUID, nil
}

func (ams AcelleMailService) NewCustomFields(company model.Company) error {
	uri, err := ams.getURL(fmt.Sprintf(PathPostAddField, company.MainListID), nil)
	if err != nil {
		return errors.Wrap(err, "acelle-mail post main list")
	}

	fields := customFields // se hace esta copia para no estar modificando el listado base
	for index := range fields {
		fields[index].ApiToken = company.ApiToken
		if err := ams.httpPost(uri, fields[index], nil); err != nil {
			return errors.Wrap(err, "acelle-mail post custom fields")
		}
	}
	return nil
}

func (ams AcelleMailService) AssignCustomerPlan(company model.Company, planUID string) error {
	uri, err := ams.getURL(fmt.Sprintf(PathPostAssignPlan, company.AccountID, planUID), nil)
	if err != nil {
		return errors.Wrap(err, "acelle-mail post assign plan")
	}
	body := map[string]string{
		"api_token": conf.GetString(config.AcelleMailToken),
		"gateway":   "direct",
	}

	if err := ams.httpPost(uri, body, nil); err != nil {
		return errors.Wrap(err, "acelle-mail post assign plan")
	}
	return nil
}

func (ams AcelleMailService) RetrievePlans() ([]model.AcelleMailPlan, error) {
	results := make([]model.AcelleMailPlan, 0)
	queryParams := map[string]string{
		"api_token": conf.GetString(config.AcelleMailToken),
	}
	if err := ams.httpGet(PathGetPlans, queryParams, &results); err != nil {
		return nil, errors.Wrap(err, "acelle-mail get plans")
	}
	return results, nil
}

func (ams AcelleMailService) SendSubscriber(company model.Company, contact model.Contact) error {
	if company.MainListID == "" {
		return errors.New("acelle-mail post subscriber. Not found main list id")
	}

	queryParam := map[string]string{
		"list_uid": company.MainListID,
	}
	uri, err := ams.getURL(PathPostSubscriber, queryParam)

	if err != nil {
		return nil
	}

	var result interface{}
	contact.ApiToken = company.ApiToken

	if err := ams.httpPost(uri, contact, &result); err != nil {
		return errors.Wrap(err, "acelle-mail post subscriber")
	}
	return nil
}
