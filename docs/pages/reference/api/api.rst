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
            "request": {
              "path": [
                {
                  "matcher": "exact",
                  "value": "/"
                }
              ],
              "method": [
                {
                  "matcher": "exact",
                  "value": "GET"
                }
              ],
              "destination": [
                {
                  "matcher": "exact",
                  "value": "myhost.io"
                }
              ],
              "scheme": [
                {
                  "matcher": "exact",
                  "value": "https"
                }
              ],
              "body": [
                {
                  "matcher": "exact",
                  "value": ""
                }
              ],
              "headers": {
                "Accept": [
                  {
                    "matcher": "glob",
                    "value": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
                  }
                ],
                "Content-Type": [
                  {
                    "matcher": "glob",
                    "value": "text/plain; charset=utf-8"
                  }
                ],
                "User-Agent": [
                  {
                    "matcher": "glob",
                    "value": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36"
                  }
                ]
              },
              "query": {
                "status": [
                  {
                    "matcher": "exact",
                    "value": "available"
                  }
                ]
              }
            },
            "response": {
              "status": 200,
              "body": "<h1>Matched on recording</h1>",
              "encodedBody": false,
              "headers": {
                "Content-Type": [
                  "text/html; charset=utf-8"
                ]
              },
              "templated": false
            }
          },
          {
            "request": {
              "path": [
                {
                  "matcher": "exact",
                  "value": "/template"
                }
              ]
            },
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
            }
          }
        ],
        "globalActions": {
          "delays": [],
          "delaysLogNormal": []
        }
      },
      "meta": {
        "schemaVersion": "v5",
        "hoverflyVersion": "v1.0.0",
        "timeExported": "2019-05-30T22:14:24+01:00"
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
            "request": {
              "path": [
                {
                  "matcher": "exact",
                  "value": "/"
                }
              ],
              "method": [
                {
                  "matcher": "exact",
                  "value": "GET"
                }
              ],
              "destination": [
                {
                  "matcher": "exact",
                  "value": "myhost.io"
                }
              ],
              "scheme": [
                {
                  "matcher": "exact",
                  "value": "https"
                }
              ],
              "body": [
                {
                  "matcher": "exact",
                  "value": ""
                }
              ],
              "headers": {
                "Accept": [
                  {
                    "matcher": "glob",
                    "value": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
                  }
                ],
                "Content-Type": [
                  {
                    "matcher": "glob",
                    "value": "text/plain; charset=utf-8"
                  }
                ],
                "User-Agent": [
                  {
                    "matcher": "glob",
                    "value": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36"
                  }
                ]
              },
              "query": {
                "status": [
                  {
                    "matcher": "exact",
                    "value": "available"
                  }
                ]
              }
            },
            "response": {
              "status": 200,
              "body": "<h1>Matched on recording</h1>",
              "encodedBody": false,
              "headers": {
                "Content-Type": [
                  "text/html; charset=utf-8"
                ]
              },
              "templated": false
            }
          },
          {
            "request": {
              "path": [
                {
                  "matcher": "exact",
                  "value": "/template"
                }
              ]
            },
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
            }
          }
        ],
        "globalActions": {
          "delays": [],
          "delaysLogNormal": []
        }
      },
      "meta": {
        "schemaVersion": "v5",
        "hoverflyVersion": "v1.0.0",
        "timeExported": "2019-05-30T22:14:24+01:00"
      }
    }

