package test

import (
	"encoding/json"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func clearCustomers() {
	// delete dependent family_list first to satisfy FK constraints
	if err := db.Where("fl_id is not null").Delete(&entity.FamilyList{}).Error; err != nil {
		log.Fatalf("Failed clear family_list data : %+v", err)
	}

	err := db.Where("cst_id is not null").Delete(&entity.Customer{}).Error
	if err != nil {
		log.Fatalf("Failed clear customer data : %+v", err)
	}
}

func ensureNationality(t *testing.T, id int) {
	n := &entity.Nationality{NationalityID: id, NationalityName: "Indonesia", NationalityCode: "ID"}
	err := db.Where("nationality_id = ?", id).FirstOrCreate(n).Error
	assert.Nil(t, err)
}

func TestCreateCustomer(t *testing.T) {
	clearCustomers()
	ensureNationality(t, 1)

	requestBody := model.CreateCustomerRequest{
		NationalityID: 1,
		CstName:       "Budi",
		CstDob:        "1990-01-01",
		CstPhoneNum:   "08123456789",
		CstEmail:      "budi@example.com",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/customers", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[*model.CustomerResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data)
	assert.Equal(t, requestBody.CstName, responseBody.Data.CstName)
	assert.Equal(t, requestBody.CstEmail, responseBody.Data.CstEmail)
}

func TestUpdateCustomer(t *testing.T) {
	// ensure a customer exists
	TestCreateCustomer(t)

	// get first customer
	cust := new(entity.Customer)
	err := db.First(cust).Error
	assert.Nil(t, err)

	requestBody := model.UpdateCustomerRequest{
		CstName:     "Budi Updated",
		CstPhoneNum: "08999999999",
		CstEmail:    "budi.updated@example.com",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	url := "/api/customers/" + strconv.Itoa(cust.CstID)
	request := httptest.NewRequest(http.MethodPut, url, strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[*model.CustomerResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data)
	assert.Equal(t, requestBody.CstName, responseBody.Data.CstName)
	assert.Equal(t, requestBody.CstEmail, responseBody.Data.CstEmail)
}

func TestCreateCustomerValidationError(t *testing.T) {
	clearCustomers()
	ensureNationality(t, 1)

	// missing required fields
	requestBody := model.CreateCustomerRequest{}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/customers", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestCreateCustomerInvalidEmail(t *testing.T) {
	clearCustomers()
	ensureNationality(t, 1)

	requestBody := model.CreateCustomerRequest{
		NationalityID: 1,
		CstName:       "Siti",
		CstDob:        "1995-05-05",
		CstPhoneNum:   "0811111111",
		CstEmail:      "not-an-email",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/customers", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestUpdateCustomerNotFound(t *testing.T) {
	clearCustomers()
	ensureNationality(t, 1)

	requestBody := model.UpdateCustomerRequest{
		CstName: "Does Not Exist",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, "/api/customers/99999", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestUpdateCustomerInvalidDob(t *testing.T) {
	// create a customer first
	TestCreateCustomer(t)

	cust := new(entity.Customer)
	err := db.First(cust).Error
	assert.Nil(t, err)

	requestBody := model.UpdateCustomerRequest{
		CstDob: "31-12-1990", // invalid format
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	url := "/api/customers/" + strconv.Itoa(cust.CstID)
	request := httptest.NewRequest(http.MethodPut, url, strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestGetCustomer(t *testing.T) {
	// create a customer first
	TestCreateCustomer(t)

	cust := new(entity.Customer)
	err := db.First(cust).Error
	assert.Nil(t, err)

	url := "/api/customers/" + strconv.Itoa(cust.CstID)
	request := httptest.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[*model.CustomerResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data)
	assert.Equal(t, cust.CstID, responseBody.Data.CstID)
}

func TestGetCustomerNotFound(t *testing.T) {
	clearCustomers()
	ensureNationality(t, 1)

	request := httptest.NewRequest(http.MethodGet, "/api/customers/999999", nil)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestDeleteCustomer(t *testing.T) {
	// create a customer first
	TestCreateCustomer(t)

	cust := new(entity.Customer)
	err := db.First(cust).Error
	assert.Nil(t, err)

	url := "/api/customers/" + strconv.Itoa(cust.CstID)
	request := httptest.NewRequest(http.MethodDelete, url, nil)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.True(t, responseBody.Data)
}

func TestDeleteCustomerNotFound(t *testing.T) {
	clearCustomers()
	ensureNationality(t, 1)

	request := httptest.NewRequest(http.MethodDelete, "/api/customers/999999", nil)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}
