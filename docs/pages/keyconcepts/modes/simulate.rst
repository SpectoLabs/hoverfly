.. _simulate_mode:

Simulate mode
=============

In this mode, Hoverfly uses its simulation data in order to simulate external APIs. Each time Hoverfly receives a request,
rather than forwarding it on to the real API, it will respond instead. No network traffic will ever reach the real external API.

.. figure:: simulate.mermaid.png

The simulation can be produced automatically via by running Hoverfly in :ref:`capture_mode`, or created manually. See :ref:`simulations` for information.
