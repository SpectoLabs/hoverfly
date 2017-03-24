.. _pairs:

Request Responses Pairs
=======================

.. todo:: @tjcunliffe to review

Hoverfly simulates APIs by `matching` incoming requests to requests that it has captured previously, and returning a response that is associated with the matched request.

Imagine scanning through a dictionary for a word, and then looking up its definition. Hoverfly does exactly that, but the “word” is the HTTP request that was “captured” 
in Capture mode, and the “definition” is the response.

Request Matcher
---------------

Hoverfly matches incoming requests to captured requests by comparing the following fields:

+--------------------+---------------------+-------------------------------------+
| HTTP Request Field | Value type          | Example                             |
+====================+=====================+=====================================+
| scheme             | string              | "https"                             |
+--------------------+---------------------+-------------------------------------+
| method             | string              | "GET"                               |
+--------------------+---------------------+-------------------------------------+
| destination        | string              | "docs.hoverfly.io"                  |
+--------------------+---------------------+-------------------------------------+
| path               | string              | "/pages/keyconcepts/templates.html" |
+--------------------+---------------------+-------------------------------------+
| query              | string              | "query=true"                        |
+--------------------+---------------------+-------------------------------------+
| body               | string              | ""                                  |
+--------------------+---------------------+-------------------------------------+
| headers            | map[string][]string | "User-Agent: ["http-client"]"       |
+--------------------+---------------------+-------------------------------------+

When Hoverfly captures a request, it creates a Request Matcher for each field in the request. A Request Matcher consists of:
 
 - the request field name 
 - the request field value 
 - the type of match that will be used to compare the captured request field value to the incoming request field value 

By default, Hoverfly will set the type of match to "exactMatch" for each field. Below is a Request Matcher set from an example Hoverfly simulation JSON file.

.. literalinclude:: ../../simulations/basic-simulation.json
   :lines: 4-24
   :linenos:
   :language: javascript


:ref:`View entire simulation file <basic_simulation>`

The matching strategy that Hoverfly uses to compare an incoming request to a captured request can be changed by editing the Request Matchers in the simulation
JSON file. 

It is not necessary to have a Request Matcher for every request field. By omitting Request Matchers, it is possible to implement **partial matching** - meaning
that more than Hoverfly will return one response for more than one incoming request. 

For example, this Request Matcher will match any request to "docs.hoverfly.io":

.. code:: json

    {
        "destination": {
            "exactMatch": "docs.hoverfly.io"
        },
    }

Although the default matcher type is "exactMatch", there are many other matchers to choose from.

Responses
---------

Each Request Matcher set has a response associated with is. If the request match is successful, Hoverfly will return the response to the client.

.. literalinclude:: ../../simulations/basic-simulation.json
   :lines: 25-32
   :linenos:
   :language: javascript

:ref:`View entire simulation file <basic_simulation>`

Editing the fields in response, combined with editing the Request Matcher set, makes it possible to configure complex request/response logic. 

Since JSON does not support binary data, binary responses are base64 encoded. This is denoted by the encodedBody field. 
Hoverfly automatically encodes and decodes the data during the export and import phases.

.. literalinclude:: ../../simulations/basic-encoded-simulation.json
   :lines: 27-28
   :linenos:
   :language: javascript

:ref:`View entire simulation file <basic_encoded_simulation>`