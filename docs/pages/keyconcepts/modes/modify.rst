.. _modify_mode:

Modify mode
***********

Modify mode is similar to :ref:`capture_mode`. Hoverfly will forward requests to the intended destinations, but it will not record the requests or the responses. Instead, it will pass each request to a :ref:`middleware` executable before forwarding it to the destination, and will do the same with the response before returning it to the client.

.. figure:: modify.mermaid.png

You could use this mode to “man in the middle” your own requests and responses. For example, you could change the API key you are using to authenticate against a third-party API.
