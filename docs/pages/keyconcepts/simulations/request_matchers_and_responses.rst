.. request_matchers_and_responses:

Request Matchers and Responses
==============================

Hoverfly simulates APIs by `matching` responses to incoming requests.

Imagine scanning through a dictionary for a word, and then looking up its definition. Hoverfly does exactly that, but the "word" is the HTTP request that was "captured" in :ref:`capture_mode`.

.. code:: json

    {
        "scheme": "http",
        "method": "GET",
        "destination": "docs.hoverfly.io",
        "path": "/pages/keyconcepts/templates.html",
        "query": "query=true",
        "body": "",
        "headers": {}
    }

These captured requests are translated into Request Matchers. This request consists all of the same fields as a request but uses matchers instead of exact values.

.. literalinclude:: ../../simulations/basic-simulation.json
   :lines: 4-24
   :linenos:
   :language: javascript

Not each of the fields is required, meaning it is possible to create partial request matchers that can be matched to more requests. For example, this request matcher will match any request to "docs.hoverfly.io".

.. code:: json

    {
        "destination": {
            "exactMatch": "docs.hoverfly.io"
        },
    }

Although the default matcher is "exactMatch", there are many other matchers to choose from.

.. todo:: Table for matchers?

.. todo:: ExactMatch - This will be an exact match (capturing requests will use this) String-To-Match -> String-To-Match

.. todo:: ExactMatch - This will be an exact match (capturing requests will use this) String-To-Match -> String-To-Match

.. todo:: XmlMatch - This will be an exact XML match (capturing request bodies with xml content-type will use this) <xml><documents><document></document></documents></xml> -> <xml><documents ><document ></document ></documents ></xml>

.. todo:: XpathMatch - This will execute an Xpath expression, matches if successful

.. todo:: JsonMatch - This will be an exact JSON match (capturing request bodies with json content-type will use this)

.. todo:: JsonPathMatch - This will execute an Json path expression, matches if successful

.. todo:: RegexMatch - This will execute an regex expression, matches if successful | String-To-Match ->

.. todo:: GlobMatch | String-To-Match -> String-*, *-To-Match, *

Request templates are defined in the :ref:`simulation_schema`.

.. literalinclude:: ../../simulations/basic-simulation.json
   :lines: 25-32
   :linenos:
   :language: javascript

.. literalinclude:: ../../simulations/basic-encoded-simulation.json
   :lines: 27-28
   :linenos:
   :language: javascript