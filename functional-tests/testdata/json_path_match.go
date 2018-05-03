package testdata

var JsonPathMatch = `{
	"data": {
		"pairs": [
			{
				"request": {
					"body": [
						{
							"matcher": "jsonpath",
							"value": "$.items[4]"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "json match",
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
		"hoverflyVersion": "v0.16.0",
		"timeExported": "2018-05-03T15:19:58+01:00"
	}
}`

var V3JsonPathMatch = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "json match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "body": {
						"jsonPathMatch": "$.items[4]"
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