POST /api/v2/simulation
"""""""""""""""""""""""

This appends the supplied simulation JSON to the existing simulation data in Hoverfly. Any pair that has request data identical to the existing ones will not be added.

**Example request body**
::

    {
      "data": {
        "pairs": [
          {
            "request": {
              "path": [
                {
                  "matcher": "exact",
                  "value": "/"
                }
              ],
              "method": [
                {
                  "matcher": "exact",
                  "value": "GET"
                }
              ],
              "destination": [
                {
                  "matcher": "exact",
                  "value": "myhost.io"
                }
              ],
              "scheme": [
                {
                  "matcher": "exact",
                  "value": "https"
                }
              ],
              "body": [
                {
                  "matcher": "exact",
                  "value": ""
                }
              ],
              "headers": {
                "Accept": [
                  {
                    "matcher": "glob",
                    "value": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
                  }
                ],
                "Content-Type": [
                  {
                    "matcher": "glob",
                    "value": "text/plain; charset=utf-8"
                  }
                ],
                "User-Agent": [
                  {
                    "matcher": "glob",
                    "value": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36"
                  }
                ]
              },
              "query": {
                "status": [
                  {
                    "matcher": "exact",
                    "value": "available"
                  }
                ]
              }
            },
            "response": {
              "status": 200,
              "body": "<h1>Matched on recording</h1>",
              "encodedBody": false,
              "headers": {
                "Content-Type": [
                  "text/html; charset=utf-8"
                ]
              },
              "templated": false
            }
          },
          {
            "request": {
              "path": [
                {
                  "matcher": "exact",
                  "value": "/template"
                }
              ]
            },
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
            }
          }
        ],
        "globalActions": {
          "delays": [],
          "delaysLogNormal": []
        }
      },
      "meta": {
        "schemaVersion": "v5",
        "hoverflyVersion": "v1.0.0",
        "timeExported": "2019-05-30T22:14:24+01:00"
      }
    }


DELETE /api/v2/simulation
"""""""""""""""""""""""""

Unsets the simulation data for Hoverfly.

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
        "cors": {
            "enabled": true,
            "allowOrigin": "*",
            "allowMethods": "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS",
            "allowHeaders": "Content-Type,Origin,Accept,Authorization,Content-Length,X-Requested-With",
            "preflightMaxAge": 1800,
            "allowCredentials": true
        },
        "destination": ".",
        "middleware": {
            "binary": "python",
            "script": "# a python script would go here",
            "remote": ""
        },
        "mode": "simulate",
        "arguments": {
            "matchingStrategy": "strongest"
        },
        "isWebServer": false,
        "usage": {
            "counters": {
                "capture": 0,
                "modify": 0,
                "simulate": 0,
                "synthesize": 0
            }
        },
        "version": "v1.3.3",
        "upstreamProxy": ""
    }

-------------------------------------------------------------------------------------------------------------

GET /api/v2/hoverfly/cors
"""""""""""""""""""""""""

Gets CORS configuration information from the running instance of Hoverfly.

**Example response body**
::

    {
        "enabled": true,
        "allowOrigin": "*",
        "allowMethods": "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS",
        "allowHeaders": "Content-Type, Origin, Accept, Authorization, Content-Length, X-Requested-With",
        "preflightMaxAge": 1800,
        "allowCredentials": true
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
        "arguments": {
            "headersWhitelist": [
                "*"
            ],
            "stateful": true,
            "overwriteDuplicate": true
        }
    }

--------------

PUT /api/v2/hoverfly/mode
"""""""""""""""""""""""""

Changes the mode of the running instance of Hoverfly. Pass additional arguments to set the mode options.

**Example request body**
::

    {
        "mode": "capture",
        "arguments": {
            "headersWhitelist": [
                "*"
            ],
            "stateful": true,
            "overwriteDuplicate": true
        }
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

Gets the PAC file configured for Hoverfly. The response contains plain text with PAC file.
If no PAC was provided before, the response is 404 with contents:

::

    {
        "error": "Not found"
    }



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
                        "path": [{
                            "matcher": "exact",
                            "value": "hoverfly.io"
                        }]
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
                "headerMatch": false,
                "closestMiss": null
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

It supports multiple parameters to limit the amount of entries returned:

- ``limit`` - Maximum amount of entries. 500 by default;
- ``from`` - Timestamp to start filtering from.

Running hoverfly with ``-logs-size=0`` disables logging and 500 response is returned with body:

::

    {
        "error": "Logs disabled"
    }


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
Gets the journal from Hoverfly. Each journal entry contains both the request Hoverfly received and the response
it served along with the mode Hoverfly was in, the time the request was received and the time taken for Hoverfly
to process the request. Latency is in milliseconds.

It supports paging using the ``offset`` and ``limit`` query parameters.

It supports multiple parameters to limit the amount of entries returned:

- ``limit`` - Maximum amount of entries. 500 by default;
- ``offset`` - Offset of the first element;
- ``to`` - Timestamp to start filtering to;
- ``from`` - Timestamp to start filtering from;
- ``sort`` - Sort results in format "field:order". Supported fields: ``timestarted`` and ``latency``. Supported orders: ``asc`` and ``desc``.

Running hoverfly with ``-journal-size=0`` disables logging and 500 response is returned with body:

::

    {
        "error": "Journal disabled"
    }


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
    ],
    "offset": 0,
    "limit": 25,
    "total": 1
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
            "destination": [{
              "matcher": "exact",
              "value": "hoverfly.io"
            }]
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


POST /api/v2/diff
"""""""""""""""""
Gets reports containing response differences from Hoverfly filtered on basis of excluded criteria provided(i.e. headers and response keys in jsonpath format to exclude). The diffs are in same format as we receive in GET request.

**Example request body**
::
  {
    "excludedHeaders":["Date"],
    "excludedResponseFields":["$.time"]
  }


-------------------------------------------------------------------------------------------------------------

DELETE /api/v2/diff
""""""""""""""""""""
Deletes all reports containing differences from Hoverfly.

DELETE /api/v2/shutdown
""""""""""""""""""""
Shuts down the hoverfly instance.
