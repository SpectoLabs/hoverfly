.. request_matchers_and_responses:

Request Matchers and Responses
==============================

.. todo:: @tjcunliffe to review

Hoverfly simulates APIs by `matching` incoming requests to requests that it has captured previously, and returning a response that is associated with the matched request.

Imagine scanning through a dictionary for a word, and then looking up its definition. Hoverfly does exactly that, but the "word" is the HTTP request that was "captured" in :ref:`capture_mode`, 
and the "definition" is the response. 

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

By default, Hoverfly will set the type of match to "exact" for each field. Below is a Request Matcher set from an example Hoverfly simulation JSON file.

.. literalinclude:: ../../simulations/basic-simulation.json
   :lines: 4-24
   :linenos:
   :language: javascript

.. raw:: html

   <div>
        <p class="include-literal-footer">
            <a href="../../simulations/basic-simulation.json">See this request in its full simulation</a>
        </p>
   </div>

Not each of the fields is required, meaning it is possible to create partial request matchers that can be matched to more requests. For example, this request matcher will match any request to "docs.hoverfly.io".

.. code:: json

    {
        "destination": {
            "exactMatch": "docs.hoverfly.io"
        },
    }

Although the default matcher is "exactMatch", there are many other matchers to choose from.

.. todo:: Finish table, not sure its gonna be able to hold all the examples and still look good

+------------------------+------------------------------------+
| Request Field Matchers | Example                            |
+========================+====================================+
| "exactMatch"           | String-To-Match == String-To-Match |
+------------------------+------------------------------------+
| "xmlMatch"             | ?                                  |
+------------------------+------------------------------------+
| "xpathMatch"           | ?                                  |
+------------------------+------------------------------------+
| "jsonMatch"            | ?                                  |
+------------------------+------------------------------------+
| "jsonPathMatch"        | ?                                  |
+------------------------+------------------------------------+
| "regexMatch"           | ?                                  |
+------------------------+------------------------------------+
| "globMatch"            | ?                                  |
+------------------------+------------------------------------+

Request templates are defined in the :ref:`simulation_schema`.

With each Request Matcher is a Response. This is what Hoverfly will serve back to the client when a match is successful.

.. literalinclude:: ../../simulations/basic-simulation.json
   :lines: 25-32
   :linenos:
   :language: javascript

.. raw:: html

   <div>
        <p class="include-literal-footer">
            <a href="../../simulations/basic-simulation.json">See this response in its full simulation</a>
        </p>
   </div>

Since JSON does not support binary data, binary responses are base64 encoded. This is denoted by the encodedBody field. Hoverfly automatically encodes and decodes the data during the export and import phases.

.. literalinclude:: ../../simulations/basic-encoded-simulation.json
   :lines: 27-28
   :linenos:
   :language: javascript

.. raw:: html

   <div>
        <p class="include-literal-footer">
            <a href="../../simulations/basic-simulation.json">See this response in its full simulation</a>
        </p>
   </div>