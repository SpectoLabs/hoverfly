.. _diff_mode:

Diff mode
========

In this mode, Hoverfly forwards a request to an external service and compares a response with currently stored simulation.
With both the stored simulation response and the real response from the external service, Hoverfly is able to detect 
differences between the two. When Hoverfly has finished comparing the two responses, the difference is stored and the
incoming request is served the real response from the external service.

The differences can be retrieved from Hoverfly using the API (`GET /api/v2/diff`).
The response contains a list of differences, containing the request and the differences of the response.

.. code:: json

  {
    "diff": [{
      "request": {
        "method": "GET",
        "host": "time.jsontest.com",
        "path": "/",
        "query": ""
      },
      "diffReports": [{
        "timestamp": "2018-03-16T17:45:40Z",
        "diffEntries": [{
          "field": "header/X-Cloud-Trace-Context",
          "expected": "[ec6c455330b682c3038ba365ade6652a]",
          "actual": "[043c9bb2eafa1974bc09af654ef15dc3]"
        }, {
          "field": "header/Date",
          "expected": "[Fri, 16 Mar 2018 17:45:34 GMT]",
          "actual": "[Fri, 16 Mar 2018 17:45:41 GMT]"
        }, {
          "field": "body/time",
          "expected": "05:45:34 PM",
          "actual": "05:45:41 PM"
        }, {
          "field": "body/milliseconds_since_epoch",
          "expected": "1.521222334104e+12",
          "actual": "1.521222341017e+12"
        }]
      }]
    }]
  }

This data is stored and kept until the Hoverfly instance is stopped or the the storage is cleaned by calling the API (`DELETE /api/v2/diff`).

.. seealso::

    For more information on the API to retrieve differences, see :ref:`rest_api`.