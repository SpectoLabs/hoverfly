Applying a delay to all responses
=================================

Let's apply a 2 second delay to all responses. First, we need to create and export a simulation.

.. literalinclude:: delays-capture.sh
   :language: sh

Take a look at the ``"globalActions"`` property within the ``simulation.json`` file you exported. It
should look like this:

.. literalinclude:: ../../../../simulations/basic-simulation.json
   :lines: 58-60
   :linenos:
   :language: javascript

Edit the file so the ``"globalActions"`` property looks like this:

.. literalinclude:: ../../../../simulations/basic-delay-simulation.json
   :lines: 26-34
   :linenos:
   :language: javascript

Hoverfly will apply a delay of 2000ms to all URLs that match the ``"urlPattern"`` value. We want
the delay to be applied to **all URLs**, so we set the ``"urlPattern"`` value to the regular expression ``"."``.

Now import the edited ``simulation.json`` file, switch Hoverfly to Simulate mode and make the requests
again.

.. literalinclude:: delays-simulate.sh
   :language: sh
   
The responses to both requests are delayed by 2 seconds.

