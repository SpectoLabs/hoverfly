package testdata

var TemplatingDisabled = `{
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
					"body": "{{ Request.QueryParam.singular }}",
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
		"timeExported": "2018-05-03T15:47:29+01:00"
	}
}`

var V3TemplatingDisabled = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "{{ Request.QueryParam.singular }}",
                    "encodedBody": false,
                    "templated" : false
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
