.. _modifyingresponses:

Using middleware to modify response payload and status code
===========================================================

.. seealso::

    Please carefully read through :ref:`middleware` alongside these tutorials to gain a high-level understanding of what we are about to cover.

We will use a python script to modify the body of a response and randomly change the status code.

Let's begin by writing our middleware. Save the following as ``middleware.py``:

.. literalinclude:: middleware.py
    :language: python

The middleware script randomly toggles the status code between 200 and 201, and changes the response body to a dictionary containing ``{'foo':'baz'}``.

.. literalinclude:: modify.sh
    :language: sh

As you can see, middleware allows you to completely modify the content of a simulated HTTP response.
