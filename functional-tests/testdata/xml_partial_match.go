package testdata

var XmlPartialMatch = `{
	"data": {
		"pairs": [
			{
				"request": {
					"body": [
						{
							"matcher": "xmlpartial",
							"value": "<items><item>123</item><item>{{regex:^[A-Z]\\d{5}$}}</item><item>{{ignore}}</item></items>"
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

