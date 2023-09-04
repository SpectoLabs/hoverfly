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
    hoverctl post-serve-action set --binary python3 --name callback-script --script <path to script directory>/post_serve_action.py --delay #delay_in_ms


Once post serve action is registered, you can check registered post serve action using below hoverctl command.

.. code::bash

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

Copy below simulation JSON content in a file ``simulation.json``:

.. code::json
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
              "headers": {
                "Access-Control-Allow-Origin": [
                  "*"
                ],
                "Content-Type": [
                  "application/json"
                ],
                "Date": [
                  "Mon, 04 Sep 2023 06:20:32 GMT"
                ],
                "Hoverfly": [
                  "Was-Here"
                ],
                "Server": [
                  "Google Frontend"
                ],
                "X-Cloud-Trace-Context": [
                  "9325c9ca551f725586fc96cc65a4aae0"
                ]
              },
              "templated": false,
              "postServeAction": "callback-script"
            }
          }
        ],
        "globalActions": {
          "delays": [

          ],
          "delaysLogNormal": [

          ]
        }
      },
      "meta": {
        "schemaVersion": "v5.2",
        "hoverflyVersion": "v1.5.3",
        "timeExported": "2023-09-04T11:50:40+05:30"
      }
    }

 Run below hoverctl command to import simulation file.

.. code:: bash
    hoverctl import <path-to-simulation-file>

Once done, make a curl call to date.jsontest.com using below curl call.

.. code::bash
    curl --proxy http://localhost:8500 http://date.jsontest.com

Check the logs using hoverctl that post serve action was invoked. You will see the message - `Output from post serve action HTTP call invoked from IP Address`.

.. code::bash

    hoverctl logs



