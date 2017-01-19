package functional_tests

import (
	"net/http"

	"github.com/dghubble/sling"
	. "github.com/onsi/gomega"
)

func DoRequest(r *sling.Sling) *http.Response {
	req, err := r.Request()
	Expect(err).To(BeNil())
	response, err := http.DefaultClient.Do(req)

	Expect(err).To(BeNil())
	return response
}
