package hoverfly_end_to_end_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"github.com/dghubble/sling"
	"strings"
	"net/http"
)

func TestHoverflyEndToEnd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hoverfly End To End Suite")
}

func SetHoverflyMode(mode, port string) {
	req := sling.New().Post("http://localhost:" + port + "/api/state").Body(strings.NewReader(`{"mode":"` + mode +`"}`))
	res := DoRequest(req)
	Expect(res.StatusCode).To(Equal(200))
}

func DoRequest(r *sling.Sling) (*http.Response) {
	req, err := r.Request()
	Expect(err).To(BeNil())
	response, err := http.DefaultClient.Do(req)

	Expect(err).To(BeNil())
	return response
}