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

.. code:: json

    {
        "scheme": {
            "exactMatch": "http"
        },
        "method": {
            "exactMatch": "GET"
        },
        "destination": {
            "exactMatch": "docs.hoverfly.io"
        },
        "path": {
            "exactMatch": "/pages/keyconcepts/templates.html"
        },
        "query": {
            "exactMatch": "query=true"
        },
        "body": {
            "exactMatch": "",
        },
        "headers": {}
    }

Not each of the fields is required, meaning it is possible to create partial request matchers that can be matched to more requests. For example, this request matcher will match any request to "docs.hoverfly.io".

.. code:: json

    {
        "destination": {
            "exactMatch": "docs.hoverfly.io"
        },
    }

Each field you want to match again may include one of the several matchers

Request templates are defined in the :ref:`simulation_schema`.

	ExactMatch - This will be an exact match (capturing requests will use this)
      String-To-Match -> String-To-Match

	XmlMatch - This will be an exact XML match (capturing request bodies with xml content-type will use this)
     <xml><documents><document></document></documents></xml> -> <xml><documents ><document ></document ></documents ></xml>

	XpathMatch - This will execute an Xpath expression, matches if successful
      ?

	JsonMatch - This will be an exact JSON match (capturing request bodies with json content-type will use this)
      ?

	JsonPathMatch - This will execute an Json path expression, matches if successful
      ?

	RegexMatch - This will execute an regex expression, matches if successful
      String-To-Match ->

	GlobMatch
      String-To-Match -> String-*, *-To-Match, *