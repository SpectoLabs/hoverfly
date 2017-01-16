.. _simulations:

Simulations
===========

The core functionality of Hoverfly is to capture HTTP(S) traffic to create API simulations which can be used in testing. Hoverfly stores captured traffic as `simulations`.

Simulations can be written to disk in two different formats: in a `BoltDB database <https://github.com/boltdb/bolt>`_, or in `JSON <https://en.wikipedia.org/wiki/JSON>`_ format.

.. figure:: simulations.mermaid.png

requests.db
...........

By default, Hoverfly stores simulation data on disk in a file called `requests.db`. This file is written to the current working directory. This means simulations do not need to be manually exported and re-imported between Hoverfly invocations.

.. warning::

    Please note that hoverctl, on the other hand, does not write the `requests.db` file to disk. It may be confusing that hoverctl and Hoverfly behave differently by default. This may be fixed in future releases.

This mechanism uses a very high performance Golang database system: `BoltDB <https://github.com/boltdb/bolt>`_.

.. note::

    It is important to remember that :ref:`delays` do not get stored in the requests.db file. This means any delay configuration you add to your simulations will not persist across invocations, unless you store them using the JSON serialisation mechanism.

<simulation>.json
.................

Simulation JSON can be exported, edited and imported in and out of Hoverfly, and can be shared among Hoverfly users or instances. Simulation JSON files must adhere to the Hoverfly :ref:`simulation_schema`.

.. seealso::

    For a hands-on tutorial of creating and editing simulations, see :ref:`simulations_io`.
