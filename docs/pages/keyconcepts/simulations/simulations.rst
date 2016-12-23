Simulations
-----------

The core functionality of Hoverfly is to capture HTTP traffic, and then simulate it; i.e play it back. To do so, it internally stores `simulations`.

Simulations can be written to disk in two different formats: a `boltdb database <https://github.com/boltdb/bolt>`_, or in `json <https://en.wikipedia.org/wiki/JSON>`_ format.

.. figure:: simulations.mermaid.png

requests.db
...........

By default hoverfly stores simulation data to disk in a file called `requests.db`, in its current working directory. This means simulations do not need to be manually exported and re-imported between Hoverfly invocations.

.. warning::

    Please note that hoverctl, on the other hand, does not store requests.db to disk. It may be confusing that hoverctl and hoverfly behave differently by default, and may be fixed in future releases.

This mechanism uses a very high performance GoLang database system: `boltdb <https://github.com/boltdb/bolt>`_. Boltdb is part of the mechanisms, and technology choices that enable Hoverfly to load very rapidly.

.. note::

    It is important to remember that :ref:`delays` do not get stored in the requests.db file. This means any delays information you add to your simulations will not persist across invocations, unless you store them using the json serialisation mechanism.

<simulation>.json
.................

Simulations can also be exported in json format.

Simulations can be exported, edited and imported in and out of Hoverfly in json format, or shared among Hoverfly users or instances. Simulation json files must adhere to its schema to remain readable.

:ref:`simulationschema`

.. seealso::

    For a hands-on tutorial of creating and editing simulations, see :ref:`simulations_io`.

