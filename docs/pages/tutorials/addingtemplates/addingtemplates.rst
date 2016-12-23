.. _addingtemplates:

Adding templates to a simulation
--------------------------------

.. seealso::
    
    Please carefully read through :ref:`templates` alongside this tutorial, to gain a high-level understanding of what we are about to cover.

In this tutorial, we are going to go through the applied steps of generating, and using a matching template.

Let's begin by capturing and exporting a simulation.

.. literalinclude:: templates.sh
    :language: bash

.. Which gives us this output:

.. .. literalinclude:: templates.output
    :language: bash

If you take a look at your ``simulation.json`` you should notice these lines in your request.

.. literalinclude:: simulation.json
    :language: json
    :lines: 33-36
    :dedent: 4

Modify them to:

.. literalinclude:: simulationimport.json
    :language: json
    :lines: 33-36
    :dedent: 4

save the file as ``simulationimport.json``, and run the following to import it and simulate it:

.. literalinclude:: templateimport.sh
    :language: bash

The simulation runs as before, even though we captured our simulation with a request to ``http://echo.jsontest.com/foo/baz/bar/spam`` and ran our simulation with url ``http://echo.jsontest.com/foo/QUX/bar/spam``. In other words, templating allows us to match URLs using `globbing <https://en.wikipedia.org/wiki/Glob_(programming)>`_.

.. note:: Key points:

    - To do templating, capture a simulation, export it, edit it
    - While editing, change requestTypes's value to 'template'
    - re-import the simulation
