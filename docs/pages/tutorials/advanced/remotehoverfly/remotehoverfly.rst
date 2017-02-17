.. _remotehoverfly:

Controlling a remote Hoverfly instance with hoverctl
====================================================

So far, the tutorials have shown how hoverctl can be used to control an instance of Hoverfly running on the same machine.

In some cases, you may wish to use hoverctl to control an instance of Hoverfly running on a remote host. 

In this example, we assume that the remote host is reachable at ``hoverfly.example.com``, and that 
ports ``8880`` and ``8555`` are available. We will also assume that the hoverfly binary is installed on the remote host.

On the **remote host**, start Hoverfly using flags to override the default admin and proxy ports.

.. literalinclude:: run-hoverfly-remote-host.sh
   :language: sh

On your **local machine**, edit your ``~/.hoverfly/config.yml`` so it looks like this:

.. literalinclude:: config.yml
   :language: none

Now that hoverctl knows the location of the remote Hoverfly instance, run the following commands
**on your local machine** to capture and simulate a URL using the remote Hoverfly:

.. literalinclude:: curl-proxy-remote-hoverfly.sh
   :language: sh

.. note::
   The ``hoverfly.host`` value in the ``config.yml`` file tells hoverctl the remote location of Hoverfly. The
   application that is making the request (in this example, cURL), **also** needs to know the remote location
   of the Hoverfly host.

If you are running Hoverfly on a remote host, you may wish to enable authentication on the Hoverfly proxy and admin API.
This is described in the next tutorial.                 
