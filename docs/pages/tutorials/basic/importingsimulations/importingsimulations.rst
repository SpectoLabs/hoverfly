.. _importing_simulation:

Importing and using a simulation
================================

In this tutorial we are going to import the simulation we created in the previous tutorial.

.. code:: bash

    hoverctl start
    hoverctl import simulation.json

Hoverfly can also import simulation data that is stored on a remote host via HTTP:

.. code:: bash

    hoverctl import https://example.com/example.json

Make a request with cURL, using Hoverfly as a proxy.

.. code:: bash

    curl --proxy localhost:8500 http://time.jsontest.com

This outputs the time at the time the request was captured.

.. code:: bash

    {
       "time": "02:07:28 PM",
       "milliseconds_since_epoch": 1482242848562,
       "date": "12-20-2016"
    }

Stop Hoverfly:

.. code:: bash

    hoverctl stop
