.. _capturingsequences:

Capturing a sequence of responses
=================================

You may want capture the same request multiple times. You may to do this because the response served changes.
An example of this could include a response that returns the time.

To record a sequence of duplicate requests, we need to enable stateful recording in capture mode.

.. code:: bash

    hoverctl mode capture --stateful

Now that we have enabled stateful recording, we can make several capture several duplicate requests.

.. code:: bash

    curl --proxy http://localhost:8500 http://time.jsontest.com

    curl --proxy http://localhost:8500 http://time.jsontest.com

Once we have finished capturing requests, we can switch Hoverfly back to simulate mode.

.. code:: bash

    hoverctl mode simulate

Now we are in simulate, we can make the same requests again, and we will see the time
update each request we make until we reach the end of our recorded sequence.

.. code:: bash

    {
        "time": "01:59:21 PM",
        "milliseconds_since_epoch": 1528120761743,
        "date": "06-04-2018"
    }
    {
        "time": "01:59:23 PM",
        "milliseconds_since_epoch": 1528120763647,
        "date": "06-04-2018"
    }

If we look at the simulation captured, can see that the requests have the ``requiresState`` fields set.
a sequence counter.

.. code:: json

    "requiresState": {
        "sequence:1": "1"
    }

    "requiresState": {
        "sequence:1": "2"
    }

We can also see that the first respone has `transitionsState`` field set.

.. code:: json

    "transitionsState": {
        "sequence:1": "2"
    }

.. seealso::

  For a more detailed explaination of how sequences work in hoverfly: see :ref:`sequences` in the :ref:`keyconcepts` section.