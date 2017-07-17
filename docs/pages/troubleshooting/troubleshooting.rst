.. _troubleshooting:

Troubleshooting
===============

Why isn't Hoverfly matching my request?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

When Hoverfly cannot match a response to an incoming request, it will return information on the closest match:

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

Here, you can see which fields did not match. In this case, it was the ``body``. 
You can also view this information by running ``hoverctl logs``.

Why isn't Hoverfly returning the closest match when it cannot match a request?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Hoverfly will only provide this information when the matching strategy is set to **strongest match** (the default). If you
are using the **first match** matching strategy, the closet match information will not be returned.  

How can I view the Hoverfly logs?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. code:: bash

    hoverctl logs
