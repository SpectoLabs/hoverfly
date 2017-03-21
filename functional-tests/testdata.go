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
                    "path": {
						"exactMatch": "/path1"        
                    },
                    "method": {
						"exactMatch": "GET"
                    },
                    "destination": {
						"exactMatch": "test-server.com"
                    },
                    "scheme": {
						"exactMatch": "http"
                    },
                    "query": {
						"exactMatch": ""
                    },
                    "body": {
						"exactMatch": ""
                    },
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
                    "destination": {
                        "exactMatch": "template-server.com"
                    }
                }
            }
        ],
        "globalActions": {
            "delays": []
        }
    },
    "meta": {
        "schemaVersion": "v2",
        "hoverflyVersion": "v0.10.2",
        "timeExported": "2017-02-23T12:43:48Z"
    }
}`

var XmlSimulation = `{
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
        "schemaVersion": "v2",
        "hoverflyVersion": "v0.10.2",
        "timeExported": "2017-02-23T12:43:48Z"
    }
}`

var XpathSimulation = `{
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
        "schemaVersion": "v2",
        "hoverflyVersion": "v0.10.2",
        "timeExported": "2017-02-23T12:43:48Z"
    }
}`

var JsonPathMatchSimulation = `{
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
        "schemaVersion": "v2",
        "hoverflyVersion": "v0.10.2",
        "timeExported": "2017-02-23T12:43:48Z"
    }
}`

var RegexMatchSimulation = `{
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
        "schemaVersion": "v2",
        "hoverflyVersion": "v0.10.2",
        "timeExported": "2017-02-23T12:43:48Z"
    }
}`

var GlobMatchSimulation = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "glob match",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "body": {
						"globMatch": "*<item field=*>*"
                    }
                }
            }
        ],
        "globalActions": {
            "delays": []
        }
    },
    "meta": {
        "schemaVersion": "v2",
        "hoverflyVersion": "v0.10.2",
        "timeExported": "2017-02-23T12:43:48Z"
    }
}`

var JsonMatchSimulation = `{
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
						"jsonMatch": "{\"test\": \"data\"}"
                    }
                }
            }
        ],
        "globalActions": {
            "delays": []
        }
    },
    "meta": {
        "schemaVersion": "v2",
        "hoverflyVersion": "v0.10.2",
        "timeExported": "2017-02-23T12:43:48Z"
    }
}`

var MultipleMatchSimulation = `{
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
        "schemaVersion": "v2",
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
                    "path": {
						"exactMatch": "/path1"
                    },
                    "method": {
						"exactMatch": "GET"
                    },
                    "destination": {
						"exactMatch": "destination1"
                    },
                    "scheme": {
						"exactMatch": "http"
                    },
                    "query": {
						"exactMatch": ""
                    },
                    "body": {
						"exactMatch": ""
                    }
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
                    "path": {
						"exactMatch": "/path2/resource"
                    },
                    "method": {
						"exactMatch": "POST"
                    },
                    "destination": {
						"exactMatch": "another-destination.com"
                    },
                    "scheme": {
						"exactMatch": "http"
                    },
                    "query": {
						"exactMatch": ""
                    },
                    "body": {
						"exactMatch": ""
                    }
                }
            }
        ],
        "globalActions": {
            "delays": []
        }
    },
    "meta": {
        "schemaVersion": "v2",
        "hoverflyVersion": "v0.10.2",
        "timeExported": "2017-02-23T12:43:48Z"
    }
}`

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

var ExactMatchPayload = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "exact match",
                    "encodedBody": false,
                    "headers": {
                        "Header": [
                            "value1"
                        ]
                    }
                },
                "request": {
                    "path": {
						"exactMatch": "/path1"        
                    },
                    "method": {
						"exactMatch": "GET"
                    },
                    "destination": {
						"exactMatch": "test-server.com"
                    },
                    "scheme": {
						"exactMatch": "http"
                    },
                    "query": {
						"exactMatch": ""
                    },
                    "body": {
						"exactMatch": ""
                    },
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
                    "destination": {
                        "exactMatch": "template-server.com"
                    }
                }
            }
        ],
        "globalActions": {
            "delays": []
        }
    },
    "meta": {
        "schemaVersion": "v2",
        "hoverflyVersion": "v0.10.2",
        "timeExported": "2017-02-23T12:43:48Z"
    }
}`
