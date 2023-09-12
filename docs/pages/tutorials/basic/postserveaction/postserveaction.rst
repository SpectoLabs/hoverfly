.. _post_serve_action_tutorial:

Adding Post Serve Action to particular response
===============================================

.. seealso::

    Please carefully read through :ref:`post_serve_action` alongside this tutorial to gain a high-level understanding of what we are about to cover.

PostServeAction allows you to execute custom code after a response has been served in simulate or spy mode.

In this tutorial, we will help you to setup post serve action and execute it after particular response is served.

Let's begin by writing our post serve action. The following python script makes an HTTP call to http://ip.jsontest.com and prints out your IP address.

We can make a call to any other URL as well in order to send webhook or initiate some processing once response is served.
Although the example script is written in python, you can also write an action using any other language as long as you have provided a binary in your local environment to execute the code.

Save the following as ``post_serve_action.py``:

.. literalinclude:: post_serve_action.py
    :language: python

Start Hoverfly and register a post serve action:

.. code:: bash

    hoverctl start
    hoverctl post-serve-action set --binary python3 --name callback-script --script <path to script directory>/post_serve_action.py --delay 3000


Once the post serve action is registered, you can confirm using below hoverctl command.

.. code:: bash

    hoverctl post-serve-action get-all

    Sample output
    +-----------------+---------+--------------------------------+-----------+
    |   ACTION NAME   | BINARY  |             SCRIPT             | DELAY(MS) |
    +-----------------+---------+--------------------------------+-----------+
    | callback-script | python3 | #!/usr/bin/env python import   |      3000 |
    |                 |         | sys import logging import      |           |
    |                 |         | random from time import sleep  |           |
    |                 |         | ...                            |           |
    +-----------------+---------+--------------------------------+-----------+

Copy this simulation JSON content to a file called ``simulation.json``:

.. code:: json

    {
      "data": {
        "pairs": [
          {
            "request": {
              "path": [
                {
                  "matcher": "exact",
                  "value": "/"
                }
              ],
              "method": [
                {
                  "matcher": "exact",
                  "value": "GET"
                }
              ],
              "destination": [
                {
                  "matcher": "exact",
                  "value": "date.jsontest.com"
                }
              ],
              "scheme": [
                {
                  "matcher": "exact",
                  "value": "http"
                }
              ],
              "body": [
                {
                  "matcher": "exact",
                  "value": ""
                }
              ]
            },
            "response": {
              "status": 200,
              "body": "01-01-1111",
              "encodedBody": false,
              "templated": false,
              "postServeAction": "callback-script"
            }
          }
        ]
      },
      "meta": {
        "schemaVersion": "v5.2",
        "hoverflyVersion": "v1.5.3",
        "timeExported": "2023-09-04T11:50:40+05:30"
      }
    }


Run this hoverctl command to import the simulation file.

.. code:: bash

    hoverctl import <path-to-simulation-file>

The simulation sets hoverfly to return a successful response and 3 seconds after that invokes the "callback-script" action.
You can try it out by making the following request to http://date.jsontest.com using cURL.

.. code:: bash

    curl --proxy http://localhost:8500 http://date.jsontest.com

You should see the message in the hoverfly logs - `Output from post serve action HTTP call invoked from IP Address`.

.. code:: bash

    hoverctl logs



