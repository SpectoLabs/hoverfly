package testdata

var JwtMatcher = `{
	"data": {
		"pairs": [
			{
				"request": {
					"headers": {
						"Authorisation": [
							{
								"matcher": "jwt",
								"value": "{\"header\":{\"alg\":\"HS256\"},\"payload\":{\"sub\":\"1234567890\",\"name\":\"John Doe\"}}"
							}
						]
					}
				},
				"response": {
					"status": 200,
					"body": "jwt matchers matches",
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
		"schemaVersion": "v5.2",
		"hoverflyVersion": "v0.17.0",
		"timeExported": "2018-05-03T14:36:59+01:00"
	}
}`
