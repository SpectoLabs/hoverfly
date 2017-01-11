.. _middleware:

Middleware
**********

When you want to modify a request, a response, or both you need *middleware*. Middleware intercepts data flowing between the client and the endpoint (whether real or virtualised), and enables you to modify it at will.

Below is a list of the modes and how they interact with middleware.

- Capture Mode: middleware affects only **outgoing** requests
- Simulate Mode: middleware affects only **responses** (cache contents remain untouched)
- Synthesize Mode: middleware **creates responses**
- Modify Mode: middleware affects **requests and responses**

You can write middleware in any language as long as it can be executed by the shell on the Hoverfly host. This includes executing a binary file, or an interpreter executing a script. e.g:

::

    php middleware.php

Middleware Interface
~~~~~~~~~~~~~~~~~~~~

When middleware is called by Hoverfly, expect to receive and return JSON. This JSON will match the request response pair structure defined in the simulation JSON schema. It is important that middleware modifies the values of this JSON schema and not the schema itself.

.. figure:: middleware.mermaid.png

Hoverfly will send the JSON object to middleware via the standard input stream. Hoverfly will then listen to the standard output stream and wait for the JSON object to be returned.


.. seealso::

    Middleware examples are covered in the tutorials section.
