.. _captured_traffic:

Captured traffic
================

Hoverfly's core functionality is to capture requests and responses ("traffic") to create API simulations.

Request response pairs
......................

When you capture traffic using Hoverfly's :ref:`capture_mode` and export resulting the simulation to JSON, you will see *request response pairs*:

.. code:: json

    {
      "response": {
        "status": 200,
        "body": "body here",
        "encodedBody": false,
        "headers": {
          "Content-Type": ["text/html; charset=utf-8"]
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
          "Accept": ["text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"],
          "User-Agent": ["Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36"]
        }
      }
    }

Notice the ``"response"`` and ``"request"`` key values in the document above. They both contain all of the fields from the original request and response. The headers are also stored.

Please also notice the ``"requestType"`` key, with its corresponding ``"recording"`` value: this denotes that the request response pair was created while using Hoverfly in :ref:`capture_mode`, and that in :ref:`simulate_mode`, Hoverfly will return this exact response when it receives this request.


.. note::

    Since JSON does not support binary data, binary responses are base64 encoded. This is denoted by the encodedBody field. Hoverfly automatically encodes and decodes the data during the export and import phases.


Matching
........

Hoverfly simulates APIs by `matching` responses to incoming requests.

Imagine scanning through a dictionary for a word, and then looking up its definition. Hoverfly does exactly that, but the "word" is the URL that was "recorded" in :ref:`capture_mode`, plus all the fields in the ``"request"`` part of the document above, with the exception of the headers.

.. note::

    The reason headers are not included in the match is because they can vary depending on the client.

This one-to-one matching strategy is extremely fast, but in some cases you may want Hoverfly to return a single response for more than one request. This is possible using :ref:`templates`.
