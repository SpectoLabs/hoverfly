package testdata

var StrongestMatchProof = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "first and weakest match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "destination": {
                        "regexMatch": "destination.*"
                    }
                }
            },
            {
                "response": {
                    "status": 200,
                    "body": "second and strongest match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "destination": {
                        "exactMatch": "destination.com"
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
