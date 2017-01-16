Applying different delays based on host
=======================================

To apply a delay of 1 second on responses from ``time.jsontest.com`` and a delay of 2 seconds on responses from ``date.jsontest.com``, save the following inside ``delays.json``.

.. literalinclude:: delays.json
   :language: javascript

Now run the following:

.. literalinclude:: delays.sh
   :language: sh

You should notice a 1 second delay on responses from ``time.jsontest.com``, and a 2 second delay on responses from ``date.jsontest.com``.


.. note::

  You can easily get into a situation where your request URL has multiple matches within your ``delays.json`` file. In this case, the first successful match wins.
