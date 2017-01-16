.. _templates:

Templates
=========

Sometimes simple one-to-one matching of responses to requests is not enough.

Request templates are defined in the :ref:`simulation_schema` by setting the ``"requestType"`` property for a request to ``"template"`` and including only the information in the request that you want Hoverfly to use in the match.

In the example below, Hoverfly will return the same response for any request with the path ``/template``:



.. code:: json

    {
        "data": {
            "pairs": [{
                "response": {
                    "status": 200,
                    "body": "<h1>Matched on template</h1>",
                    "encodedBody": false,
                    "headers": {
                        "Content-Type": ["text/html; charset=utf-8"]
                    }
                },
                "request": {
                    "requestType": "template",
                    "path": "/template"
                }
            }]
        }
    }

For looser matching on URL paths, it is possible to use a wildcard to substitute characters. This is achieved by using the ``*`` symbol. This will match any number of characters, and is case sensitive.

In the next example, Hoverfly will return the same response for requests with the path ``/api/v1/template`` or ``/api/v2/template``.

.. code:: json

    {
        "data": {
            "pairs": [{
                "response": {
                    "status": 200,
                    "body": "<h1>Matched on template</h1>",
                    "encodedBody": false,
                    "headers": {
                        "Content-Type": ["text/html; charset=utf-8"]
                    }
                },
                "request": {
                    "requestType": "template",
                    "path": "/api/*/template"
                }
            }]
        }
    }

The :ref:`simulation_schema` can contain both request recordings and request templates:

.. code:: json

    {
        "data": {
            "pairs": [{
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
            }, {
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
            }],
            "globalActions": {
                "delays": []
            }
        }
    }

A standard workflow might be:

1. Capture some traffic
2. Export it to JSON
3. Edit the JSON to set certain requests to templates, removing the properties for these requests that should be excluded from the match or substituting characters in the URL path for wildcards
4. Re-import the JSON to Hoverfly

.. note::

    If the ``"requestType"`` property is not defined or not recognized, Hoverfly will treat a request as a ``"recording"``.

.. seealso::

    Templating is best understood with a practical example, so please refer to :ref:`addingtemplates` to get hands on experience with templating.
