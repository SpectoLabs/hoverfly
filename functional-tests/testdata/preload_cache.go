package testdata

var PreloadCache = `{
	"data": {
		"pairs": [
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/path1"
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
							"value": "test-server.com"
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
					]
				},
				"response": {
					"status": 200,
					"body": "exact match",
					"encodedBody": false,
					"headers": {
						"Header": [
							"value1",
							"value2"
						]
					},
					"templated": false
				}
			},
			{
				"request": {
					"destination": [
						{
							"matcher": "exact",
							"value": "destination-server.com"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "destination matched",
					"encodedBody": false,
					"templated": false
				}
			}
		],
		"globalActions": {
			"delays": []
		}
	},
	"meta": {
		"schemaVersion": "v5",
		"hoverflyVersion": "v0.17.0",
		"timeExported": "2018-05-03T15:28:40+01:00"
	}
}`

var V3PreloadCache = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "exact match",
                    "encodedBody": false,
                    "headers": {
                        "Header": [
                            "value1",
                            "value2"
                        ]
                    }
                },
                "request": {
                    "path": {
						"exactMatch": "/path1"
                    },
                    "method": {
						"exactMatch": "GET"
                    },
                    "destination": {
						"exactMatch": "test-server.com"
                    },
                    "scheme": {
						"exactMatch": "http"
                    },
                    "query": {
						"exactMatch": ""
                    },
                    "body": {
						"exactMatch": ""
                    }
                }
            },
            {
                "response": {
                    "status": 200,
                    "body": "destination matched",
                    "encodedBody": false
                },
                "request": {
                    "destination": {
                        "exactMatch": "destination-server.com"
                    }
                }
            }
        ],
        "globalActions": {
            "delays": []
        }
    },
    "meta": {
        "schemaVersion": "v3",
        "hoverflyVersion": "v0.10.2",
        "timeExported": "2017-02-23T12:43:48Z"
    }
}`
