package testdata

var JsonMatch = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "json match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "body": {
						"jsonMatch": "{\"test\": \"data\"}"
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
