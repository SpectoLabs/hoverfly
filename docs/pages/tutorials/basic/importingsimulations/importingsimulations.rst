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


.. note:: Importing multiple simulations:

    The above example shows importing one simulation into Hoverfly. You can also import multiple simulations:

    .. code:: bash

        hoverctl simulation add foo.json bar.json


    You can specify one or more simulations when starting Hoverfly using ``hoverfly`` or ``hoverctl`` commands:

    .. code:: bash

        hoverfly -import foo.json -import bar.json
        hoverctl start --import foo.json --import bar.json

    Hoverfly appends any unique pair to the existing simulation by comparing the equality of the request JSON objects.
    If a conflict occurs, the pair is not added.