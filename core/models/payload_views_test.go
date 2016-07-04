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

func TestPayloadViewData_ConvertToPayloadDataWithoutEncoding(t *testing.T) {
	RegisterTestingT(t)

	view := PayloadView{
		Request: RequestDetailsView{
			Path: "A",
			Method: "A",
			Destination: "A",
			Scheme: "A",
			Query: "A",
			Body: "A",
			Headers: map[string][]string{
				"A" : []string{"B"},
				"C" : []string{"D"},
			},
		},
		Response: ResponseDetailsView{
			Status: 1,
			Body: "1",
			EncodedBody: false,
			Headers: map[string][]string{
				"1" : []string{"2"},
				"3" : []string{"4"},
			},
		},
	}

	payload := view.ConvertToPayload()

	Expect(payload).To(Equal(Payload{
		Request: RequestDetails{
			Path: "A",
			Method: "A",
			Destination: "A",
			Scheme: "A",
			Query: "A",
			Body: "A",
			Headers: map[string][]string{
				"A" : []string{"B"},
				"C" : []string{"D"},
			},
		},
		Response: ResponseDetails{
			Status: 1,
			Body: "1",
			Headers: map[string][]string{
				"1" : []string{"2"},
				"3" : []string{"4"},
			},
		},
	}))
}

func TestPayloadViewData_ConvertToPayloadDataWithEncoding(t *testing.T) {
	RegisterTestingT(t)

	view := PayloadView{
		Response: ResponseDetailsView{
			Body: "ZW5jb2RlZA==",
			EncodedBody: true,
		},
	}

	payload := view.ConvertToPayload()

	Expect(payload.Response.Body).To(Equal("encoded"))
}