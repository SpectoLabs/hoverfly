.. _webservertutorial:

Running Hoverfly as a webserver
===============================

.. seealso::

    Please carefully read through :ref:`webserver` alongside this tutorial to gain a high-level understanding of what we are about to cover.

Below is a complete example how to capture data with Hoverfly running as a proxy, and how to save it in a simulation file.

.. literalinclude:: asawebserver.sh
    :language: bash
    :lines: 3-7

Now we can use Hoverfly as a **webserver** in Simulate mode.

.. literalinclude:: asawebserver.sh
    :language: bash
    :lines: 9-12

Hoverfly returned a response to our request while running as a webserver, not as a proxy. 

Notice that we issued a cURL command to ``http://localhost:8500/a/b`` instead of ``http://echo.jsontest.com/a/b``.
This is because when running as a webserver, Hoverfly strips the domain from the endpoint URL in 
the simulation. 

This is explained in more detail in the :ref:`webserver` section. 

.. note::

  Hoverfly starts in Simulate mode by default.
