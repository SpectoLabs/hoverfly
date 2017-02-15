Applying different delays based on HTTP method
==============================================

Let's apply a delay of 2 seconds on responses to **GET requests only** made to ``echo.jsontest.com/b/c``.

Run the following to create and export a simulation.

.. literalinclude:: delays-capture.sh
   :language: sh

Edit the ``simulation.json`` file so that the ``"globalActions"`` property looks like this:

.. literalinclude:: global-actions-with-delay-method.json
   :language: javascript

Now run the following to import the edited ``simulation.json`` file and run the simulation:

.. literalinclude:: delays-simulate.sh
   :language: sh

You should notice a 2 second delay on the response to the GET request and no delay on the response to the POST request.
