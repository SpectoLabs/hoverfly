package testdata

var SingleRequestMatcherToResponse = `{
	"data": {
		"pairs": [
			{
				"request": {
					"destination": [
						{
							"matcher": "exact",
							"value": "miss"
						}
					]
				},
				"response": {
					"status": 0,
					"body": "body",
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
		"timeExported": "2018-05-03T15:31:46+01:00"
	}
}`

var V3SingleRequestMatcherToResponse = `{
	"data": {
		"pairs": [
			{
				"response": {
					"body": "body"
				},
				"request": {
					"destination": {
						"exactMatch": "miss"
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
