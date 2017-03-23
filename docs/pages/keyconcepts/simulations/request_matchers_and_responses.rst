.. request_matchers_and_responses:

Request Matchers and Responses
==============================

A request matcher is a request within  the simulation that is used to match against incoming requests. A request matcher may include each of the fields available on a request. These are

::
    {
        "scheme": "http",
        "method": "GET",
        "destination": "docs.hoverfly.io",
        "path": "/pages/keyconcepts/templates.html",
        "query": "query=true",
        "body": "",
        "headers": {}
    }
    



Request templates are defined in the :ref:`simulation_schema`.

	ExactMatch
	XmlMatch
	XpathMatch
	JsonMatch
	JsonPathMatch
	RegexMatch
	GlobMatch

