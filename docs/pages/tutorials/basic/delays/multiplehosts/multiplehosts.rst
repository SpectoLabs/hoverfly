Applying different delays based on host
=======================================

Now let's apply a delay of 1 second on responses from ``time.jsontest.com`` and a delay of 2 seconds on responses from ``date.jsontest.com``.

Run the following to create and export a simulation.

.. literalinclude:: delays-capture.sh
   :language: sh

Edit the ``simulation.json`` file so that the ``"globalActions"`` property looks like this:

.. literalinclude:: ../../../../simulations/multiple-hosts-delay-simulation.json
   :lines: 26-39
   :linenos:
   :language: javascript

Now run the following to import the edited ``simulation.json`` file and run the simulation:

.. literalinclude:: delays-simulate.sh
   :language: sh

You should notice a 1 second delay on responses from ``time.jsontest.com``, and a 2 second delay on responses from ``date.jsontest.com``.


.. note::

  You can easily get into a situation where your request URL has multiple matches. In this case, the first successful match wins.
