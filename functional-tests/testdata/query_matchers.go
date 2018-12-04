package testdata

var QueryMatchers = `{
	"data": {
		"pairs": [
			{
				"request": {
					"query": {
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
					"body": "query matchers matches",
					"encodedBody": false,
					"templated": false
				}
			},
			{
				"request": {
					"query": {
						"test": [
							{
								"matcher": "exact",
								"value": "test1;test2"
							}
						]
					}
				},
				"response": {
					"status": 200,
					"body": "query matchers matches",
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
		"timeExported": "2018-05-03T12:47:30+01:00"
	}
}`

var V4QueryMatchers = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "query matchers matches",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "queriesWithMatchers": {
                        "test": {
                            "exactMatch": "test"
                        }
                    }
                }
            },
            {
                "response": {
                    "status": 200,
                    "body": "query matchers matches",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "queriesWithMatchers": {
                        "test": {
                            "exactMatch": "test1;test2"
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
