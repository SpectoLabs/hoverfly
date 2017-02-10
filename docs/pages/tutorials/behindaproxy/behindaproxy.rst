.. _importing_simulation:

Using Hoverfly behind a proxy
================================

In some environments, you may already be running behind a proxy. Without configuration,
Hoverfly will not be able to forward requests.

Hoverfly has the ability to configure an upstream proxy. This configuration value
can be easily set when starting an instance of Hoverfly.

.. code:: bash

    hoverctl start --upstream-proxy http://localhost:8080

Currently, Hoverfly will only work with unauthenticated proxies.
