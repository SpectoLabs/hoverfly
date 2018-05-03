package testdata

var StatePayload = `{
	"data": {
		"pairs": [
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/basket"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "empty",
					"encodedBody": false,
					"templated": false
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/basket"
						}
					],
					"requiresState": {
						"eggs": "present"
					}
				},
				"response": {
					"status": 200,
					"body": "eggs",
					"encodedBody": false,
					"templated": false
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/basket"
						}
					],
					"requiresState": {
						"bacon": "present"
					}
				},
				"response": {
					"status": 200,
					"body": "bacon",
					"encodedBody": false,
					"templated": false
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/basket"
						}
					],
					"requiresState": {
						"bacon": "present",
						"eggs": "present"
					}
				},
				"response": {
					"status": 200,
					"body": "eggs, bacon",
					"encodedBody": false,
					"templated": false
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/add-eggs"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "added eggs",
					"encodedBody": false,
					"templated": false,
					"transitionsState": {
						"eggs": "present"
					}
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/add-bacon"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "added bacon",
					"encodedBody": false,
					"templated": false,
					"transitionsState": {
						"bacon": "present"
					}
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/remove-eggs"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "removed eggs",
					"encodedBody": false,
					"templated": false,
					"removesState": [
						"eggs"
					]
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/remove-bacon"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "removed bacon",
					"encodedBody": false,
					"templated": false,
					"removesState": [
						"bacon"
					]
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
		"timeExported": "2018-05-03T15:40:22+01:00"
	}
}`

var V4StatePayload = `{
	"data": {
		"pairs": [{
				"request": {
					"path": {
						"exactMatch": "/basket"
					}
				},
				"response": {
					"status": 200,
					"body": "empty"
				}
			},
			{
				"request": {
					"path": {
						"exactMatch": "/basket"
					},
					"requiresState": {
						"eggs": "present"
					}
				},
				"response": {
					"status": 200,
					"body": "eggs"
				}
			},
			{
				"request": {
					"path": {
						"exactMatch": "/basket"
					},
					"requiresState": {
						"bacon": "present"
					}
				},
				"response": {
					"status": 200,
					"body": "bacon"
				}
			},
			{
				"request": {
					"path": {
						"exactMatch": "/basket"
					},
					"requiresState": {
						"eggs": "present",
						"bacon": "present"
					}
				},
				"response": {
					"status": 200,
					"body": "eggs, bacon"
				}
			},
			{
				"request": {
					"path": {
						"exactMatch": "/add-eggs"
					}
				},
				"response": {
					"status": 200,
					"body": "added eggs",
					"transitionsState": {
						"eggs": "present"
					}
				}
			},
			{
				"request": {
					"path": {
						"exactMatch": "/add-bacon"
					}
				},
				"response": {
					"status": 200,
					"body": "added bacon",
					"transitionsState": {
						"bacon": "present"
					}
				}
			},
			{
				"request": {
					"path": {
						"exactMatch": "/remove-eggs"
					}
				},
				"response": {
					"status": 200,
					"body": "removed eggs",
					"removesState": ["eggs"]
				}
			},
			{
				"request": {
					"path": {
						"exactMatch": "/remove-bacon"
					}
				},
				"response": {
					"status": 200,
					"body": "removed bacon",
					"removesState": ["bacon"]
				}
			}
		],
		"globalActions": {
			"delays": []
		}
	},
	"meta": {
		"schemaVersion": "v4",
		"hoverflyVersion": "v0.10.2",
		"timeExported": "2017-02-23T12:43:48Z"
	}
}`
