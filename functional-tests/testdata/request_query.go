package testdata

var EmptyQuery = `{
	"data": {
		"pairs": [
			{
				"request": {
					"method": [
						{
							"matcher": "exact",
							"value": "GET"
						}
					],
					"query": {}
				},
				"response": {
					"status": 200,
					"body": "hello"
				}
			}
		],
		"globalActions": {
			"delays": []
		}
	},
	"meta": {
		"schemaVersion": "v5",
		"hoverflyVersion": "v0.17.3"
	}
}`

var NoQuery = `{
	"data": {
		"pairs": [
			{
				"request": {
					"method": [
						{
							"matcher": "exact",
							"value": "GET"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "hello"
				}
			}
		],
		"globalActions": {
			"delays": []
		}
	},
	"meta": {
		"schemaVersion": "v5",
		"hoverflyVersion": "v0.17.3"
	}
}`
