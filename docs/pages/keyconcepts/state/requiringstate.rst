.. _requiringstate:


Requiring State in order to Match
=================================

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