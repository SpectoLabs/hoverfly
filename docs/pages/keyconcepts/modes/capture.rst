Capture Mode
~~~~~~~~~~~~

Capture mode is used for creating API simulations. 

.. figure:: capture.mermaid.png

In capture mode, Hoverfly (running as a proxy server) intercepts communication between the client application and the external API. It transparently records outgoing requests from the client and records the incoming responses from the API.
 
Most commonly, requests to the external API are triggered by running automated tests against the application that consumes the API. During subsequent test runs, Hoverfly can be set to run in simulate mode, removing the dependency on the real external API. Alternatively, requests can be generated using a manual process. 

Usually, capture mode is used as the starting point in the process of creating a simulation. Captured data is then exported and modified before being re-imported into Hoverfly for use as a simulation.

Hoverfly cannot be set to capture mode when running as a webserver (see Hoverfly as a webserver).
