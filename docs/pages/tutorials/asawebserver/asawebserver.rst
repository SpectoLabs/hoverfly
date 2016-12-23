.. _webservertutorial:

Running Hoverfly as a webserver
-------------------------------

.. seealso::
    
    Please carefully read through :ref:`webserver` alongside this tutorial, to gain a high-level understanding of what we are about to cover.

Below you'll find a complete example of capturing data with Hoverfly in proxy mode, and saving it in a simulation file.

.. literalinclude:: asawebserver.sh
    :language: bash
    :lines: 3-7

We've now captured our request. Now let's start Hoverfly in webserver mode, and simulate.

.. literalinclude:: asawebserver.sh
    :language: bash
    :lines: 9-12

Hoverfly simulated our request successfully, but it did so as a webserver, and not as a proxy.