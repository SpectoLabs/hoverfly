package testdata

var TemplatingEnabledWithStateInBody = `{
	"data": {
		"pairs": [
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/one"
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
					"body": "setting state",
					"encodedBody": false,
					"templated": true,
					"transitionsState": {
						"eggs": "state for eggs"
					}
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/two"
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
					"body": "{{ State.eggs }}",
					"encodedBody": false,
					"templated": true,
					"transitionsState": {
						"eggs": "present"
					}
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
		"timeExported": "2018-05-03T15:45:52+01:00"
	}
}`

var V3TemplatingEnabledWithStateInBody = `{
	"data": {
		"pairs": [{
				"request": {
					"method": {
						"exactMatch": "GET"
					},
					"path": {
						"exactMatch": "/one"
					}
				},
				"response": {
					"status": 200,
					"body": "setting state",
					"templated": true,
					"transitionsState": {
						"eggs": "state for eggs"
					}
				}
			},
			{
				"request": {
					"method": {
						"exactMatch": "GET"
					},
					"path": {
						"exactMatch": "/two"
					}
				},
				"response": {
					"status": 200,
					"body": "{{ State.eggs }}",
					"templated": true,
					"transitionsState": {
						"eggs": "present"
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
