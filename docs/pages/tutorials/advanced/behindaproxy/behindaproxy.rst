.. _behind_a_proxy:

Using Hoverfly behind a proxy
================================

In some environments, you may only be able to access the internet via a proxy. For example,
your organization may route all traffic through a proxy for security reasons.

If this is the case, you will need to configure Hoverfly to work with the 'upstream' proxy.  

This configuration value can be easily set when starting an instance of Hoverfly.

For example, if the 'upstream' proxy is running on port ``8080`` on host ``corp.proxy``:  

.. code:: bash

    hoverctl start --upstream-proxy http://corp.proxy:8080

Upstream proxy authentication
-----------------------------

If the proxy you are using uses HTTP basic authentication, you can provide the authentication credentials as part of the upstream proxy configuration setting.

For example:

.. code:: bash

    hoverctl start --upstream-proxy http://my-user:my-pass@corp.proxy:8080
   
Currently, HTTP basic authentication is the only supported authentication method for an authenticated proxy.
