package testdata

var IncorrectContentLength = `{
	"data": {
		"pairs": [
			{
				"request": {
					"destination": [
						{
							"matcher": "exact",
							"value": "hoverfly.io"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "json match",
					"encodedBody": false,
					"templated": false,
					"headers": {
						"Content-Length": ["10000000000"]
					}
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
