package testdata

var JsonPayloadV1 = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "v1 match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "destination": "v1-simulation.com"
                }
            },
            {
                "response": {
                    "status": 200,
                    "body": "v1 match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "requestType": "recording",
                    "destination": "v1-simulation.com",
                    "method": "GET",
                    "scheme": "http",
                    "path": "/path",
                    "query": "",
                    "body": "",
                    "headers": {
                        "Accept-Encoding": [
						    "gzip"
						],
					    "User-Agent": [
						    "Go-http-client/1.1"
						]
                    }
                }
            }
        ],
        "globalActions": {
            "delays": []
        }
    },
    "meta": {
        "schemaVersion": "v1",
        "hoverflyVersion": "v0.10.2",
        "timeExported": "2017-02-23T12:43:48Z"
    }
}`
