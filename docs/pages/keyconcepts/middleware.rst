.. _middleware:

Middleware
==========

.. todo:: Rewrite this so its a bit more generic and then talk about the different ways of running middleware; webserver, or locally with a binary or a script or both

Middleware intercepts traffic between the client and the API (whether real or simulated), and allowing you to manipulate it.

You can use middleware to manipulate data in simulated responses, or to inject unpredictable performance characteristics into your simulation.

Middleware works differently depending on the Hoverfly mode.

- Capture mode: middleware affects only **outgoing requests**
- Simulate/Spy mode: middleware affects only **responses** (cache contents remain untouched)
- Synthesize mode: middleware **creates responses**
- Modify mode: middleware affects **requests and responses**

.. note::

    Middleware is applied after rendering the templating functions (see :ref:`templating`) in the response body.


You can write middleware in any language. There are two different types of middleware.

Local Middleware
----------------
Hoverfly has the ability to invoke middleware by executing a script or binary file on a host operating system.
The only requires are that the provided middleware can be executed and sends the Middleware JSON schema to stdout
when the Middleware JSON schema is received on stdin.

HTTP Middleware
---------------
Hoverfly can also send middleware requests to a HTTP server instead of running a process locally. The benefits of this
are that Hoverfly does not initiate the process, giving more control to the user. The only requirements are that Hoverfly can
POST the Middleware JSON schema to middleware URL provided and the middleware HTTP server responses with a 200 and the
Middleware JSON schema is in the response.


Middleware Interface
--------------------

When middleware is called by Hoverfly, it expects to receive and return JSON (see :ref:`simulation_schema`). Middleware can be used to modify the values in the JSON but **must not** modify the schema itself.

.. figure:: middleware.mermaid.png

Hoverfly will send the JSON object to middleware via the standard input stream. Hoverfly will then listen to the standard output stream and wait for the JSON object to be returned.


.. seealso::

    Middleware examples are covered in the tutorials section. See :ref:`randomlatency` and :ref:`modifyingresponses`.


Security and availability of the Set Middleware API
---------------------------------------------------

By default, the admin endpoint to set middleware (PUT /api/v2/hoverfly/middleware) is disabled. To enable it:

- When starting Hoverfly directly: run with the flag: -enable-middleware-api
- When starting via hoverctl: use the same flag on start: hoverctl start --enable-middleware-api

Network binding and remote access
---------------------------------

By default, Hoverfly binds its Admin and Proxy ports to the loopback interface only (127.0.0.1). This means the Admin API is not reachable from remote hosts out of the box.

.. warning::

   Exposing the Admin API outside localhost increases risk, especially if the Set Middleware API is enabled, because it allows executing arbitrary scripts/binaries on the host (for local middleware) or invoking remote middleware services.

   If you expose the Admin API and enable the Set Middleware API, you should:

   - Run Hoverfly only on trusted/private networks.
   - Restrict access to the Admin API to trusted callers and networks (e.g., via firewalls, security groups, VPNs, reverse proxy ACLs).
   - Prefer binding to localhost unless there is a strong need to expose it, and scope exposure to the minimum required interfaces.
   - :ref:`proxyauth` if appropriate and avoid exposing the Admin port publicly.

The guidance above applies whether you configure middleware as a local executable/script or as HTTP middleware.
