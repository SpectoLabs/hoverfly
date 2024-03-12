package testdata

var NegationMatch = `{
	"data": {
		"pairs": [
			{
				"request": {
					"path": [
						{
							"matcher": "negate",
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
							"value": "test.com"
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
					]
				},
				"response": {
					"status": 200,
					"body": "destination matched",
					"encodedBody": false,
					"headers": {
						"Header": [
							"value1"
						]
					},
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
		"timeExported": "2018-05-03T12:12:46+01:00"
	}
}`
