.. _managingstate:


Managing state via Hoverctl
===========================

It could be tricky to reason about the current state of Hoverfly, or to get Hoverfly in a state that you desire for testing.
This is why Hoverctl comes with commands that let you orchestrate it's state. Some useful commands are:

.. code:: bash

    $ hoverctl state --help
    $ hoverctl state get-all
    $ hoverctl state get key
    $ hoverctl state set key value
    $ hoverctl state delete-all
