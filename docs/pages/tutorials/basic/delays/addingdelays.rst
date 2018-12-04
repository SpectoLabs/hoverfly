.. _adding_delays:

Adding delays to a simulation
=============================

Simulating API latency during development allows you to write code that will deal with it gracefully. 

In Hoverfly, this is done by applying "delays" or "delaysLogNormal" to responses in a simulation.

Delays are applied by editing the Hoverfly simulation JSON file. Delays 
can be applied selectively according to request URL pattern and/or HTTP method.

.. toctree::
    :maxdepth: 3

    allresponses/allresponses
    multiplehosts/multiplehosts
    multiplelocations/multiplelocations
    httpmethods/httpmethods
    lognormal/lognormal
