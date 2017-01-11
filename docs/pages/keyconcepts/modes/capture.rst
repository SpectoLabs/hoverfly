.. _capture_mode:

Capture mode
************

Capture Mode is used for creating API simulations.

.. figure:: capture.mermaid.png

In Capture Mode, Hoverfly (running as a proxy server - see :ref:`proxy_server`) intercepts communication between the client application and the external service. It transparently records outgoing requests from the client and records the incoming responses from the service API.

Most commonly, requests to the external service API are triggered by running automated tests against the application that consumes the API. During subsequent test runs, Hoverfly can be set to run in :ref:`simulate_mode`, removing the dependency on the real external service. Alternatively, requests can be generated using a manual process.

Usually, Capture Mode is used as the starting point in the process of creating an API simulation. Captured data is then exported and modified before being re-imported into Hoverfly for use as a simulation.

.. note::

    Hoverfly cannot be set to capture mode when running as a webserver (see :ref:`webserver`).
