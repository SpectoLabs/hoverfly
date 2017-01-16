.. _hoverctl_commands:

hoverctl commands
=================

This page contains the output of:

.. code:: bash

    hoverctl --help

The command's help content has been placed here for convenience.

::

    usage: hoverctl [<flags>] <command> [<args> ...]

    Flags:
          --help                     Show context-sensitive help (also try
                                     --help-long and --help-man).
      -v, --verbose                  Verbose mode.
          --host=HOST                Set the host of Hoverfly
          --admin-port=ADMIN-PORT    Set the admin port of Hoverfly
          --proxy-port=PROXY-PORT    Set the proxy port of Hoverfly
          --certificate=CERTIFICATE  Supply path for custom certificate
          --key=KEY                  Supply path for custom key
          --disable-tls              Disable TLS verification
          --version                  Show application version.

    Commands:
      help [<command>...]
        Show help.

      mode [<name>]
        Get Hoverfly's current mode

      destination [<flags>] [<name>]
        Get Hoverfly's current destination

      middleware [<path>]
        Get Hoverfly's middleware

      start [<server type>]
        Start a local instance of Hoverfly

      stop
        Stop a local instance of Hoverfly

      export <name>
        Exports data out of Hoverfly

      import [<flags>] <name>
        Imports data into Hoverfly

      delete [<resource>]
        Delete test data from Hoverfly

      delays [<path>]
        Get per-host response delay config currently loaded in Hoverfly

      logs [<flags>]
        Get the logs from Hoverfly

      templates [<path>]
        Get set of request templates currently loaded in Hoverfly

      config
        Get the config being used by hoverctl and Hoverfly
