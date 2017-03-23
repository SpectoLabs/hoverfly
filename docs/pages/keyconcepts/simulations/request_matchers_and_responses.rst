.. request_matchers_and_responses:

Request Matchers and Responses
==============================

A request matcher is a request within  the simulation that is used to match against incoming requests. A request matcher may include each of the fields available on a request. These are

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