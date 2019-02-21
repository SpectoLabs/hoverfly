.. _capturingsequences:

Capturing a stateful sequence of responses
==========================================

By default Hoverfly will store a given request/response pair once only.
If the same request returns different responses you may want capture the sequence of changing request/response pairs.
You may want to do this because the API is stateful rather stateless.

A simple example of this is an API that returns the time.

To record a sequence of request/responses where the request is the same but the response is different,
we need to enable stateful recording in capture mode.

.. code:: bash

    hoverctl start
    hoverctl mode capture --stateful

Now that we have enabled stateful recording, we can capture several request/response pairs.

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
You will see Hoverfly has added a state variable called sequence:1 that acts as a counter.

.. code:: json

    "requiresState": {
        "sequence:1": "1"
    }

    "requiresState": {
        "sequence:1": "2"
    }

We can also see that the first response has `transitionsState`` field set.

.. code:: json

    "transitionsState": {
        "sequence:1": "2"
    }

Note that Hoverfly will automatically set the state for any "sequence:" key to "1" on import.
If you want to use a more meaningful key name you will need to initialise the state as follows:

.. code:: bash

    hoverctl state set shopping-basket empty

.. seealso::

  For a more detailed explaination of how sequences work in hoverfly: see :ref:`sequences` in the :ref:`keyconcepts` section.
