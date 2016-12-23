Recording
~~~~~~~~~

Hoverfly's core functionality is to capture, and later simulate HTTP requests and responses.

Request Response Pairs
......................

When you capture traffic in Hoverfly, and export the simulation, notice it gets stored as *request response pairs*:

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

Please notice the ``"response"`` and ``"request"`` key values in the document above. Notice they both contain all of the fields from the original request and response that was made over HTTP. The headers are also stored, with invidual key / value pairs.

Please also notice the ``"requestType"`` key, with its corresponding ``"recording"`` value: they denote that this request response pair is a recording.

.. note::

    Since JSON does not support binary data, binary responses are base64 encoded. This is denoted by the encodedBody field. Hoverfly automatically encodes and decodes the data during the export and import phases.


Matching
........

Remember that Hoverfly's core functionality is capturing HTTP interactions, and then playing them back. To do this, it needs to `match` requests, with your current request, to give you the appropriate response.

Imagine scanning through a dictionary for a word, and then looking up its definition. Well, Hoverfly needs to do exactly that, but the word in question is the URL you supplied during the capture phase, as well as all the fields in the ``"request"`` part of the document above, with the exception of the headers.

This is what is referred to as matching.

.. note::

    The reason headers are not included during matching is because they can vary based on the client.