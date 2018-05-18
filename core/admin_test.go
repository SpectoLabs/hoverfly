package hoverfly

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-zoo/bone"
	. "github.com/onsi/gomega"
)

var adminApi = AdminApi{}

func Test_AdminAPI_Serves404(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	router := bone.New()
	m := adminApi.addAdminApiRoutes(router, unit)

	req, err := http.NewRequest("GET", "/api/sdiughvksjv", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	Expect(respRec.Code, http.StatusNotFound)
}
