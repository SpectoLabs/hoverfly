package testdata

var XpathMatch = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "xpath match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "body": {
						"xpathMatch": "//item[count(preceding::item) < 5]"
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
