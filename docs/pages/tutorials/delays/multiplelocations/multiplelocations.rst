Applying different delays based on URI
======================================

Now let's apply different delays based on location. Run the following to create and export a simulation.

.. literalinclude:: delays-capture.sh
   :language: sh

Edit the ``simulation.json`` file so that the ``"globalActions"`` property looks like this:

.. literalinclude:: global-actions-with-delay-location.json
   :language: javascript

Now run the following to import the edited ``simulation.json`` file and run the simulation:

.. literalinclude:: delays-simulate.sh
   :language: sh

You should notice a 2 second delay on responses from ``echo.jsontest.com/a/b`` and ``echo.jsontest.com/b/c``, and a 3 second delay on the response from ``echo.jsontest.com/c/d``.