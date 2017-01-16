.. _simulation_schema:

Simulation schema
=================

.. code:: json


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
