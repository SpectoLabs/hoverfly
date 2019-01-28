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
