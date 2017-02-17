.. _addingtemplates:

Adding templates to a simulation
================================

.. seealso::

    Please carefully read through :ref:`templates` alongside this tutorial to gain a high-level understanding of what we are about to cover.

In this tutorial, we are going to go through the steps required to generate and use a matching template.

Let's begin by capturing some traffic and exporting a simulation.

.. literalinclude:: templates.sh
    :language: bash

Which gives us this output:

 .. literalinclude:: templates.output
    :language: none

If you take a look at your ``simulation.json`` you should notice these lines in your request.

.. literalinclude:: simulation.json
    :language: javascript
    :lines: 33-36
    :dedent: 4

Modify them to:

.. literalinclude:: simulationimport.json
    :language: javascript
    :lines: 33-36
    :dedent: 4

Save the file as ``simulationimport.json`` and run the following command to import it and cURL the simulated endpoint:

.. literalinclude:: templateimport.sh
    :language: bash

The same response is returned, even though we created our simulation with a request to ``http://echo.jsontest.com/foo/baz/bar/spam`` in Capture mode and then sent a request to ``http://echo.jsontest.com/foo/QUX/bar/spam`` in Simulate mode.

As you can see, templating allows us to match URLs using `globbing <https://en.wikipedia.org/wiki/Glob_(programming)>`_.

.. note:: Key points:

    - To do templating, capture a simulation, export it, edit it
    - While editing, change ``"requestTypes"`` value to ``"template"``
    - Substitute strings in URLs with the wildcard ``*`` to return one response for more than one request
    - Re-import the simulation
