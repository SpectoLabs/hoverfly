package v2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
)

var csvDataSource = CSVDataSourceView{
	"Test-CSV",
	"id,name,marks\n1,Test1,55\n2,Test2,56\n",
}

const path = "/api/v2/hoverfly/templating-data-source/csv"

type HoverflyTemplateDataSourceStub struct{}

func (HoverflyTemplateDataSourceStub) GetAllDataSources() TemplateDataSourceView {

	csvDataSources := []CSVDataSourceView{csvDataSource}
	return TemplateDataSourceView{DataSources: csvDataSources}
}

func (HoverflyTemplateDataSourceStub) SetCsvDataSource(string, string) error {
	return nil
}

func (HoverflyTemplateDataSourceStub) DeleteDataSource(string) {
}

func Test_TemplateDataSourceHandler_DeleteDataSource(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyTemplateDataSourceStub{}
	unit := HoverflyTemplateDataSourceHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("DELETE", path+"/test-source", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Delete, request)
	Expect(response.Code).To(Equal(http.StatusOK))
}

func Test_TemplateDataSourceHandler_SetDataSource(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyTemplateDataSourceStub{}
	unit := HoverflyTemplateDataSourceHandler{Hoverfly: stubHoverfly}

	bodyBytes, err := json.Marshal(csvDataSource)
	Expect(err).To(BeNil())

	request, err := http.NewRequest("PUT", path, ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)
	Expect(response.Code).To(Equal(http.StatusOK))
}

func Test_TemplateDataSourceHandler_GetAllDataSources(t *testing.T) {
	RegisterTestingT(t)
	stubHoverfly := &HoverflyTemplateDataSourceStub{}
	unit := HoverflyTemplateDataSourceHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", path, nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)
	Expect(response.Code).To(Equal(http.StatusOK))

	responseBody := response.Body.String()
	expectedResponseBodyBytes, _ := json.Marshal(TemplateDataSourceView{DataSources: []CSVDataSourceView{csvDataSource}})

	Expect(responseBody).To(Equal(string(expectedResponseBodyBytes)))

}
