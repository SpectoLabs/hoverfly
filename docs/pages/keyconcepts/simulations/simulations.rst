.. _simulations:

Simulations
===========

The core functionality of Hoverfly is to capture HTTP(S) traffic to create API simulations which can be used in testing. Hoverfly stores captured traffic as `simulations`. 

Simulation JSON can be exported, edited and imported in and out of Hoverfly, and can be shared among Hoverfly users or instances. Simulation JSON files must adhere to the Hoverfly :ref:`simulation_schema`.

Simulations consist of **Request Matchers and Responses**, **Delays** and **Metadata** ("Meta").

.. toctree::

    pairs
    delays
    meta

.. seealso::

    For a hands-on tutorial of creating and editing simulations, see :ref:`simulations_io`.