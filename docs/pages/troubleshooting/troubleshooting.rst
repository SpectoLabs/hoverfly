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

Hoverfly will only provide this information when the matching strategy is set to **strongest match** 
(the default). If you are using the **first match** matching strategy, the closet match information 
will not be returned.

How can I view the Hoverfly logs?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. code:: bash

    hoverctl logs


Why does my simulation have a ``deprecatedQuery`` field?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Older simulations that have been upgraded through newer versions of Hoverfly may now contain a field 
on requests called ``deprecatedQuery``. With the v5 simulation schema, the request query field was
updated to more fully represent request query paramters. This involves storing queries based on
query keys, similarly to how headers are stored in a simulation.

Currently the ``deprecatedQuery`` field will work and works alongside the ``query`` field and support
for this field will eventually be dropped.

If you have ``deprecatedQuery`` field, you should remove it by splitting it by query keys.


.. code:: json

    "deprecatedQuery": "page=20&pageSize=15"

.. code:: json

    "query": {
        "page": [
            {
                "matcher": "exact",
                "value": "20"
            }
        ],
        "pageSize": [
            {
                "matcher": "exact",
                "value": "15"
            }
        ],
    }

If you cannot update your ``deprecatedQuery`` from your simulation for a technical reason, feel free to 
raise an issue on Hoverfly.

Why am I not able to access my Hoverfly remotely?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

That's because Hoverfly is bind to loopback interface by default, meaning that you can only access 
to it on localhost. To access it remotely, you can specify the IP address it listens on. For example, 
setting ``0.0.0.0`` to listen on all network interfaces.

.. code:: bash

    hoverfly -listen-on-host 0.0.0.0

