.. _webserver:

Hoverfly as a webserver
=======================

Sometimes you may not be able to configure your client to use a proxy, or you may want to explicitly point your application at Hoverfly. For this reason, Hoverfly can run as a webserver.

.. figure:: webserver.mermaid.png

.. note::

    When running as a webserver, Hoverfly cannot capture traffic (see :ref:`capture_mode`) - it can only be used to simulate and synthesize APIs (see :ref:`simulate_mode`, :ref:`modify_mode` and :ref:`synthesize_mode`). For this reason, when you use Hoverfly as a webserver, you should have Hoverfly simulations ready to be loaded.

When running as a webserver, Hoverfly strips the domain from the endpoint URL. For example, if you made requests to the following URL while capturing traffic with Hoverfly running as a proxy:

.. code::

      http://echo.jsontest.com/key/value

And Hoverfly is running in simulate mode as a webserver on:

.. code::

      http://localhost:8888

Then the URL you would use to retrieve the data from Hoverfly would be:

.. code::

      http://localhost:8500/key/value

.. seealso::

    Please refer to the :ref:`webservertutorial` tutorial for a step-by-step example.
