package testdata

var MatcherChaining = `{
	"data": {
		"pairs": [
			{
				"request": {
					"body": [
						{
							"matcher": "jsonpath",
							"value": "$.items[4]",
							"doMatching": {
								"matcher": "jsonPartial",
								"value": "{\"fruit\": \"pineapple\"}"
							}
						}
					]
				},
				"response": {
					"status": 200,
					"body": "matcher chaining",
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
