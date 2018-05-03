package testdata

var QueryMatchers = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "query matchers matches",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "queriesWithMatchers": {
                        "test": {
                            "exactMatch": "test"
                        }
                    }
                }
            },
            {
                "response": {
                    "status": 200,
                    "body": "query matchers matches",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "queriesWithMatchers": {
                        "test": {
                            "exactMatch": "test1;test2"
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
        "schemaVersion": "v4",
        "hoverflyVersion": "v0.10.2",
        "timeExported": "2017-02-23T12:43:48Z"
    }
}`
