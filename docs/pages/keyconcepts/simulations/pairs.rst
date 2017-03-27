.. _pairs:

Request Responses Pairs
=======================

Hoverfly simulates APIs by matching **incoming requests** from the client to **stored requests**. Stored requests have an associated
**stored response** which is returned to the client if the match is successful.  

The matching logic that Hoverfly uses to compare incoming requests with stored requests can be configured using **Request Matchers**. 

Request Matchers
----------------

When Hoverfly captures a request, it creates a Request Matcher for each field in the request. A Request Matcher consists of
the request field name, the type of match which will be used to compare the field in the incoming request to the field in the stored
request, and the request field value.
 
By default, Hoverfly will set the type of match to :code:`exactMatch` for each field. 

.. seealso::

    There are many types of Request Matcher. Please refer to :ref:`request_matchers` for a list of the types available, and
    examples of how to use them.

    
An example Request Matcher Set might look like this:

+-------------+-------------+------------------------------------+
| Field       | Matcher Type| Value                              |
+=============+=============+====================================+
| scheme      | exactMatch  | "https"                            |
+-------------+-------------+------------------------------------+
| method      | exactMatch  | "GET"                              |
+-------------+-------------+------------------------------------+
| destination | exactMatch  | "docs.hoverfly.io"                 |         
+-------------+-------------+------------------------------------+
| path        | exactMatch  | "/pages/keyconcepts/templates.html"|
+-------------+-------------+------------------------------------+
| query       | exactMatch  | "query=true"                       |
+-------------+-------------+------------------------------------+
| body        | exactMatch  | ""                                 |
+-------------+-------------+------------------------------------+
| headers     | exactMatch  |                                    |
+-------------+-------------+------------------------------------+

In the Hoverfly simulation JSON file, this Request Matcher Set would be represented like this: 

.. literalinclude:: ../../simulations/basic-simulation.json
   :lines: 4-24
   :linenos:
   :language: javascript


:ref:`View entire simulation file <basic_simulation>`

The matching strategy that Hoverfly uses to compare an incoming request to a stored request can be changed by editing the Request Matchers in the simulation
JSON file. 

It is not necessary to have a Request Matcher for every request field. By omitting Request Matchers, it is possible to implement **partial matching** - meaning
that Hoverfly will return one stored response for multiple incoming requests. 

For example, this Request Matcher will match any incoming request to the :code:`docs.hoverfly.io` destination:

.. literalinclude:: ../../simulations/all-matchers-simulation.json
   :lines: 4-8
   :linenos:
   :language: javascript


:ref:`View entire simulation file <all_matchers_simulation>`


In the example below, the :code:`globMatch` Request Matcher type is used to match any subdomain of :code:`hoverfly.io`:

.. literalinclude:: ../../simulations/all-matchers-simulation.json
   :lines: 18-22
   :linenos:
   :language: javascript
   

:ref:`View entire simulation file <all_matchers_simulation>`

It is also possible to use more than one Request Matcher for each field.

In the example below, a :code:`regexMatch` **and** a :code:`globMatch` are used on the :code:`destination` field. 

This will match on any subdomain of :code:`hoverfly.io` which begins with the letter :code:`d`. This means that
incoming requests to :code:`docs.hoverfly.io` and :code:`dogs.hoverfly.io` will be matched, but requests to 
:code:`cats.hoverfly.io` will not be matched.

.. literalinclude:: ../../simulations/all-matchers-simulation.json
   :lines: 32-37
   :linenos:
   :language: javascript
   

:ref:`View entire simulation file <all_matchers_simulation>`


.. seealso::

    There are many types of Request Matcher. Please refer to :ref:`request_matchers` for a list of the types available, and
    examples of how to use them.

    For a practical example of how to use a Request Matcher, please refer to :ref:`loosematching` in the tutorials section.


Responses
---------

Each Request Matcher Set has a response associated with is. If the request match is successful, Hoverfly will return the response to the client.

.. literalinclude:: ../../simulations/basic-simulation.json
   :lines: 25-32
   :linenos:
   :language: javascript

:ref:`View entire simulation file <basic_simulation>`

Editing the fields in response, combined with editing the Request Matcher set, makes it possible to configure complex request/response logic. 

Binary data in responses
~~~~~~~~~~~~~~~~~~~~~~~~

Since JSON does not support binary data, binary responses are base64 encoded. This is denoted by the encodedBody field. 
Hoverfly automatically encodes and decodes the data during the export and import phases.

.. literalinclude:: ../../simulations/basic-encoded-simulation.json
   :lines: 27-28
   :linenos:
   :language: javascript

:ref:`View entire simulation file <basic_encoded_simulation>`