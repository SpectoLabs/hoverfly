package testdata

var TemplatingEnabled = `{
	"data": {
		"pairs": [
			{
				"request": {
					"method": [
						{
							"matcher": "exact",
							"value": "GET"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.QueryParam.one }}",
					"encodedBody": false,
					"templated": true
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
		"timeExported": "2018-05-03T15:44:32+01:00"
	}
}`

var V3TemplatingEnabled = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "{{ Request.QueryParam.one }}",
                    "encodedBody": false,
                    "templated" : true
                },
                "request": {
                    "method": {
						"exactMatch": "GET"
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

var JournalTemplatingWithQueryParamIndexEnabled = `{
	"data": {
		"pairs": [
			{
				"request": {
					"method": [
						{
							"matcher": "exact",
							"value": "GET"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ journal 'Request.QueryParam.id' '123' 'Response' 'jsonpath' '$.name' }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"method": [
						{
							"matcher": "exact",
							"value": "GET"
						}
					],
					"path": [
						{
							"matcher": "exact",
							"value": "/checkJournalKey"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "123: {{ hasJournalKey 'Request.QueryParam.id' '123' }} 345: {{ hasJournalKey 'Request.QueryParam.id' '345' }}",
					"encodedBody": false,
					"templated": true
				}
			}
		],
		"globalActions": {
			"delays": [],
			"delaysLogNormal": []
		}
	},
	"meta": {
		"schemaVersion": "v5.2",
		"hoverflyVersion": "v1.9.3"
	}
}`

var JournalTemplatingWithBodyIndexEnabled = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "{{ journal 'Request.Body jsonpath $.id' '1234' 'Response' 'jsonpath' '$.name' }}",
                    "encodedBody": false,
                    "templated" : true
                },
                "request": {
                    "method": {
						"exactMatch": "GET"
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
