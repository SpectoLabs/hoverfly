Simulate Mode
~~~~~~~~~~~~~

In this mode, Hoverfly uses either traffic captured previously, or manually created simulation JSON files to mimic external APIs.

.. figure:: simulate.mermaid.png

Each time Hoverfly receives a request from the application, instead of forwarding it to the intended destination, it looks in the simulation data for a matching response. If it finds a match, it returns the response to the application.

This matching strategy won’t always be appropriate however. In some cases you may want Hoverfly to return a single response for a number of different possible requests. This can be done by editing the simulation data (see Templates). 


