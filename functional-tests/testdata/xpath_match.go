package testdata

var XpathMatch = `{
	"data": {
		"pairs": [
			{
				"request": {
					"body": [
						{
							"matcher": "xpath",
							"value": "//item[count(preceding::item) < 5]"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "xpath match",
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
		"timeExported": "2018-05-03T14:42:34+01:00"
	}
}`

var V3XpathMatch = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "xpath match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "body": {
						"xpathMatch": "//item[count(preceding::item) < 5]"
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
