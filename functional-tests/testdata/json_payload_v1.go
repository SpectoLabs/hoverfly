package testdata

var V5JsonPayload = `{
	"data": {
		"pairs": [
			{
				"request": {
					"destination": [
						{
							"matcher": "exact",
							"value": "v1-simulation.com"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "v1 match",
					"encodedBody": false,
					"templated": false
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/path"
						}
					],
					"method": [
						{
							"matcher": "exact",
							"value": "GET"
						}
					],
					"destination": [
						{
							"matcher": "exact",
							"value": "v1-simulation.com"
						}
					],
					"scheme": [
						{
							"matcher": "exact",
							"value": "http"
						}
					],
					"deprecatedQuery": [
						{
							"matcher": "exact",
							"value": ""
						}
					],
					"body": [
						{
							"matcher": "exact",
							"value": ""
						}
					],
					"headers": {
                        "Accept-Encoding": [
						    {
								"matcher": "glob",
								"value": "gzip"
							}
						],
					    "User-Agent": [
						    {
								"matcher": "glob",
								"value": "Go-http-client/1.1"
							}
						]
                    }
				},
				"response": {
					"status": 200,
					"body": "v1 match",
					"encodedBody": false,
					"templated": false
				}
			}
		],
		"globalActions": {
			"delays": [],
			"delaysLogNormal": []
		}
	},
	"meta": {
		"schemaVersion": "v5",
		"hoverflyVersion": "v0.17.0",
		"timeExported": "2018-05-03T15:22:00+01:00"
	}
}`

var V1JsonPayload = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "v1 match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "destination": "v1-simulation.com"
                }
            },
            {
                "response": {
                    "status": 200,
                    "body": "v1 match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "requestType": "recording",
                    "destination": "v1-simulation.com",
                    "method": "GET",
                    "scheme": "http",
                    "path": "/path",
                    "query": "",
                    "body": "",
                    "headers": {
                        "Accept-Encoding": [
						    "gzip"
						],
					    "User-Agent": [
						    "Go-http-client/1.1"
						]
                    }
                }
            }
        ],
        "globalActions": {
            "delays": []
        }
    },
    "meta": {
        "schemaVersion": "v1",
        "hoverflyVersion": "v0.10.2",
        "timeExported": "2017-02-23T12:43:48Z"
    }
}`
