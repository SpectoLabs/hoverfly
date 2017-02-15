.. _hoverctl_commands:

hoverctl commands
=================

This page contains the output of:

.. code:: bash

    hoverctl --help

The command's help content has been placed here for convenience.

::
    hoverctl is the command line tool for Hoverfly

    Usage:
      hoverctl [command]

    Available Commands:
      config      Show hoverctl configuration information
      delete      Delete Hoverfly simulation
      destination Get and set Hoverfly destination
      export      Export a simulation from Hoverfly
      import      Import a simulation into Hoverfly
      logs        Get the logs from Hoverfly
      middleware  Get and set Hoverfly middleware
      mode        Get and set the Hoverfly mode
      start       Start Hoverfly
      stop        Stop Hoverfly
      version     Get the version of hoverctl

    Flags:
          --admin-port string       A port number for the Hoverfly API/GUI. Overrides the default Hoverfly admin port (8888)
          --certificate string      A path to a certificate file. Overrides the default Hoverfly certificate
          --database string         A database type [memory|boltdb]. Overrides the default Hoverfly database type (memory)
          --disable-tls             Disables TLS verification
          --host string             A host on which a Hoverfly instance is running. Overrides the default Hoverfly host (localhost)
          --key string              A path to a key file. Overrides the default Hoverfly TLS key
          --proxy-port string       A port number for the Hoverfly proxy. Overrides the default Hoverfly proxy port (8500)
          --upstream-proxy string   A host for which Hoverfly will proxy its requests to
      -v, --verbose                 Verbose logging from hoverctl

    Use "hoverctl [command] --help" for more information about a command.
