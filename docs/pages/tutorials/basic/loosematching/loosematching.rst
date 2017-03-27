.. _loosematching:

Loose request matching using a Request Matcher
==============================================

.. seealso::

    Please carefully read through :ref:`pairs` alongside this tutorial to gain a high-level understanding of what we are about to cover.

In some cases you may want Hoverfly to return the same stored response for more than one incoming request. This can be done using
Request Matchers. 

Let's begin by capturing some traffic and exporting a simulation. This step saves us having to manually create a simulation ourselves and gives us a request to work with.

.. literalinclude:: loosematching.sh
    :language: bash

If you take a look at your ``simulation.json`` you should notice these lines in your request.

.. literalinclude:: simulation.json
    :language: javascript

Modify them to:

.. literalinclude:: simulationimport.json
    :language: javascript

Save the file as ``simulationimport.json`` and run the following command to import it and cURL the simulated endpoint:

.. literalinclude:: loosematchingimport.sh
    :language: bash

The same response is returned, even though we created our simulation with a request to ``http://echo.jsontest.com/foo/baz/bar/spam`` in Capture mode and then sent a request to ``http://echo.jsontest.com/foo/QUX/bar/spam`` in Simulate mode.

.. seealso::

   In this example we used the :code:`globMatch` Request Matcher type. For a list of other Request Matcher types and examples
   of how to use them, please see the :ref:`request_matchers` section.



.. note:: Key points:

    - To change how incoming requests are matched to stored responses, capture a simulation, export it, edit it
    - While editing, choose a request field to match on, select a Request Matcher type and a matcher value 
    - Re-import the simulation
    - Requests can be manually added without capturing the request
