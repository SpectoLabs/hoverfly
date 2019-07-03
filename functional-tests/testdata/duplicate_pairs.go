package testdata

var DuplicatePairs = `{
	"data": {
		"pairs": [
			{
				"request": {
					"destination": [
						{
							"matcher": "exact",
							"value": "destination.com"
						},
						{
							"matcher": "glob",
							"value": "*.com"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "multiple matches",
					"encodedBody": false,
					"templated": false
				}
			},
			{
				"request": {
					"destination": [
						{
							"matcher": "exact",
							"value": "destination.com"
						},
						{
							"matcher": "glob",
							"value": "*.com"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "multiple matches 2",
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
		"hoverflyVersion": "v0.17.0",
		"timeExported": "2018-05-03T15:27:30+01:00"
	}
}`
