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
