package functional_tests

var Middleware = `#!/usr/bin/env python
import sys
import json
import logging

logging.basicConfig(filename='middleware_request.log', level=logging.DEBUG)
logging.debug('Middleware "modify_request" called')


def main():
    data = sys.stdin.readlines()
    # this is a json string in one line so we are interested in that one line
    payload = data[0]
    logging.debug(payload)

    payload_dict = json.loads(payload)

    payload_dict['response']['body'] = "CHANGED_RESPONSE_BODY"
    payload_dict['request']['body'] = "CHANGED_REQUEST_BODY"
    payload_dict['response']['status'] = 200
    payload_dict['response']['headers'] = {'Content-Length': ["21"]}

    # returning new payload
    print(json.dumps(payload_dict))

if __name__ == "__main__":
    main()`

var JsonPayload = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "exact match",
                    "encodedBody": false,
                    "headers": {
                        "Header": [
                            "value1",
                            "value2"
                        ]
                    }
                },
                "request": {
                    "path": "/path1",
                    "method": "GET",
                    "destination": "test-server.com",
                    "scheme": "http",
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
            },
            {
                "response": {
                    "status": 200,
                    "body": "template match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "requestType": "template",
                    "destination": "template-server.com"
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

var JsonSimulationGetAndPost = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 201,
                    "body": "body1",
                    "encodedBody": false,
                    "headers": {
                        "Header": [
                            "value1",
                            "value2"
                        ]
                    }
                },
                "request": {
                    "path": "/path1",
                    "method": "GET",
                    "destination": "destination1",
                    "scheme": "http",
                    "query": "",
                    "body": ""
                }
            },
            {
                "response": {
                    "status": 200,
                    "body": "POST body response",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "path": "/path2/resource",
                    "method": "POST",
                    "destination": "another-destination.com",
                    "scheme": "http",
                    "query": "",
                    "body": ""
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
