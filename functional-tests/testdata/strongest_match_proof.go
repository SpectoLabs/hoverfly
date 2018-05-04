package testdata

var StrongestMatchProof = `{
	"data": {
		"pairs": [
			{
				"request": {
					"destination": [
						{
							"matcher": "regex",
							"value": "destination.*"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "first and weakest match",
					"encodedBody": false,
					"templated": false
				}
			},
			{
				"request": {
					"destination": [
						{
							"matcher": "exact",
							"value": "destination.com"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "second and strongest match",
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
		"timeExported": "2018-05-03T15:51:35+01:00"
	}
}`

var V3StrongestMatchProof = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "first and weakest match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "destination": {
                        "regexMatch": "destination.*"
                    }
                }
            },
            {
                "response": {
                    "status": 200,
                    "body": "second and strongest match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "destination": {
                        "exactMatch": "destination.com"
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
