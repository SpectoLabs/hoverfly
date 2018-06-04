.. _settingstate:


Setting State when Performing a Match
=====================================

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
