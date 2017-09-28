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
                    "body": "destination matched",
                    "encodedBody": false
                },
                "request": {
                    "destination": {
                        "exactMatch": "destination-server.com"
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

var JsonPayloadPreloadCache = `{
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
                    }
                }
            },
            {
                "response": {
                    "status": 200,
                    "body": "destination matched",
                    "encodedBody": false
                },
                "request": {
                    "destination": {
                        "exactMatch": "destination-server.com"
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

var JsonPayloadWithDelays = `{
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
                    "body": "destination matched",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "destination": {
                        "exactMatch": "destination-server.com"
                    }
                }
            }
        ],
        "globalActions": {
            "delays": [
                    {
                        "urlPattern": "test-server\\.com",
                        "delay": 100
                    },
                    {
                        "urlPattern": "test-server\\.com",
                        "delay": 110
                    },
                    {
                        "urlPattern": "localhost(.*)",
                        "delay": 110
                    }
            ]
        }
    },
    "meta": {
        "schemaVersion": "v3",
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
        "schemaVersion": "v3",
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
        "schemaVersion": "v3",
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
        "schemaVersion": "v3",
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
        "schemaVersion": "v3",
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
        "schemaVersion": "v3",
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
        "schemaVersion": "v3",
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
        "schemaVersion": "v3",
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
        "schemaVersion": "v3",
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

var ExactMatchPayload = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "exact match 1",
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
                        "Header": [
						    "value1"
						]
                    }
                }
            },
            {
                "response": {
                    "status": 200,
                    "body": "exact match 2",
                    "encodedBody": false,
                    "headers": {
                        "Header": [
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
                        "Header": [
						    "value2"
						]
                    }
                }
            },
            {
                "response": {
                    "status": 200,
                    "body": "destination matched",
                    "encodedBody": false,
                    "headers": {}
                },
                "request": {
                    "destination": {
                        "exactMatch": "destination-server.com"
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

var StrongestMatchProofSimulation = `{
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
                        "exactMatch": "destination.com"
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
                        "exactMatch": "destination.com",
                        "globMatch" : "dest*"
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

var ClosestMissProofSimulation = `{
	"data": {
		"pairs": [{
				"response": {
					"status": 200
				},
				"request": {
					"destination": {
						"exactMatch": "destination.com"
					},
					"body": {
						"exactMatch": "body"
					}
				}
			},
			{
				"response": {
					"status": 200
				},
				"request": {
					"destination": {
						"exactMatch": "destination.com"
					},
					"body": {
						"exactMatch": "body"
					},
					"path": {
						"exactMatch": "/closest-miss"
					}
				}
			},
			{
				"response": {
					"status": 200
				},
				"request": {
					"destination": {
						"exactMatch": "destination.com"
					},
					"body": {
						"exactMatch": "body"
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

var SingleRequestMatcherToResponse = `{
	"data": {
		"pairs": [
			{
				"response": {
					"body": "body"
				},
				"request": {
					"destination": {
						"exactMatch": "miss"
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

var TemplatingEnabled = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "{{ Request.QueryParam.one }}",
                    "encodedBody": false,
                    "templated" : true
                },
                "request": {
                    "method": {
						"exactMatch": "GET"
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

var TemplatingEnabledWithStateInBody = `{
	"data": {
		"pairs": [{
				"request": {
					"method": {
						"exactMatch": "GET"
					},
					"path": {
						"exactMatch": "/one"
					}
				},
				"response": {
					"status": 200,
					"body": "setting state",
					"templated": true,
					"transitionsState": {
						"eggs": "state for eggs"
					}
				}
			},
			{
				"request": {
					"method": {
						"exactMatch": "GET"
					},
					"path": {
						"exactMatch": "/two"
					}
				},
				"response": {
					"status": 200,
					"body": "{{ State.eggs }}",
					"templated": true,
					"transitionsState": {
						"eggs": "present"
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

var TemplatingDisabled = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "{{ Request.QueryParam.singular }}",
                    "encodedBody": false,
                    "templated" : false
                },
                "request": {
                    "method": {
						"exactMatch": "GET"
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

var TemplatingDisabledByDefault = `{
    "data": {
        "pairs": [
            {
                "response": {
                    "status": 200,
                    "body": "{{ Request.QueryParam.one }}",
                    "encodedBody": false
                },
                "request": {
                    "method": {
						"exactMatch": "GET"
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

var Issue607 = `
	{
	"data": {
		"pairs": [
			{
				"response": {
					"status": 200,
					"body": "",
					"encodedBody": false,
					"headers": {
						"Connection": [
							"keep-alive"
						],
						"Content-Type": [
							"application/json"
						],
						"Date": [
							"Tue, 13 Jun 2017 17:54:51 GMT"
						],
						"Hoverfly": [
							"Was-Here"
						],
						"Transfer-Encoding": [
							"chunked"
						]
					}
				},
				"request": {
					"path": {
						"exactMatch": "/billing/v1/servicequotes/123456"
					},
					"method": {
						"exactMatch": "GET"
					},
					"destination": {
						"exactMatch": "domain.com"
					},
					"scheme": {
						"exactMatch": "https"
					},
					"query": {
						"exactMatch": "saleschannel=RETAIL"
					},
					"body": {
						"jsonMatch": ""
					},
					"headers": {
						"Accept": [
							"application/json"
						],
						"Activityid": [
							"ChangeMSISDN_CR_PushtoBill(Get)-200"
						],
						"Applicationid": [
							"ACUI"
						],
						"Authorization": [
							"Bearer token"
						],
						"Cache-Control": [
							"no-cache"
						],
						"Channelid": [
							"RETAIL"
						],
						"Content-Type": [
							"application/json"
						],
						"Interactionid": [
							"123456787"
						],
						"Senderid": [
							"ACUI"
						],
						"User-Agent": [
							"curl/7.54.0"
						],
						"Workflowid": [
							"CHANGEMSISDN"
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
		"schemaVersion": "v3",
		"hoverflyVersion": "v0.12.0",
		"timeExported": "2017-06-13T10:55:12-07:00"
	}
}
`

var StatePayload = `{
		"data": {
			"pairs": [{
					"request": {
						"path": {
							"exactMatch": "/basket"
						}
					},
					"response": {
						"status": 200,
						"body": "empty"
					}
				},
				{
					"request": {
						"path": {
							"exactMatch": "/basket"
						},
						"requiresState": {
							"eggs": "present"
						}
					},
					"response": {
						"status": 200,
						"body": "eggs"
					}
				},
				{
					"request": {
						"path": {
							"exactMatch": "/basket"
						},
						"requiresState": {
							"bacon": "present"
						}
					},
					"response": {
						"status": 200,
						"body": "bacon"
					}
				},
				{
					"request": {
						"path": {
							"exactMatch": "/basket"
						},
						"requiresState": {
							"eggs": "present",
							"bacon": "present"
						}
					},
					"response": {
						"status": 200,
						"body": "eggs, bacon"
					}
				},
				{
					"request": {
						"path": {
							"exactMatch": "/add-eggs"
						}
					},
					"response": {
						"status": 200,
						"body": "added eggs",
						"transitionsState": {
							"eggs": "present"
						}
					}
				},
				{
					"request": {
						"path": {
							"exactMatch": "/add-bacon"
						}
					},
					"response": {
						"status": 200,
						"body": "added bacon",
						"transitionsState": {
							"bacon": "present"
						}
					}
				},
				{
					"request": {
						"path": {
							"exactMatch": "/remove-eggs"
						}
					},
					"response": {
						"status": 200,
						"body": "removed eggs",
						"removesState": ["eggs"]
					}
				},
				{
					"request": {
						"path": {
							"exactMatch": "/remove-bacon"
						}
					},
					"response": {
						"status": 200,
						"body": "removed bacon",
						"removesState": ["bacon"]
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
