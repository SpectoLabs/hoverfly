package testdata

var TemplatingEnabled = `{
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
					"body": "{{ Request.QueryParam.one }}",
					"encodedBody": false,
					"templated": true
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
		"timeExported": "2018-05-03T15:44:32+01:00"
	}
}`

var V3TemplatingEnabled = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "{{ Request.QueryParam.one }}",
                    "encodedBody": false,
                    "templated" : true
                },
                "request": {
                    "method": {
						"exactMatch": "GET"
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
