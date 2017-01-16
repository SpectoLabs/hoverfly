.. _webservertutorial:

Running Hoverfly as a webserver
===============================

.. seealso::

    Please carefully read through :ref:`webserver` alongside this tutorial to gain a high-level understanding of what we are about to cover.

Below you'll find a complete example how to capture data with Hoverfly running as a proxy, and how to save it in a simulation file.

.. literalinclude:: asawebserver.sh
    :language: bash
    :lines: 3-7

Now let's start Hoverfly as a webserver in Simulate mode.

.. literalinclude:: asawebserver.sh
    :language: bash
    :lines: 9-12

Hoverfly simulated our request successfully, but it did so as a webserver, not as a proxy.

.. note::

  Hoverfly starts in Simulate mode by default.
