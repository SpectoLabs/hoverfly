.. _diff_mode:

Diff mode
========

In this mode, Hoverfly forwards a request to an external service and compares a response with currently stored simulation
to detect differences between the simulated version and what is in the current service response. When the comparison is done,
then Hoverfly returns the real service's response to the client and stores the differences.

The stored differences can be retrieved by calling the running Hoverfly instance on the path: `/api/v2/diff`.
The response then contains lists of messages with described differences grouped by the same request:

.. code:: json

 {
   "diff":[
      {
         "request":{
            "method":"GET",
            "host":"my.service.com",
            "path":"/users/myaccount",
            "query":""
         },
         "diffMessage":[
            "(1)The \"body/email\" parameter is not same - the expected value was [expected@email.com], but the actual one [actual@email.com]"
         ]
      }
   ]
 }

All reports containing the differences are stored and kept until the Hoverfly instance is stopped or the the storage is cleaned by calling
`DELETE` on the path `/api/v2/diff`.
