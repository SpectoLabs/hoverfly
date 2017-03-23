.. _simulate_mode:

Simulate mode
=============

.. todo:: @tjcunliffe to review

In this mode, Hoverfly uses traffic captured using :ref:`capture_mode` (which may also have been manually edited) to mimic external APIs.

.. figure:: simulate.mermaid.png

Each time Hoverfly receives a request from the application, instead of forwarding it to the intended destination, Hoverfly will resolve it. It will first look in a cache (which is configurable). This cache is populated with successful requests and responses that Hoverfly has previously served. If Hoverfly has a matching request in the cache, it will serve the response.

If the request is not found in the cache, the request is then compared against all the requests in the simulation.
