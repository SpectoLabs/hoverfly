package hoverfly_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//"io/ioutil"
	"github.com/SpectoLabs/hoverfly/models"
	"github.com/dghubble/sling"
	"strings"
	"io/ioutil"
)

var _ = Describe("Interacting with the API", func() {

	BeforeEach(func() {
		requestCache.DeleteData()
		pl1 := models.Payload{
			Request: models.RequestDetails{
				Path:"path1",
				Method:"method1",
				Destination:"destination1",
				Scheme:"scheme1",
				Query:"query1",
				Body:"body1",
				Headers:map[string][]string{"header": []string{"value1"}},
			},
			Response: models.ResponseDetails{
				Status: 201,
				Body: "body1",
				Headers:map[string][]string{"header": []string{"value1"}},
			},
		}
		encoded, _ := pl1.Encode()
		requestCache.Set([]byte(pl1.Id()), encoded)
		pl2 := models.Payload{
			Request: models.RequestDetails{
				Path:"path2",
				Method:"method2",
				Destination:"destination2",
				Scheme:"scheme2",
				Query:"query2",
				Body:"body2",
				Headers:map[string][]string{"header": []string{"value2"}},
			},
			Response: models.ResponseDetails{
				Status: 202,
				Body: "body2",
				Headers:map[string][]string{"header": []string{"value2"}},
			},
		}
		encoded, _ = pl2.Encode()
		requestCache.Set([]byte(pl2.Id()), encoded)
	})

	Context("GET /api/records", func() {

		It("Should retrieve the records", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/records")
			res := DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			recordsJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(recordsJson).To(MatchJSON(
				`{
				  "data": [
				    {
				      "response": {
					"status": 201,
					"body": "body1",
					"encodedBody": false,
					"headers": {
					  "header": [
					    "value1"
					  ]
					}
				      },
				      "request": {
					"path": "path1",
					"method": "method1",
					"destination": "destination1",
					"scheme": "scheme1",
					"query": "query1",
					"body": "body1",
					"headers": {
					  "header": [
					    "value1"
					  ]
					}
				      }
				    },
				    {
				      "response": {
					"status": 202,
					"body": "body2",
					"encodedBody": false,
					"headers": {
					  "header": [
					    "value2"
					  ]
					}
				      },
				      "request": {
					"path": "path2",
					"method": "method2",
					"destination": "destination2",
					"scheme": "scheme2",
					"query": "query2",
					"body": "body2",
					"headers": {
					  "header": [
					    "value2"
					  ]
					}
				      }
				    }
				  ]
				}`))
		})
	})

	Context("DELETE /api/records", func() {

		It("Should delete the records", func() {
			req := sling.New().Delete(hoverflyAdminUrl + "/api/records")
			res := DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			Expect(requestCache.RecordsCount()).To(Equal(0))
		})
	})

	Context("POST /api/records", func() {

		Context("When no records exist", func() {
			It("Should create the records", func() {
				res := DoRequest(sling.New().Post(hoverflyAdminUrl + "/api/records").Set("Content-Type", "application/json").Body(
					strings.NewReader(`
					{
						"data": [{
							"response": {
								"status": 201,
								"body": "body1",
								"encodedBody": false,
								"headers": {
									"header": [
										"value1"
									]
								}
							},
							"request": {
								"path": "path1",
								"method": "method1",
								"destination": "destination1",
								"scheme": "scheme1",
								"query": "query1",
								"body": "body1",
								"headers": {
									"header": [
										"value1"
									],
									"Content-Type": [
										"application/json"
									]
								}
							}
						}]
					}
					`)))
				Expect(res.StatusCode).To(Equal(200))

				data := models.PayloadViewData{
					Data: []models.PayloadView{models.PayloadView{
						Request: models.RequestDetailsView{
							Path:"path1",
							Method:"method1",
							Destination:"destination1",
							Scheme:"scheme1",
							Query:"query1",
							Body:"body1",
							Headers:map[string][]string{
								"header": []string{"value1"},
								"Content-Type": []string{"application/json"},
							},
						},
						Response: models.ResponseDetailsView{
							Status: 201,
							Body: "body1",
							EncodedBody: false,
							Headers:map[string][]string{"header": []string{"value1"}},
						},
					}},
				}

				expectedPayload := data.Data[0].ConvertToPayload()
				outputBytes, err := requestCache.Get([]byte(expectedPayload.Id()))
				Expect(err).To(BeNil())
				payload, err := models.NewPayloadFromBytes(outputBytes)
				Expect(*payload).To(Equal(expectedPayload))
			})
		})

		Context("When a record already exists", func() {

			BeforeEach(func() {
				p := models.Payload{
					Request: models.RequestDetails{
						Path:"path2",
						Method:"method2",
						Destination:"destination2",
						Scheme:"scheme2",
						Query:"query2",
						Body:"body2",
						Headers:map[string][]string{"header": []string{"value1"}},
					},
					Response: models.ResponseDetails{
						Status: 201,
						Body: "body2",
						Headers:map[string][]string{"header": []string{"value1"}},
					},
				}
				bytes, err := p.Encode()
				Expect(err).To(BeNil())
				requestCache.Set([]byte(p.Id()), bytes)
			})

			It("Should append the records to the existing ones", func() {
				res := DoRequest(sling.New().Post(hoverflyAdminUrl + "/api/records").Set("Content-Type", "application/json").Body(
					strings.NewReader(`
					{
						"data": [{
							"response": {
								"status": 201,
								"body": "body1",
								"encodedBody": false,
								"headers": {
									"header": [
										"value1"
									]
								}
							},
							"request": {
								"path": "path1",
								"method": "method1",
								"destination": "destination1",
								"scheme": "scheme1",
								"query": "query1",
								"body": "body1",
								"headers": {
									"header": [
										"value1"
									],
									"Content-Type": [
										"application/json"
									]
								}
							}
						}]
					}
					`)))
				Expect(res.StatusCode).To(Equal(200))

				data := models.PayloadViewData{
					Data: []models.PayloadView{models.PayloadView{
						Request: models.RequestDetailsView{
							Path:"path1",
							Method:"method1",
							Destination:"destination1",
							Scheme:"scheme1",
							Query:"query1",
							Body:"body1",
							Headers:map[string][]string{
								"header": []string{"value1"},
								"Content-Type": []string{"application/json"},
							},
						},
						Response: models.ResponseDetailsView{
							Status: 201,
							Body: "body1",
							EncodedBody: false,
							Headers:map[string][]string{"header": []string{"value1"}},
						},
					}},
				}

				expectedPayload := data.Data[0].ConvertToPayload()
				outputBytes, err := requestCache.Get([]byte(expectedPayload.Id()))
				Expect(err).To(BeNil())
				payload, err := models.NewPayloadFromBytes(outputBytes)
				Expect(*payload).To(Equal(expectedPayload))

				other := models.Payload{
					Request: models.RequestDetails{
						Path:"path2",
						Method:"method2",
						Destination:"destination2",
						Scheme:"scheme2",
						Query:"query2",
						Body:"body2",
						Headers:map[string][]string{"header": []string{"value1"}},
					},
					Response: models.ResponseDetails{
						Status: 201,
						Body: "body2",
						Headers:map[string][]string{"header": []string{"value1"}},
					},
				}

				outputBytes, err = requestCache.Get([]byte(other.Id()))
				Expect(err).To(BeNil())
				payload, err = models.NewPayloadFromBytes(outputBytes)
				Expect(*payload).To(Equal(other))
			})
		})
	})
})
