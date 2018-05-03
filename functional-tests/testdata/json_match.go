package testdata

var JsonMatch = `{
	"data": {
		"pairs": [
			{
				"request": {
					"body": [
						{
							"matcher": "json",
							"value": "{\"test\": \"data\"}"
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
		"hoverflyVersion": "v0.17.0",
		"timeExported": "2018-05-03T15:11:24+01:00"
	}
}`

var V3JsonMatch = `{
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
						"jsonMatch": "{\"test\": \"data\"}"
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
