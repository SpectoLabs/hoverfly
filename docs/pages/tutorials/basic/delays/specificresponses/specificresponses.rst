Applying a delay to specific responses
======================================

When you want to add a delay to specific responses and global host / uri / http method matching is not the case,
you can specify **fixedDelay** in response. Here we apply a delay of *3000ms* only to request **/api/profile** with
**X-API-Version: v1** header:

.. code:: json

    {
        "request": {
            "path": [
                {"matcher": "exact", "value": "/api/profile"}
            ],
            "headers": {
                "X-API-Version": [
                    {"matcher": "exact", "value": "v1"}
                ]
            }
        },
        "response": {
            "status": 404,
            "body": "Page not found",
            "fixedDelay": 3000
        }
    }

It's also possible to apply a logNormal delay:

.. code:: json

    {
        "request": {
            "path": [
                {"matcher": "exact", "value": "/api/profile"}
            ],
            "headers": {
                "X-API-Version": [
                    {"matcher": "exact", "value": "v1"}
                ]
            }
        },
        "response": {
            "status": 404,
            "body": "Page not found",
            "logNormalDelay": {
                "min": 100,
                "max": 10000,
                "mean": 5000,
                "median": 500
            }
        }
    }

Like global delays, when both `fixedDelay` and `logNormalDelay` are provided they are applied one after another.