.. _rest_api:


REST API
========

GET /api/v2/simulation
""""""""""""""""""""""

Gets all simulation data. The simulation JSON contains all the information Hoverfly can hold; this includes recordings, templates, delays and metadata.

**Example response body**
::

    {
      "data": {
        "pairs": [
          {
            "response": {
              "status": 200,
              "body": "<h1>Matched on recording</h1>",
              "encodedBody": false,
              "headers": {
                "Content-Type": [
                  "text/html; charset=utf-8"
                ]
              }
            },
            "request": {
              "path": {
	        "exactMatch": "/"
	      },
              "method": {
	        "exactMatch": "GET"
              },
	      "destination": {
	        "exactMatch": "myhost.io"
              },
	      "scheme": {
	        "exactMatch": "https"
              },
	      "query": {
	        "exactMatch": ""
	      },
              "body": {
	        "exactMatch": ""
	      },
              "headers": {
                "Accept": [
                  "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
                ],
                "Content-Type": [
                  "text/plain; charset=utf-8"
                ],
                "User-Agent": [
                  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36"
                ]
              }
            }
          },
          {
            "response": {
              "status": 200,
              "body": "<h1>Matched on template</h1>",
              "encodedBody": false,
              "headers": {
                "Content-Type": [
                  "text/html; charset=utf-8"
                ]
              },
              "templated": false
            },
            "request": {
              "path": {
	        "exactMatch": "/template"
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
        "hoverflyVersion": "v0.11.0",
        "timeExported": "2016-11-11T11:53:52Z"
      }
    }


PUT /api/v2/simulation
""""""""""""""""""""""

This puts the supplied simulation JSON into Hoverfly, overwriting any existing simulation data.

