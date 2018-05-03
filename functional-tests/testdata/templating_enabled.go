package testdata

var TemplatingEnabled = `{
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
