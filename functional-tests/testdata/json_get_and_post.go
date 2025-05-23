package testdata

var JsonGetAndPost = `{
	"data": {
		"pairs": [
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
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
							"value": "destination1"
						}
					],
					"scheme": [
						{
							"matcher": "exact",
							"value": "http"
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
					"status": 201,
					"body": "body1",
					"encodedBody": false,
					"headers": {
						"Header": [
							"value1",
							"value2"
						]
					},
					"templated": false
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/path2/resource"
						}
					],
					"method": [
						{
							"matcher": "exact",
							"value": "POST"
						}
					],
					"destination": [
						{
							"matcher": "exact",
							"value": "another-destination.com"
						}
					],
					"scheme": [
						{
							"matcher": "exact",
							"value": "http"
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
					"body": "POST body response",
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
		"timeExported": "2018-05-03T15:09:36+01:00"
	}
}`

