.. _post_serve_action_tutorial:

Adding Post Serve Action to particular response
===============================================

.. seealso::

    Please carefully read through :ref:`post_serve_action` alongside this tutorial to gain a high-level understanding of what we are about to cover.

PostServeAction allows you to execute custom code after a response has been served in simulate or spy mode.

In this tutorial, we will help you to setup post serve action and execute it after particular response is served.

Let's begin by writing our post serve action. This script just makes HTTP call to find outbound IP.

We can make call to any other URL as well in order to send webhook or initiate some processing once response is served. Save the following as ``post_serve_action.py``:

.. literalinclude:: post_serve_action.py
    :language: python

Start Hoverfly and register post serve action

.. code:: bash

    hoverctl start
    hoverctl post-serve-action set --binary python3 --name outbound-http-call --script <path to script directory>/post_serve_action.py --delay #delay_in_ms

Now, you need to import simulation file containing post serve action details in response part.  Once done, you need to make API call that you stubbed. Custom code will get executed with delay provided post serving the response.

.. code:: json
    {
        "data": {
            "pairs": [
                {
                    "request": {
                        ...
                        "destination": [
                            {
                                "matcher": "exact",
                                "value": "helloworld-test.com"
                            }
                        ]
                        ...
                    },
                    "response": {
                        "status": 200,
                        "postServeAction": "outbound-http-call",
                        "body": "Hello World",
                        "encodedBody": false,
                        ...
                    }
                }
            ],
        ...
        },
        "meta": {
            "schemaVersion": "v5.2",
            "hoverflyVersion": "v1.6.0",
            "timeExported": "2023-09-02T13:10:04+05:30"
        }
    }