**Example request body**
::

    {
      "data": {
        "pairs": [
          {
            "response": {
              "status": 200,
              "body": "<h1>Matched on recording</h1>",
              "encodedBody": false,
              "headers": {
                "Content-Type": [
                  "text/html; charset=utf-8"
                ]
              }
            },
            "request": {
              "path": {
	        "exactMatch": "/"
	      },
              "method": {
	        "exactMatch": "GET"
              },
	      "destination": {
	        "exactMatch": "myhost.io"
              },
	      "scheme": {
	        "exactMatch": "https"
              },
	      "query": {
	        "exactMatch": ""
	      },
              "body": {
	        "exactMatch": ""
	      },
              "headers": {
                "Accept": [
                  "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
                ],
                "Content-Type": [
                  "text/plain; charset=utf-8"
                ],
                "User-Agent": [
                  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36"
                ]
              }
            }
          },
          {
            "response": {
              "status": 200,
              "body": "<h1>Matched on template</h1>",
              "encodedBody": false,
              "headers": {
                "Content-Type": [
                  "text/html; charset=utf-8"
                ]
              }
            },
            "request": {
              "path": {
	        "exactMatch": "/template"
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
        "hoverflyVersion": "v0.11.0",
        "timeExported": "2016-11-11T11:53:52Z"
      }
    }

-------------------------------------------------------------------------------------------------------------

GET /api/v2/simulation/schema
"""""""""""""""""""""""""""""
Gets the JSON Schema used to validate the simulation JSON.


-------------------------------------------------------------------------------------------------------------

GET /api/v2/hoverfly
""""""""""""""""""""

Gets configuration information from the running instance of Hoverfly.

**Example response body**
::

    {
        "destination": ".",
        "middleware": {
        "binary": "python",
		"script": "# a python script would go here",
		"remote": ""
	},
        "mode": "simulate",
        "usage": {
            "counters": {
                "capture": 0,
                "modify": 0,
                "simulate": 0,
                "synthesize": 0
            }
        }
    }


-------------------------------------------------------------------------------------------------------------


GET /api/v2/hoverfly/destination
""""""""""""""""""""""""""""""""

Gets the current destination setting for the running instance of
Hoverfly.

**Example response body**
::

    {
        destination: "."
    }


PUT /api/v2/hoverfly/destination
""""""""""""""""""""""""""""""""

Sets a new destination for the running instance of Hoverfly, overwriting
the existing destination setting.

**Example request body**
::

    {
        destination: "new-destination"
    }


-------------------------------------------------------------------------------------------------------------


GET /api/v2/hoverfly/middleware
"""""""""""""""""""""""""""""""

Gets the middleware settings for the running instance of Hoverfly. This
could be either an executable binary, a script that can be executed with 
a binary or a URL to remote middleware.

**Example response body**
::

    {
        "binary": "python",
	      "script": "#python code goes here",
	      "remote": ""
    }


PUT /api/v2/hoverfly/middleware
"""""""""""""""""""""""""""""""

Sets new middleware, overwriting the existing middleware
for the running instance of Hoverfly. The middleware being set
can be either an executable binary located on the host, a script
and the binary to execute it or the URL to a remote middleware.

**Example request body**
::

    {
        "binary": "python",
	      "script": "#python code goes here",
	      "remote": ""
    }


-------------------------------------------------------------------------------------------------------------


GET /api/v2/hoverfly/mode
"""""""""""""""""""""""""

Gets the mode for the running instance of Hoverfly.

**Example response body**
::

    {
        "mode": "capture",
        "arguments": {}
    }

--------------

PUT /api/v2/hoverfly/mode
"""""""""""""""""""""""""

Changes the mode of the running instance of Hoverfly.

**Example request body**
::

    {
        "mode": "capture"
    }


-------------------------------------------------------------------------------------------------------------


GET /api/v2/hoverfly/usage
""""""""""""""""""""""""""

Gets metrics information for the running instance of Hoverfly.

**Example response body**
::

    {
        "metrics": {
            "counters": {
                "capture": 0,
                "modify": 0,
                "simulate": 0,
                "synthesize": 0
            }
        }
    }


-------------------------------------------------------------------------------------------------------------


GET /api/v2/hoverfly/version
""""""""""""""""""""""""""""

Gets the version of Hoverfly.

**Example response body**
::

    {
        "version": "v0.10.1"
    }


-------------------------------------------------------------------------------------------------------------


GET /api/v2/hoverfly/upstream-proxy
"""""""""""""""""""""""""""""""""""

Gets the upstream proxy configured for Hoverfly.

**Example response body**
::

    {
        "upstreamProxy": "proxy.corp.big-it-company.org:8080"
    }


-------------------------------------------------------------------------------------------------------------


GET /api/v2/hoverfly/pac
""""""""""""""""""""""""

Gets the PAC file configured for Hoverfly.


-------------------------------------------------------------------------------------------------------------


PUT /api/v2/hoverfly/pac
""""""""""""""""""""""""

Sets the PAC file for Hoverfly.


-------------------------------------------------------------------------------------------------------------


DELETE /api/v2/hoverfly/pac
""""""""""""""""""""""""

Unsets the PAC file configured for Hoverfly.

-------------------------------------------------------------------------------------------------------------


GET /api/v2/cache
""""""""""""""""""""
Gets the requests and responses stored in the cache.

**Example response body**
::

    {
        "cache": [
            {
                "key": "2fc8afceec1b6bcf99ff1f547c1f5b11",
                "matchingPair": {
                    "request": {
                        "path": {
                            "exactMatch": "hoverfly.io"
                        }
                    },
                    "response": {
                        "status": 200,
                        "body": "response body",
                        "encodedBody": false,
                        "headers": {
                            "Hoverfly": [
                                "Was-Here"
                            ]
                        }
                    }
                },
                "headerMatch": false
            }
        ]
    }

-------------------------------------------------------------------------------------------------------------


DELETE /api/v2/cache
""""""""""""""""""""
Delete all requests and responses stored in the cache.


-------------------------------------------------------------------------------------------------------------


GET /api/v2/logs
""""""""""""""""""""
Gets the logs from Hoverfly.

**Example response body**
::

    {
        "logs": [
            {
                "level": "info",
                "msg": "serving proxy",
                "time": "2017-03-13T12:22:39Z"
            },  
            {
                "destination": ".",
                "level": "info",
                "mode": "simulate",
                "msg": "current proxy configuration",
                "port": "8500",
                "time": "2017-03-13T12:22:39Z"
            },  
            {
                "destination": ".",
                "Mode": "simulate",
                "ProxyPort": "8500",
                "level": "info",
                "msg": "Proxy prepared...",
                "time": "2017-03-13T12:22:39Z"
            },  
        ]
    }


-------------------------------------------------------------------------------------------------------------


GET /api/v2/journal
"""""""""""""""""""
Gets the journal from Hoverfly. Each journal entry contains both the request Hoverfly recieved and the response 
it served along with the mode Hoverfly was in, the time the request was recieved and the time taken for Hoverfly
to process the request. Latency is in milliseconds.

**Example response body**
::
  {
    "journal": [
      {
        "request": {
          "path": "/",
          "method": "GET",
          "destination": "hoverfly.io",
          "scheme": "http",
          "query": "",
          "body": "",
          "headers": {
            "Accept": [
              "*/*"
            ],
            "Proxy-Connection": [
              "Keep-Alive"
            ],
            "User-Agent": [
              "curl/7.50.2"
            ]
          }
        },
        "response": {
          "status": 502,
          "body": "Hoverfly Error!\n\nThere was an error when matching\n\nGot error: Could not find a match for request, create or record a valid matcher first!",
          "encodedBody": false,
          "headers": {
            "Content-Type": [
              "text/plain"
            ]
          }
        },
        "mode": "simulate",
        "timeStarted": "2017-07-17T10:41:59.168+01:00",
        "latency": 0.61334
      }
    ]
  }


-------------------------------------------------------------------------------------------------------------


DELETE /api/v2/journal
""""""""""""""""""""""
Delete all entries stored in the journal.


-------------------------------------------------------------------------------------------------------------


POST /api/v2/journal
""""""""""""""""""""
Filter and search entries stored in the journal.

**Example request body**
::
    {
        "request": {
            "destination": {
              "exactMatch": "hoverfly.io"
            }
        }
    }


-------------------------------------------------------------------------------------------------------------


GET /api/v2/state
"""""""""""""""""
Gets the state from Hoverfly. State is represented as a set of key value pairs.

**Example response body**
::
  {
    "state": {
      "page_state": "CHECKOUT"
    }
  }


-------------------------------------------------------------------------------------------------------------

DELETE /api/v2/state
""""""""""""""""""""
Deletes all state from Hoverfly.

-------------------------------------------------------------------------------------------------------------

PUT /api/v2/state
"""""""""""""""""""
Deletes all state from Hoverfly and then sets the state to match the state in the request body.

**Example request body**
::
  {
    "state": {
      "page_state": "CHECKOUT"
    }
  }

-------------------------------------------------------------------------------------------------------------

PATCH /api/v2/state
"""""""""""""""""""
Updates state in Hoverfly. Will update each state key referenced in the request body. 

**Example request body**
::
  {
    "state": {
      "page_state": "CHECKOUT"
    }
  }

-------------------------------------------------------------------------------------------------------------


GET /api/v2/diff
"""""""""""""""""
Gets all reports containing response differences from Hoverfly. The diffs are represented as lists of strings grouped by the same requests.

**Example response body**
::
  {
    "diff": [{
      "request": {
        "method": "GET",
        "host": "time.jsontest.com",
        "path": "/",
        "query": ""
      },
      "diffReports": [{
        "timestamp": "2018-03-16T17:45:40Z",
        "diffEntries": [{
          "field": "header/X-Cloud-Trace-Context",
          "expected": "[ec6c455330b682c3038ba365ade6652a]",
          "actual": "[043c9bb2eafa1974bc09af654ef15dc3]"
        }, {
          "field": "header/Date",
          "expected": "[Fri, 16 Mar 2018 17:45:34 GMT]",
          "actual": "[Fri, 16 Mar 2018 17:45:41 GMT]"
        }, {
          "field": "body/time",
          "expected": "05:45:34 PM",
          "actual": "05:45:41 PM"
        }, {
          "field": "body/milliseconds_since_epoch",
          "expected": "1.521222334104e+12",
          "actual": "1.521222341017e+12"
        }]
      }]
    }]
  }

-------------------------------------------------------------------------------------------------------------

DELETE /api/v2/diff
""""""""""""""""""""
Deletes all reports containing differences from Hoverfly.
