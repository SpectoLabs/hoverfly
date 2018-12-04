package testdata

var ClosestMissProof = `{
	"data": {
		"pairs": [
			{
				"request": {
					"destination": [
						{
							"matcher": "exact",
							"value": "destination.com"
						}
					],
					"body": [
						{
							"matcher": "exact",
							"value": "body"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "",
					"encodedBody": false,
					"templated": false
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/closest-miss"
						}
					],
					"destination": [
						{
							"matcher": "exact",
							"value": "destination.com"
						}
					],
					"body": [
						{
							"matcher": "exact",
							"value": "body"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "",
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
		"timeExported": "2018-05-03T12:08:35+01:00"
	}
}`

var V3ClosestMissProof = `{
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
