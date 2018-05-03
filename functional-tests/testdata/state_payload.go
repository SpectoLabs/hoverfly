package testdata

var StatePayload = `{
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
