package testdata

var ArrayMatcher = `{
	"data": {
		"pairs": [
			{
				"request": {
					"headers": {
						"test1": [
							{
								"matcher": "array",
								"value": ["a", "b", "c"]
							}
						]
					}
				},
				"response": {
					"status": 200,
					"body": "array matchers matches",
					"encodedBody": false,
					"templated": false
				}
			},
			{
				"request": {
					"query": {
						"test": [
							{
								"matcher": "array",
								"value": ["value1", "value2", "value3"],
								"config": {
									"ignoreOrder": true
								}
							}
						]
					}
				},
				"response": {
					"status": 200,
					"body": "array matchers matches query",
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
		"schemaVersion": "v5.2",
		"hoverflyVersion": "v0.17.0",
		"timeExported": "2018-05-03T14:36:59+01:00"
	}
}`
