.. _corstutorial:

Enable CORS support on Hoverfly
===============================

By enabling CORS (Cross-Origin Resource Sharing) support, your web application running on the browser can make requests to
Hoverfly even it's not on the same domain.

Starting Hoverfly with CORS enabled is simple:

.. code:: bash

    hoverfly -cors


Or using `hoverctl`:

.. code:: bash

    hoverctl start --cors


You can check if CORS is enabled on the Hoverfly by querying the status:

.. code:: bash

    hoverctl status


When CORS is enabled, Hoverfly intercepts any pre-flight request, and returns an empty 200 response with the following default CORS headers:

- Access-Control-Allow-Origin: (same value as the ``Origin`` header from the request)
- Access-Control-Allow-Methods: GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS
- Access-Control-Allow-Headers: (same value as the ``Access-Control-Request-Headers`` header from the request)
- Access-Control-Max-Age: 1800
- Access-Control-Allow-Credentials: true

Hoverfly also intercepts the actual CORS requests, and add the following default CORS headers to the response:

- Access-Control-Allow-Origin: (same value as the ``Origin`` header from the request)
- Access-Control-Allow-Credentials: true

Support for customizing the CORS headers will be added in the future release.

.. note::

    Two points to notice when Hoverfly is in capture mode and CORS is enabled:
    1. Pre-flight requests handling and CORS headers provided by Hoverfly are not recorded in the simulation.
    2. Hoverfly preserves the CORS headers from the remote server if they are present.


