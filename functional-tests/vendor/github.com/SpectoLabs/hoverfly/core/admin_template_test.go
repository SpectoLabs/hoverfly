package hoverfly

import (
	"bytes"
	"encoding/json"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/views"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllTemplates(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestMatcher.TemplateStore.Wipe()
	m := getBoneRouter(dbClient)

	req, err := http.NewRequest("GET", "/api/templates", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	Expect(respRec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(respRec.Body)

	pair := views.RequestResponsePairPayload{}
	err = json.Unmarshal(body, &pair)

	Expect(pair.Data).To(HaveLen(0))
}

func TestGetAllTemplatesWTemplates(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestMatcher.TemplateStore.Wipe()

	response := models.ResponseDetails{
		Body: "test-body",
	}
	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}
	destination := "testhost.com"
	method := "GET"
	path := "/a/1"
	query := "q=test"
	templateEntry := matching.RequestTemplateResponsePair{
		RequestTemplate: matching.RequestTemplate{
			Headers:     headers,
			Destination: &destination,
			Path:        &path,
			Method:      &method,
			Query:       &query,
		},
		Response: response,
	}
	dbClient.RequestMatcher.TemplateStore = matching.RequestTemplateStore{templateEntry, templateEntry}

	// performing query
	m := getBoneRouter(dbClient)

	req, err := http.NewRequest("GET", "/api/templates", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	Expect(respRec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(respRec.Body)

	pair := matching.RequestTemplateResponsePairPayload{}
	err = json.Unmarshal(body, &pair)

	// check the json given is correct to construct the request template store
	result := pair.ConvertToRequestTemplateStore()

	Expect(result).To(HaveLen(2))
}

func TestExportImportTemplates(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestMatcher.TemplateStore.Wipe()
	m := getBoneRouter(dbClient)

	// inserting some payloads
	response := models.ResponseDetails{
		Body: "test-body",
	}
	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}
	destination := "testhost.com"
	method := "GET"
	path := "/a/1"
	query := "q=test"
	templateEntry := matching.RequestTemplateResponsePair{
		RequestTemplate: matching.RequestTemplate{
			Headers:     headers,
			Destination: &destination,
			Path:        &path,
			Method:      &method,
			Query:       &query,
		},
		Response: response,
	}
	dbClient.RequestMatcher.TemplateStore = matching.RequestTemplateStore{templateEntry, templateEntry}

	req, err := http.NewRequest("GET", "/api/templates", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	Expect(respRec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(respRec.Body)

	// deleting records
	dbClient.RequestMatcher.TemplateStore.Wipe()
	Expect(dbClient.RequestMatcher.TemplateStore).To(HaveLen(0))

	// using body to import records again
	importReq, err := http.NewRequest("POST", "/api/templates", ioutil.NopCloser(bytes.NewBuffer(body)))
	//The response recorder used to record HTTP responses
	importRec := httptest.NewRecorder()

	m.ServeHTTP(importRec, importReq)
	Expect(importRec.Code).To(Equal(http.StatusOK))

	// records should be there
	Expect(dbClient.RequestMatcher.TemplateStore).To(HaveLen(2))
}

func TestDeleteTemplates(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestMatcher.TemplateStore.Wipe()
	m := getBoneRouter(dbClient)

	// inserting some payloads
	response := models.ResponseDetails{
		Body: "test-body",
	}
	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}
	destination := "testhost.com"
	method := "GET"
	path := "/a/1"
	query := "q=test"
	templateEntry := matching.RequestTemplateResponsePair{
		RequestTemplate: matching.RequestTemplate{
			Headers:     headers,
			Destination: &destination,
			Path:        &path,
			Method:      &method,
			Query:       &query,
		},
		Response: response,
	}
	dbClient.RequestMatcher.TemplateStore = matching.RequestTemplateStore{templateEntry, templateEntry}

	// checking whether we have records
	Expect(dbClient.RequestMatcher.TemplateStore).To(HaveLen(2))

	// deleting through handler
	deleteReq, _ := http.NewRequest("DELETE", "/api/templates", nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, deleteReq)
	Expect(rec.Code).To(Equal(http.StatusOK))
	Expect(dbClient.RequestMatcher.TemplateStore).To(HaveLen(0))

}
