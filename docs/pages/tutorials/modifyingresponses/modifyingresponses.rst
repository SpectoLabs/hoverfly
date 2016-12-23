.. _modifyingresponses:

Using middleware to modify responses
------------------------------------

.. seealso::
    
    Please carefully read through :ref:`middleware` alongside these tutorials, to gain a high-level understanding of what we are about to cover.

Let's begin by writing our middleware, please save the following as ``middleware.py``:

.. literalinclude:: middleware.py
    :language: python

The middleware above randomly toggles the status code between 200, and 201. Moreover it sets the response's body to a dictionary containing ``{'foo':'baz'}``.

.. literalinclude:: modify.sh
    :language: sh

As you can see, middleware has the ability to completely modify the content of http traffic.