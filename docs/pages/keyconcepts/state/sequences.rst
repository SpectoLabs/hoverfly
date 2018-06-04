.. _sequences:


Sequences
=========
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