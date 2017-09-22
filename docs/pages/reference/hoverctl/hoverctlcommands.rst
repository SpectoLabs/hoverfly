.. _hoverctl_commands:

hoverctl commands
=================

This page contains the output of:

.. code:: bash

    hoverctl --help

The command's help content has been placed here for convenience.

.. literalinclude:: hoverctl.output
   :language: none

hoverctl auto completion
------------------------

hoverctl supplies auto completion for Bash. Run the following command to install the completions.

.. code:: bash

    hoverctl completion

This will create the completion file in your hoverfly directory and create a symbolic link in your bash_completion.d folder.

Optionally you can supply a location for the symbolic link as an argument to the completion command.

.. code:: bash

    hoverctl completion /usr/local/etc/bash_completion.d/hoverctl
