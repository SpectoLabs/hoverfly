package testdata

var HeaderMatchers = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "header matchers matches",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "headersWithMatchers": {
                        "test": {
                            "exactMatch": "test"
                        }
                    }
                }
            },
            {
                "response": {
                    "status": 200,
                    "body": "header matchers matches",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "headersWithMatchers": {
                        "test2": {
                            "exactMatch": "one;two;three"
                        }
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
