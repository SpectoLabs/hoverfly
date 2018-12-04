package testdata

var HeaderMatchers = `{
	"data": {
		"pairs": [
			{
				"request": {
					"headers": {
						"test": [
							{
								"matcher": "exact",
								"value": "test"
							}
						]
					}
				},
				"response": {
					"status": 200,
					"body": "header matchers matches",
					"encodedBody": false,
					"templated": false
				}
			},
			{
				"request": {
					"headers": {
						"test2": [
							{
								"matcher": "exact",
								"value": "one;two;three"
							}
						]
					}
				},
				"response": {
					"status": 200,
					"body": "header matchers matches",
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
		"timeExported": "2018-05-03T14:36:59+01:00"
	}
}`

var V4HeaderMatchers = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "header matchers matches",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "headersWithMatchers": {
                        "test": {
                            "exactMatch": "test"
                        }
                    }
                }
            },
            {
                "response": {
                    "status": 200,
                    "body": "header matchers matches",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "headersWithMatchers": {
                        "test2": {
                            "exactMatch": "one;two;three"
                        }
                    }
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
