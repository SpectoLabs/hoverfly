package testdata

var Delays = `{
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
			"delays": [
				{
					"urlPattern": "test-server\\.com",
					"httpMethod": "",
					"delay": 100
				},
				{
					"urlPattern": "test-server\\.com",
					"httpMethod": "",
					"delay": 110
				},
				{
					"urlPattern": "localhost(.*)",
					"httpMethod": "",
					"delay": 110
				}
			],
			"delaysLogNormal": []
		}
	},
	"meta": {
		"schemaVersion": "v5",
		"hoverflyVersion": "v0.17.0",
		"timeExported": "2018-05-03T11:56:26+01:00"
	}
}`

var V3Delays = `{
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
                    },
                    "headers": {
                        "Accept-Encoding": [
						    "gzip"
						],
					    "User-Agent": [
						    "Go-http-client/1.1"
						]
                    }
                }
            },
            {
                "response": {
                    "status": 200,
                    "body": "destination matched",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "destination": {
                        "exactMatch": "destination-server.com"
                    }
                }
            }
        ],
        "globalActions": {
            "delays": [
                    {
                        "urlPattern": "test-server\\.com",
                        "delay": 100
                    },
                    {
                        "urlPattern": "test-server\\.com",
                        "delay": 110
                    },
                    {
                        "urlPattern": "localhost(.*)",
                        "delay": 110
                    }
            ]
        }
    },
    "meta": {
        "schemaVersion": "v3",
        "hoverflyVersion": "v0.10.2",
        "timeExported": "2017-02-23T12:43:48Z"
    }
}`
