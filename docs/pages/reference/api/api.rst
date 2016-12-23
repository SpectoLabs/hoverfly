API
===

GET /api/v2/simulation
""""""""""""""""""""""

Gets the entire simulation data. Please remember that the simulation contains all the information Hoverfly can hold; this includes recordings, templates, delays and metadata.

Example response body

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
              "requestType": "recording",
              "path": "/",
              "method": "GET",
              "destination": "myhost.io",
              "scheme": "https",
              "query": "",
              "body": "",
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
              "requestType": "template",
              "path": "/template",
              "method": null,
              "destination": null,
              "scheme": null,
              "query": null,
              "body": null,
              "headers": null
            }
          }
        ],
        "globalActions": {
          "delays": []
        }
      },
      "meta": {
        "schemaVersion": "v1",
        "hoverflyVersion": "v0.9.0",
        "timeExported": "2016-11-11T11:53:52Z"
      }


PUT /api/v2/simulation
""""""""""""""""""""""

This puts the given simulation into Hoverfly, overwriting any existing simulation data.

Example request body

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
              "requestType": "recording",
              "path": "/",
              "method": "GET",
              "destination": "myhost.io",
              "scheme": "https",
              "query": "",
              "body": "",
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
              "requestType": "template",
              "path": "/template",
              "method": null,
              "destination": null,
              "scheme": null,
              "query": null,
              "body": null,
              "headers": null
            }
          }
        ],
        "globalActions": {
          "delays": []
        }
      },
      "meta": {
        "schemaVersion": "v1",
        "hoverflyVersion": "v0.9.0",
        "timeExported": "2016-11-11T11:53:52Z"
      }


-------------------------------------------------------------------------------------------------------------

GET /api/v2/hoverfly
""""""""""""""""""""

Gets configuration information from the running instance of Hoverfly.

Example response body

::

    {
        destination: ".",
        middleware: "",
        mode: "simulate",
        usage: {
            counters: {
                capture: 0,
                modify: 0,
                simulate: 0,
                synthesize: 0
            }
        }
    }


-------------------------------------------------------------------------------------------------------------


GET /api/v2/hoverfly/destination
""""""""""""""""""""""""""""""""

Gets the current destination setting for the running instance of
Hoverfly.

Example response body

::

    {
        destination: "."
    }


PUT /api/v2/hoverfly/destination
""""""""""""""""""""""""""""""""

Sets a new destination for the running instance of Hoverfly, overwriting
the existing destination setting.

Example request body

::

    {
        destination: "new-destination"
    }


-------------------------------------------------------------------------------------------------------------


GET /api/v2/hoverfly/middleware
"""""""""""""""""""""""""""""""

Gets the middleware value for the running instance of Hoverfly. This
will be either an executable command, or an executable command with a
path to a middleware script.

Example response body

::

    {
        "middleware": "python ./middleware.py"
    }


PUT /api/v2/hoverfly/middleware
"""""""""""""""""""""""""""""""

Sets a new middleware value, overwriting the existing middleware value
for the running instance of Hoverfly. The middleware value should be an
executable command, or an executable command with a path to a middleware
script. The command and the file must be available on the Hoverfly host
machine.

Example request body

::

    {
        "middleware": "python ./new-middleware.py"
    }


-------------------------------------------------------------------------------------------------------------


GET /api/v2/hoverfly/mode
"""""""""""""""""""""""""

Gets the mode for the running instance of Hoverfly.

Example response body

::

    {
        mode: "simulate"
    }

--------------

PUT /api/v2/hoverfly/mode
"""""""""""""""""""""""""

Changes the mode of the running instance of Hoverfly.

Example request body

::

    {
        mode: "simulate"
    }


-------------------------------------------------------------------------------------------------------------


GET /api/v2/hoverfly/usage
""""""""""""""""""""""""""

Gets metrics information for the running instance of Hoverfly.

Example response body

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
