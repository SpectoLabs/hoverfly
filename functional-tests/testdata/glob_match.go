package testdata

var GlobMatch = `{
	"data": {
		"pairs": [
			{
				"request": {
					"body": [
						{
							"matcher": "glob",
							"value": "*<item field=*>*"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "glob match",
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
		"timeExported": "2018-05-03T12:15:49+01:00"
	}
}
`

var V3GlobMatch = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "glob match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "body": {
						"globMatch": "*<item field=*>*"
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
