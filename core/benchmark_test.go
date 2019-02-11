package hoverfly

import (
	"encoding/json"
	"fmt"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func BenchmarkProcessRequest(b *testing.B) {

	RegisterTestingT(b)

	hoverfly := NewHoverflyWithConfiguration(&Configuration{
		Webserver: true,
		ProxyPort: "8500",
		AdminPort: "8888",
		Mode:      "simulate",
	})

	simulation := v2.SimulationViewV5{}
	_ = json.Unmarshal([]byte(`{
	"data": {
		"pairs": [{
			"response": {
				"status": 200,
				"body": "5iNe8dxWH5Ca8pZqAfEHv3rgC0SsvKNLu6o3K",
				"encodedBody": false,
				"headers": {
					"Accept-Ranges": ["bytes"],
					"Cache-Control": ["max-age=3600"],
					"Connection": ["keep-alive"],
					"Content-Type": ["text/html; charset=utf-8"],
					"Date": ["Tue, 30 May 2017 13:23:09 GMT"],
					"Etag": ["\"82b6bafbc0c4af5e3886d07802f4b62d\""],
					"Hoverfly": ["Was-Here"],
					"Last-Modified": ["Fri, 19 May 2017 10:59:12 GMT"],
					"Server": ["nginx"],
					"Strict-Transport-Security": ["max-age=31556926"],
					"Transfer-Encoding": ["chunked"],
					"Vary": ["Accept-Encoding"]
				}
			},
			"request": {
				"path": [{
					"matcher": "exact",
					"value": "/bar"
				}],
				"method": [{
					"matcher": "exact",
					"value": "GET"
				}],
				"query": {},
				"body": [{
					"matcher": "exact",
					"value": ""
				}]
			}
		}],
		"globalActions": {
			"delays": []
		}
	},
	"meta": {
		"schemaVersion": "v5",
		"hoverflyVersion": "v0.17.0",
		"timeExported": "2017-05-30T14:23:44+01:00"
	}
}`), &simulation)

	templated := v2.SimulationViewV5{}
	_ = json.Unmarshal([]byte(`{
	"data": {
		"pairs": [{
			"response": {
				"status": 200,
				"body": "{\"st\": 1,\"sid\": 418,\"tt\": \"{{ Request.Path.[0] }}\",\"gr\": 0,\"uuid\": \"{{ randomUuid }}\",\"ip\": \"127.0.0.1\",\"ua\": \"user_agent\",\"tz\": -6,\"v\": 1}",
				"encodedBody": false,
				"templated": true,
				"headers": {
					"Accept-Ranges": ["bytes"],
					"Cache-Control": ["max-age=3600"],
					"Connection": ["keep-alive"],
					"Content-Type": ["text/html; charset=utf-8"],
					"Date": ["Tue, 30 May 2017 13:23:09 GMT"],
					"Etag": ["\"82b6bafbc0c4af5e3886d07802f4b62d\""],
					"Hoverfly": ["Was-Here"],
					"Last-Modified": ["Fri, 19 May 2017 10:59:12 GMT"],
					"Server": ["nginx"],
					"Strict-Transport-Security": ["max-age=31556926"],
					"Transfer-Encoding": ["chunked"],
					"Vary": ["Accept-Encoding"]
				}
			},
			"request": {
				"path": [{
					"matcher": "exact",
					"value": "/bar"
				}],
				"method": [{
					"matcher": "exact",
					"value": "GET"
				}],
				"query": {},
				"body": [{
					"matcher": "exact",
					"value": ""
				}]
			}
		}],
		"globalActions": {
			"delays": []
		}
	},
	"meta": {
		"schemaVersion": "v5",
		"hoverflyVersion": "v0.17.0",
		"timeExported": "2017-05-30T14:23:44+01:00"
	}
}`), &templated)

	bytes, _ := ioutil.ReadFile("../testdata/large_response_body.json")
	largeResponse := v2.SimulationViewV5{}
	_ = json.Unmarshal(bytes, &largeResponse)

	benchmarks := []struct{
		name string
		simulation v2.SimulationViewV5
	} {
		{"Simple simulation", simulation},
		{"Templated simulation", templated},
		{"Large response body", largeResponse},
	}

	fmt.Println(hoverfly.StartProxy())
	time.Sleep(time.Second)
	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8500/bar", nil)
	var resp *http.Response

	for _, bm := range benchmarks {
		hoverfly.DeleteSimulation()
		hoverfly.PutSimulation(bm.simulation)

		b.Run(bm.name, func(b *testing.B) {

			for n := 0; n < b.N; n++ {
				resp = hoverfly.processRequest(request)
			}
		})
	}


	Expect(resp.StatusCode).To(Equal(200))

	hoverfly.StopProxy()
}
