.. _randomlatency:

Using middleware to simulate network latency
============================================

.. todo:: I'm not sure I like the idea of show people middleware by duplicating functionality that is in Hoverfly...

.. seealso::

    Please carefully read through :ref:`middleware` alongside these tutorials to gain a high-level understanding of what we are about to cover.

We will use a Python script to apply a random delay of less than one second to every response in a simulation.

Before you proceed, please ensure that you have Python installed. 

Let's begin by writing our middleware. Save the following as ``middleware.py``:

.. literalinclude:: middleware.py
    :language: python

The middleware script delays each response by a random value of less than one second.

.. literalinclude:: modify.sh
    :language: sh

Middleware gives you control over the behaviour of a simulation, as well as the data.

.. note::

  Middleware gives you flexibility when simulating network latency - allowing you to randomize the delay value
  for example - but a new process is spawned every time the middleware script is executed. This can impact
  Hoverfly's performance under load.

  If you need to simulate latency during a load test, it is recommended that you use Hoverfly's native :ref:`delays`
  functionality to simulate network latency (see :ref:`adding_delays`) instead of writing middleware. The delays functionality sacrifices flexibility for performance.
