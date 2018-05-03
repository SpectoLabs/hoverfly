package testdata

var JsonPathMatch = `{
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
						"jsonPathMatch": "$.items[4]"
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
