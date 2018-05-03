package testdata

var RegexMatch = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "regex match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "body": {
						"regexMatch": "<item field=(.*)>"
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
