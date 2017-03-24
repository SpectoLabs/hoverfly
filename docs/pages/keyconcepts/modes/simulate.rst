.. _simulate_mode:

Simulate mode
=============

.. todo:: what is meant by "cache (which is configurable)"? What can be configured?

In this mode, Hoverfly uses traffic captured using :ref:`capture_mode` (which may also have been manually edited) to mimic external APIs.

.. figure:: simulate.mermaid.png

Each time Hoverfly receives a request from the application, instead of forwarding it to the intended destination, Hoverfly will resolve it. 
First it will look in a cache which is populated with requests and responses that it has served previously. If Hoverfly finds a matching request in the cache, 
it will serve the associated response.

If a matching request is not found in the cache, Hoverfly will look for a matching request within the simulation.
