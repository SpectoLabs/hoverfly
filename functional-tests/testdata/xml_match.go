package testdata

var XmlMatch = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "xml match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "body": {
						"xmlMatch": "<items><item>one</item></items>"
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
