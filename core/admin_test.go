package hoverfly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	. "github.com/onsi/gomega"
)

var adminApi = AdminApi{}

func TestSetMetadata(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	// preparing to set mode through rest api
	var reqBody v1.SetMetadata
	reqBody.Key = "some_key"
	reqBody.Value = "some_val"

	reqBodyBytes, err := json.Marshal(&reqBody)
	Expect(err).To(BeNil())

	// deleting through handler
	req, err := http.NewRequest("PUT", "/api/metadata", ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes)))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusCreated))

	// checking mode
	metaValue, err := unit.MetadataCache.Get([]byte("some_key"))
	Expect(err).To(BeNil())
	Expect(string(metaValue)).To(Equal("some_val"))
}

func TestSetMetadataBadBody(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	// deleting through handler
	req, err := http.NewRequest("PUT", "/api/metadata", ioutil.NopCloser(bytes.NewBuffer([]byte("you shall not decode me!!"))))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusBadRequest))
}

func TestSetMetadataMissingKey(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	// preparing to set mode through rest api
	var reqBody v1.SetMetadata
	// missing key
	reqBody.Value = "some_val"

	reqBodyBytes, err := json.Marshal(&reqBody)
	Expect(err).To(BeNil())

	// deleting through handler
	req, err := http.NewRequest("PUT", "/api/metadata", ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes)))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusBadRequest))

	// checking response body
	body, err := ioutil.ReadAll(rec.Body)
	mr := v1.MessageResponse{}
	err = json.Unmarshal(body, &mr)

	Expect(mr.Message).To(Equal("Key not provided."))
}

func TestGetMetadata(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)
	// adding some metadata
	for i := 0; i < 3; i++ {
		k := fmt.Sprintf("key_%d", i)
		v := fmt.Sprintf("val_%d", i)
		err := unit.MetadataCache.Set([]byte(k), []byte(v))
		Expect(err).To(BeNil())
	}

	req, err := http.NewRequest("GET", "/api/metadata", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)

	sm := v1.StoredMetadata{}
	err = json.Unmarshal(body, &sm)

	Expect(len(sm.Data)).To(Equal(3))

	for i := 0; i < 3; i++ {
		k := fmt.Sprintf("key_%d", i)
		v := fmt.Sprintf("val_%d", i)
		Expect(sm.Data[k]).To(Equal(v))
	}
}

func TestDeleteMetadata(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)
	// adding some metadata
	for i := 0; i < 3; i++ {
		k := fmt.Sprintf("key_%d", i)
		v := fmt.Sprintf("val_%d", i)
		err := unit.MetadataCache.Set([]byte(k), []byte(v))
		Expect(err).To(BeNil())
	}

	// checking that metadata is there
	allMeta, err := unit.MetadataCache.GetAllEntries()
	Expect(err).To(BeNil())
	Expect(len(allMeta)).To(Equal(3))

	// deleting it
	req, err := http.NewRequest("DELETE", "/api/metadata", nil)
	Expect(err).To(BeNil())
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	// checking metadata again, should be zero
	allMeta, err = unit.MetadataCache.GetAllEntries()
	Expect(err).To(BeNil())
	Expect(len(allMeta)).To(Equal(0))
}

func TestDeleteMetadataEmpty(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	// deleting it
	req, err := http.NewRequest("DELETE", "/api/metadata", nil)
	Expect(err).To(BeNil())
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	// checking metadata again, should be zero
	allMeta, err := unit.MetadataCache.GetAllEntries()
	Expect(err).To(BeNil())
	Expect(len(allMeta)).To(Equal(0))
}

func TestGetMissingURL(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	req, err := http.NewRequest("GET", "/api/sdiughvksjv", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	Expect(respRec.Code, http.StatusNotFound)
}
