.. _delays:

Delays
======

Once you have created a simulated service by capturing traffic between
your application and an external service, you may wish to make the
simulation more "realistic" by applying latency to the responses
returned by Hoverfly.

Hoverfly can be configured to apply delays to responses based on URL pattern matching or HTTP
method. This is done using a regular expression to match against the URL, a delay value in milliseconds,
and an optional HTTP method value.

.. seealso::

  This functionality is best understood via a practical example: see :ref:`adding_delays` in the :ref:`tutorials` section.

  You can also apply delays to simulations using :ref:`middleware` (see the :ref:`randomlatency` tutorial).
  Using middleware to apply delays sacrifices performance for flexibility. 
