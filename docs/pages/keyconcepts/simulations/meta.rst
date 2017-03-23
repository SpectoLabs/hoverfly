.. _meta:

Meta
====

The last part of the simulation schema is the meta object. Its purpose is to store metadata that is relevant to your simulation. This includes the simulation schema version, the version of Hoverfly used to export the simulation and the date and time at which the simulation was exported.

.. todo:: simulations/basic-simulation.json lines 40 to 44

.. code:: json

    "meta": {
      "schemaVersion": "v2",
      "hoverflyVersion": "v0.11.0",
      "timeExported": "2016-11-11T11:53:52Z"
    }
