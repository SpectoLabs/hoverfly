package models

import (
	"testing"
	. "github.com/onsi/gomega"
)

func TestRequestDetailsView_ConvertToRequestDetails(t *testing.T) {
	RegisterTestingT(t)

	requestDetailsView := RequestDetailsView{
		Path: "/",
		Method: "GET",
		Destination: "/",
		Scheme: "scheme",
		Query: "", Body: "",
		Headers: map[string][]string{"Content-Encoding": []string{"gzip"}}}

	requestDetails := requestDetailsView.ConvertToRequestDetails()

	Expect(requestDetails.Path).To(Equal(requestDetailsView.Path))
	Expect(requestDetails.Method).To(Equal(requestDetailsView.Method))
	Expect(requestDetails.Destination).To(Equal(requestDetailsView.Destination))
	Expect(requestDetails.Scheme).To(Equal(requestDetailsView.Scheme))
	Expect(requestDetails.Query).To(Equal(requestDetailsView.Query))
	Expect(requestDetails.Headers).To(Equal(requestDetailsView.Headers))
}