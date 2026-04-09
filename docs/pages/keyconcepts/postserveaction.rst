.. _post_serve_action:

Post Serve Action
=================

Overview
--------

- PostServeAction allows you to execute custom code or invoke endpoint with request-response pair after a response has been served in simulate or spy mode.

- It is custom script that can be written in any language. Hoverfly has the ability to invoke a script or binary file on a host operating system or on the remote host. Custom code is executed/remote host is invoked after a provided delay(in ms) once simulated response is served.

- We can register multiple post serve actions.

- In order to register local post serve action, it takes mainly four parameters - binary to invoke script, script content/location, delay(in ms) post which it will be executed and name of that action.

- In order to register remote post serve action, it takes mainly three parameters - remote host to be invoked with request-response pair, delay(in ms) post which it will be executed and name of that action.

Ways to register a Post Serve Action
------------------------------------

- At time of startup by passing single/multiple -post-serve-action flag(s) as mentioned in the `hoverfly command page <https://docs.hoverfly.io/en/latest/pages/reference/hoverfly/hoverflycommands.html>`_.

- Via PUT API to register new post serve action as mentioned in the `API page <https://docs.hoverfly.io/en/latest/pages/reference/api/api.html>`_.

- Using command line hoverctl post-serve-action set command as mentioned in the `hoverctl command page <https://docs.hoverfly.io/en/latest/pages/reference/hoverctl/hoverctlcommands.html>`_.


- Once post serve action is registered, we can trigger particular post serve action by putting it in response part of request-response pair in simulation JSON.

**Example Simulation JSON**

.. code:: json

    {
        "response": {
            "postServeAction": "<name of post serve action we want to invoke>"
        }
    }


Choosing Between Local and Remote Execution
--------------------------------------------

Hoverfly supports two execution modes for post-serve actions. Understanding the
difference is important when running under significant load.

**Local execution**

Hoverfly forks a new subprocess for every request that triggers the action. The
binary or script is re-executed from scratch each time:

.. code::

    hoverfly -post-serve-action "my-action python3 /path/to/script.py 0"

This is the easiest option to get started, but it does not scale. At high request
rates (e.g. 40+ rps with a Python script), you will quickly accumulate dozens of
concurrent processes — each with its own interpreter startup cost, memory footprint,
and open connections. This commonly leads to OOMKilled in containerised environments.

**Remote execution (recommended for high throughput)**

Instead of forking a subprocess, Hoverfly makes an HTTP ``POST`` request to a
server that you run separately. That server stays alive indefinitely — only one
process, shared across all requests:

.. code::

    hoverfly -post-serve-action "my-action http://localhost:8080/trigger 0"

This means you can use async frameworks such as ``aiohttp`` or ``FastAPI`` in their
natural form: one running event loop handling many concurrent webhooks efficiently,
without paying the startup cost on every request.

Running a Remote Post Serve Action: Step by Step
--------------------------------------------------

**Step 1 — Write your server**

Your server must accept HTTP ``POST`` requests and return ``HTTP 200``. Hoverfly
considers any other status code a failure and logs an error. Here is a minimal
Python example:

.. code:: python

    # server.py
    from http.server import BaseHTTPRequestHandler, HTTPServer
    import json

    class Handler(BaseHTTPRequestHandler):
        def do_POST(self):
            length = int(self.headers.get("Content-Length", 0))
            payload = json.loads(self.rfile.read(length))

            # payload contains two keys: "request" and "response"
            # use them to implement your webhook/callback logic
            print(f"[action] path={payload['request']['path']}")

            self.send_response(200)
            self.end_headers()

        def log_message(self, format, *args):
            pass  # suppress default access log noise

    if __name__ == "__main__":
        print("Listening on :8080")
        HTTPServer(("", 8080), Handler).serve_forever()

Start it:

.. code::

    python3 server.py

**Step 2 — Understand the payload Hoverfly sends**

On every matching request, Hoverfly POSTs a JSON body to your server containing
the full request-response pair:

.. code:: json

    {
        "request": {
            "path": [{"matcher": "exact", "value": "/api/orders"}],
            "method": [{"matcher": "exact", "value": "POST"}],
            "destination": [{"matcher": "exact", "value": "example.com"}],
            "body": [{"matcher": "exact", "value": ""}],
            "headers": {}
        },
        "response": {
            "status": 200,
            "body": "Hello World",
            "encodedBody": false
        }
    }

Your server can read any field from this payload to drive its logic — for example,
extracting an order ID from the request body to send a webhook.

**Step 3 — Configure Hoverfly to call your server**

Pass the remote URL as the second token in ``-post-serve-action``:

.. code::

    hoverfly -post-serve-action "my-action http://localhost:8080 0" -import simulation.json

The format is: ``"<name> <url> <delay-ms>"``.

- ``my-action`` must match the ``postServeAction`` field in your simulation JSON.
- The URL must be reachable from Hoverfly at runtime.
- The delay (in milliseconds) is applied before Hoverfly calls the endpoint.

Alternatively, register it at runtime via hoverctl:

.. code::

    hoverctl post-serve-action set --name my-action --remote http://localhost:8080 --delay 0

**Step 4 — Confirm the action is registered**

.. code::

    curl http://localhost:8888/api/v2/hoverfly/post-serve-action

You should see your action listed with its remote URL.

**Running in Docker Compose**

When both Hoverfly and your action server run as containers, use the service name
as the hostname. Make sure Hoverfly starts after the action server is ready:

.. code:: yaml

    services:
      hoverfly:
        image: spectolabs/hoverfly
        command: >
          -post-serve-action "my-action http://postaction:8080 0"
          -import /app/simulation.json
        depends_on:
          - postaction

      postaction:
        build: ./postaction-server
        ports:
          - "8080:8080"


