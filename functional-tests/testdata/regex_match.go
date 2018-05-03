package testdata

var RegexMatch = `{
	"data": {
		"pairs": [
			{
				"request": {
					"body": [
						{
							"matcher": "regex",
							"value": "<item field=(.*)>"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "regex match",
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
		"timeExported": "2018-05-03T15:29:45+01:00"
	}
}`

var V3RegexMatch = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "regex match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "body": {
						"regexMatch": "<item field=(.*)>"
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
