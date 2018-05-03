package testdata

var SingleRequestMatcherToResponse = `{
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
