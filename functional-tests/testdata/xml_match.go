package testdata

var XmlMatch = `{
	"data": {
		"pairs": [
			{
				"request": {
					"body": [
						{
							"matcher": "xml",
							"value": "<items><item>one</item></items>"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "xml match",
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
		"timeExported": "2018-05-03T14:45:16+01:00"
	}
}`

var V3XmlMatch = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "xml match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "body": {
						"xmlMatch": "<items><item>one</item></items>"
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
