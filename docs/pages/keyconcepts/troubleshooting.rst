.. _troubleshooting:

Troubleshooting
===============

1. Why does Hoverfly not match my request?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

When Hoverfly misses, it will return the closest match in the response body of the request like so:

::

    Hoverfly Error!

    There was an error when matching

    Got error: Could not find a match for request, create or record a valid matcher first!

    The following request was made, but was not matched by Hoverfly:

    {
        "Path": "/closest-miss",
        "Method": "GET",
        "Destination": "destination.com",
        "Scheme": "http",
        "Query": "",
        "Body": "",
        "Headers": {
            "Accept-Encoding": [
                "gzip"
            ],
            "User-Agent": [
                "Go-http-client/1.1"
            ]
        }
    }

    The matcher which came closest was:

    {
        "path": {
            "exactMatch": "/closest-miss"
        },
        "destination": {
            "exactMatch": "destination.com"
        },
        "body": {
            "exactMatch": "body"
        }
    }

    But it did not match on the following fields:

    [body]

    Which if hit would have given the following response:

    {
        "status": 200,
        "body": "",
        "encodedBody": false,
        "headers": null
    }`

From this, you are told which fields did not match. In the above case, it was the body. You can also view this information through hoverctl logs.

2. Why doesn't Hoverfly returning the closest matcher when there is a miss?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Hoverfly will only have this information when the matching strategy is "strongest". If you are in simulate mode and have set the matching strategy
to "first" then the information will not be available.

3. How can I access the logs?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. code:: bash

    hoverctl logs
