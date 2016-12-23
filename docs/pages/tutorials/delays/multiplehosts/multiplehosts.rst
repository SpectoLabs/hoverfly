Delays on multiple hosts
........................

To apply a delay of 1 second to ``time.jsontest.com`` and a delay of 2 seconds to ``date.jsontest.com``, once again save the following inside ``delays.json``.

.. literalinclude:: delays.json
   :language: json

Now run the following:

.. literalinclude:: delays.sh
   :language: json

Once again, you should notice a 1 second delay taking place on ``time.jsontest.com``, and a 2 seconds delay taking place on ``date.jsontest.com``.


