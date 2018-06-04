.. _pairs:

Request Responses Pairs
=======================

Hoverfly simulates APIs by matching **incoming requests** from the client to **stored requests**. Stored requests have an associated
**stored response** which is returned to the client if the match is successful.  

The matching logic that Hoverfly uses to compare incoming requests with stored requests can be configured using **Request Matchers**. 

.. _matchers:

Request Matchers
----------------

When Hoverfly captures a request, it creates a Request Matcher for each field in the request. A Request Matcher consists of
the request field name, the type of match which will be used to compare the field in the incoming request to the field in the stored
request, and the request field value.
 
By default, Hoverfly will set the type of match to :code:`exact` for each field. 

.. seealso::

    There are many types of Request Matcher. Please refer to :ref:`request_matchers` for a list of the types available, and
    examples of how to use them.

    There alse are two different matching strategies: **strongest match** (default) and **first match** (legacy). Please refer to
    :ref:`matching` for more information. 

    
An example Request Matcher Set might look like this:

+-------------+-------------+------------------------------------+
| Field       | Matcher Type| Value                              |
+=============+=============+====================================+
| scheme      | exact       | "https"                            |
+-------------+-------------+------------------------------------+
| method      | exact       | "GET"                              |
+-------------+-------------+------------------------------------+
| destination | exact       | "docs.hoverfly.io"                 |         
+-------------+-------------+------------------------------------+
| path        | exact       | "/pages/keyconcepts/templates.html"|
+-------------+-------------+------------------------------------+
| query       | exact       | "query=true"                       |
+-------------+-------------+------------------------------------+
| body        | exact       | ""                                 |
+-------------+-------------+------------------------------------+
| headers     | exact       |                                    |
+-------------+-------------+------------------------------------+

In the Hoverfly simulation JSON file, this Request Matcher Set would be represented like this: 

.. literalinclude:: ../../simulations/basic-simulation.json
   :lines: 5-44
   :linenos:
   :language: javascript


:ref:`View entire simulation file <basic_simulation>`

The matching logic that Hoverfly uses to compare an incoming request to a stored request can be changed by editing the Request Matchers in the simulation
JSON file. 

It is not necessary to have a Request Matcher for every request field. By omitting Request Matchers, it is possible to implement **partial matching** - meaning
that Hoverfly will return one stored response for multiple incoming requests. 

For example, this Request Matcher will match any incoming request to the :code:`docs.hoverfly.io` destination:

.. literalinclude:: ../../simulations/all-matchers-simulation.json
   :lines: 6-11
   :linenos:
   :language: javascript


:ref:`View entire simulation file <all_matchers_simulation>`


In the example below, the :code:`globMatch` Request Matcher type is used to match any subdomain of :code:`hoverfly.io`:

.. literalinclude:: ../../simulations/all-matchers-simulation.json
   :lines: 27-32
   :linenos:
   :language: javascript
   

:ref:`View entire simulation file <all_matchers_simulation>`

It is also possible to use more than one Request Matcher for each field.

In the example below, a :code:`regexMatch` **and** a :code:`globMatch` are used on the :code:`destination` field. 

This will match on any subdomain of :code:`hoverfly.io` which begins with the letter :code:`d`. This means that
incoming requests to :code:`docs.hoverfly.io` and :code:`dogs.hoverfly.io` will be matched, but requests to 
:code:`cats.hoverfly.io` will not be matched.

.. literalinclude:: ../../simulations/all-matchers-simulation.json
   :lines: 48-57
   :linenos:
   :language: javascript
   

:ref:`View entire simulation file <all_matchers_simulation>`


.. seealso::

    There are many types of Request Matcher. Please refer to :ref:`request_matchers` for a list of the types available, and
    examples of how to use them.

    For a practical example of how to use a Request Matcher, please refer to :ref:`loosematching` in the tutorials section.

    There alse are two different matching strategies: **strongest match** (default) and **first match** (legacy). Please refer to
    :ref:`matching` for more information. 


Responses
---------

Each Request Matcher Set has a response associated with is. If the request match is successful, Hoverfly will return the response to the client.

.. literalinclude:: ../../simulations/basic-simulation.json
   :lines: 45-55
   :linenos:
   :language: javascript

:ref:`View entire simulation file <basic_simulation>`

Editing the fields in response, combined with editing the Request Matcher set, makes it possible to configure complex request/response logic. 

Binary data in responses
~~~~~~~~~~~~~~~~~~~~~~~~

JSON is a text-based file format so it has no intrinsic support for binary data. Therefore if Hoverfly a response body contains binary data (images, gzipped, etc), the response body will be base64 encoded and the `encodedBody` field set to true.

.. literalinclude:: ../../simulations/basic-encoded-simulation.json
   :lines: 47-48
   :linenos:
   :language: javascript

:ref:`View entire simulation file <basic_encoded_simulation>`
