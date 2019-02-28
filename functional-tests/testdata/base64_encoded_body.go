package testdata

var Base64EncodedBody = `{
	"data": {
		"pairs": [
			{
				"request": {
					"destination": [
						{
							"matcher": "exact",
							"value": "test-server.com"
						}
					],
					"path": [
						{
							"matcher": "exact",
							"value": "/image.png"
						}
					],
					"method": [
						{
							"matcher": "exact",
							"value": "GET"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAAAAAA6fptVAAAACklEQVR4nGP6DwABBQECz6AuzQAAAABJRU5ErkJggg==",
					"encodedBody": true,
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
		"hoverflyVersion": "v1.0.0-rc.2",
		"timeExported": "2019-02-28T14:48:28Z"
	}
}`