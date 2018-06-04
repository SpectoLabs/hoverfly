.. _state:


State
-----

Hoverfly contains a map of keys and values which it uses to store it's internal state. Some :ref:`request_matchers` can be made to only
match when Hoverfly is in a certain state, and other matchers can be set to mutate Hoverfly's state.


Requiring State in order to Match
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

A matcher can include a field `requiresState`, which dictates the state Hoverfly must be in for there to be a match:

.. code:: json

    "request": {
        "path": [
            {
                "matcher": "exact",
                "value": "/basket"
            }
        ]
        "requiresState": {
            "eggs": "present",
            "bacon" : "large"
        }
    },
    "response": {
        "status": 200,
        "body": "eggs and large bacon"
    }

In the above case, the following matches results would occur when making a request to `/basket`:

+-------------------------------+----------+----------------------------------------------------+
| Current State of Hoverfly     | matches? | reason                                             |
+===============================+==========+====================================================+
| eggs=present,bacon=large      | true     | Required and current state are equal               |
+-------------------------------+----------+----------------------------------------------------+
| eggs=present,bacon=large,f=x  | true     | Additional state 'f=x' is not used by this matcher |
+-------------------------------+----------+----------------------------------------------------+
| eggs=present                  | false    | Bacon is missing                                   |
+-------------------------------+----------+----------------------------------------------------+
| eggs=present,bacon=small      | false    | Bacon is has the wrong value                       |
+-------------------------------+----------+----------------------------------------------------+

Setting State when Performing a Match
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

A response includes two fields, `transitionsState` and `removesState` which alter Hoverflies internal state during a match:

.. code:: json

    "request": {
        "path": [
            {
                "matcher": "exact",
                "value": "/pay"
            }
        ]
    },
    "response": {
        "status": 200,
        "body": "eggs and large bacon",
        "transitionsState" : {
            "payment-flow" : "complete",
        },
        "removesState" : [
            "basket"
        ]
    }

In the above case, the following changes to Hoverflies internal state would be made on a match:

+----------------------------------+------------------------+----------------------------------------------------+
| Current State of Hoverfly        | New State of Hoverfly? | reason                                             |
+==================================+========================+====================================================+
| payment-flow=pending,basket=full | payment-flow=complete  | Payment value transitions, basket deleted by key   |
+----------------------------------+------------------------+----------------------------------------------------+
| basket=full                      | payment-flow=complete  | Payment value created, basket deleted by key       |
+----------------------------------+------------------------+----------------------------------------------------+
|                                  | payment-flow=complete  | Payment value created, basket already absent       |
+----------------------------------+------------------------+----------------------------------------------------+

Managing state via Hoverctl
~~~~~~~~~~~~~~~~~~~~~~~~~~~

It could be tricky to reason about the current state of Hoverfly, or to get Hoverfly in a state that you desire for testing.
This is why Hoverctl comes with commands that let you orchestrate it's state. Some useful commands are:

.. code:: bash

    $ hoverctl state --help
    $ hoverctl state get-all
    $ hoverctl state get key
    $ hoverctl state set key value
    $ hoverctl state delete-all

Stateful sequences
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Using state, it is possible to recreate a sequence of different responses that may come back given a single request. 
This can be useful when trying to test stateful endpoints.

When defining state for request response pairs, if you prefix your state key with the string ``sequence:``, Hoverfly 
will acknowledge the pair as being part of a stateful sequence. When simulating this sequence, Hoverfly will keep track
of the user's position in the sequence and move them forwards.

Once Hoverfly has reached the end of the sequence, it will continue to return the final response.

.. code:: json

    {
        "request": {
            "requiresState": {
                "sequence:1": "1"
            }
        },
        "response": {
            "status": 200,
            "body": "First response",
            "transitionsState" : {
                "sequence:1" : "2",
            }
        }
        "request": {
            "requiresState": {
                "sequence:1": "2"
            }
        },
        "response": {
            "status": 200,
            "body": "Second response",
        }
    }