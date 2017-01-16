.. _delays:

Delays
======

Once you have created a simulated service by capturing traffic between
your application and an external service, you may wish to make the
simulation more "realistic" by applying latency to the responses
returned by Hoverfly.

This could be done using :ref:`middleware`. However, if you do not want to go to the effort of writing a
middleware script, you can use a JSON file to apply a set of fixed
response delays to the Hoverfly simulation.

This method is useful if Hoverfly is being used in a load test to
simulate an external service, and there is a requirement to simulate
external service latency. Under high load, the overhead of executing
middleware scripts will impact the performance of Hoverfly, making the
middleware approach to adding latency unsuitable.

Delays based on URL or HTTP method
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Hoverfly can be configured to apply delays to responses based on URL pattern matching or HTTP
method. This is done using a regular expression to match against the URL, a delay value in milliseconds,
and an optional HTTP method value.

.. seealso::

  This functionality is best understood via a practical example: see :ref:`adding_delays` in the :ref:`tutorials` section.
