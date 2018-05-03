package testdata

var MultipleMatch = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "multiple matches",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "body": {
						"globMatch": "*<item field=*>*",
                        "regexMatch": "something"
                    }
                }
            },
            {
                "response": {
                    "status": 200,
                    "body": "multiple matches 2",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "destination": {
                        "exactMatch": "destination.com",
                        "globMatch": "*.com"
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
