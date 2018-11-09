.. _lognormal:

Log-normal distributions of delay
=================================

The `Log-normal distribution`_ is a pretty accurate description of a server latency.
The Log-normal distribution defines by 2 parameters μ and σ. We will compute these parameters from `mean`_ and `median`_
of a server response time. These values you can see in your monitoring of the production server.
If need you can adjust response time by min and max parameters.

.. _Log-normal distribution: https://en.wikipedia.org/wiki/Log-normal_distribution
.. _mean: https://en.wikipedia.org/wiki/Expected_value
.. _median: https://en.wikipedia.org/wiki/Median
Let's apply a random log-normal distributed delay to all responses. First, we need to create and export a simulation.

.. literalinclude:: delays-capture.sh
   :language: sh

Take a look at the ``"globalActions"`` property within the ``simulation.json`` file you exported. It
should look like this:

.. literalinclude:: ../../../../simulations/basic-simulation.json
   :lines: 58-60
   :linenos:
   :language: javascript

Edit the file so the ``"globalActions"`` property looks like this:

.. literalinclude:: ../../../../simulations/log-normal-delay-simulation.json
   :lines: 26-37
   :linenos:
   :language: javascript

Hoverfly will apply a delay to all URLs that match the ``"urlPattern"`` value. We want
the delay to be applied to **all URLs**, so we set the ``"urlPattern"`` value to the regular expression ``"."``.

Now import the edited ``simulation.json`` file, switch Hoverfly to Simulate mode and make the requests
again.

.. literalinclude:: delays-simulate.sh
   :language: sh

