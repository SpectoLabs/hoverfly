package testdata

var ClosestMissProof = `{
	"data": {
		"pairs": [{
				"response": {
					"status": 200
				},
				"request": {
					"destination": {
						"exactMatch": "destination.com"
					},
					"body": {
						"exactMatch": "body"
					}
				}
			},
			{
				"response": {
					"status": 200
				},
				"request": {
					"destination": {
						"exactMatch": "destination.com"
					},
					"body": {
						"exactMatch": "body"
					},
					"path": {
						"exactMatch": "/closest-miss"
					}
				}
			},
			{
				"response": {
					"status": 200
				},
				"request": {
					"destination": {
						"exactMatch": "destination.com"
					},
					"body": {
						"exactMatch": "body"
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
