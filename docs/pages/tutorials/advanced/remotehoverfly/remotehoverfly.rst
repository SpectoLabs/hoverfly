.. _remotehoverfly:

Controlling a remote Hoverfly instance with hoverctl
====================================================

So far, the tutorials have shown how hoverctl can be used to control an instance of Hoverfly running on the same machine.

In some cases, you may wish to use hoverctl to control an instance of Hoverfly running on a remote host. With hoverctl,
you can do this using the **targets** feature.

In this example, we assume that the remote host is reachable at ``hoverfly.example.com``, and that 
ports ``8880`` and ``8555`` are available. We will also assume that the Hoverfly binary is installed on the remote host.

On the **remote host**, start Hoverfly using flags to override the default admin port (``-ap``) and proxy port (``-pp``).

.. literalinclude:: run-hoverfly-remote-host.sh
   :language: sh

.. seealso:: 
   For a full list of all Hoverfly flags, please refer to :ref:`hoverfly_commands` in the :ref:`reference` section.
   

On your **local machine**, you can create a **target** named ``remote`` using hoverctl. This target will be configured to communicate
with Hoverfly.

.. literalinclude:: create-hoverctl-target.sh
   :language: none

Now that hoverctl knows the location of the ``remote`` Hoverfly instance, run the following commands
**on your local machine** to capture and simulate a URL using this instance:

.. literalinclude:: curl-proxy-remote-hoverfly.sh
   :language: sh

You will now need to specify the ``remote`` target every time you want to interact with this Hoverfly instance.
If you are only working with this remote instance, you can set it to be the default target instance for hoverctl.

.. literalinclude:: default-hoverctl-target.sh
   :language: sh

.. note::
   The ``--host`` value of the hoverctl target allows hoverctl to interact with the **admin API** 
   of the remote Hoverfly instance.
   
   The application that is making the request (in this case, cURL), **also** needs to be configured to 
   use the remote Hoverfly instance as a proxy. In this example, it is done using cURL's ``--proxy`` flag.  

If you are running Hoverfly on a remote host, you may wish to enable authentication on the Hoverfly proxy and admin API.
This is described in the :ref:`proxyauth` tutorial.                 