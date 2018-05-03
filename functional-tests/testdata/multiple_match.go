package testdata

var MultipleMatch = `{
	"data": {
		"pairs": [
			{
				"request": {
					"body": [
						{
							"matcher": "glob",
							"value": "*<item field=*>*"
						},
						{
							"matcher": "regex",
							"value": "something"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "multiple matches",
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
						},
						{
							"matcher": "glob",
							"value": "*.com"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "multiple matches 2",
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
		"timeExported": "2018-05-03T15:27:30+01:00"
	}
}`

var V3MultipleMatch = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "multiple matches",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "body": {
						"globMatch": "*<item field=*>*",
                        "regexMatch": "something"
                    }
                }
            },
            {
                "response": {
                    "status": 200,
                    "body": "multiple matches 2",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "destination": {
                        "exactMatch": "destination.com",
                        "globMatch": "*.com"
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
