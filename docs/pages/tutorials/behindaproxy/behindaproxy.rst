.. _behind_a_proxy:

Using Hoverfly behind a proxy
================================

In some environments, you may only be able to access the internet via a proxy. For example,
your organization may route all traffic through a proxy for security reasons.

If this is the case, you will need to configure Hoverfly to work with the 'upstream' proxy.  

This configuration value can be easily set when starting an instance of Hoverfly.

.. code:: bash

    hoverctl start --upstream-proxy http://localhost:8080

Currently, Hoverfly will only work with unauthenticated proxies.
