.. _proxyauth:

Enabling authentication for the Hoverfly proxy and API
======================================================

Sometimes, you may need authentication on Hoverfly. An example
use case would be when running Hoverfly on a remote host.

Hoverfly can provide authentication for both the admin API
(using a JWT token) and the proxy (using HTTP basic auth).

Setting Hoverfly authentication credentials
-------------------------------------------

To start a Hoverfly instance with authentication enabled, you need to run
the ``hoverctl start`` command with the authentication (``--auth``) flag.

.. literalinclude:: hoverfly-start-proxy-auth.sh
   :language: sh 

Running this command will prompt you to set a username and 
password for the Hoverfly instance. 

This can be bypassed by providing the ``--username``
and ``--password`` flags, although **this will leave credentials
in your terminal history**.

.. warning::
  
   By default, hoverctl will start Hoverfly with authentication disabled. If you require authentication
   you must make sure the ``--auth`` flag are supplied every time Hoverfly is started. 

Logging in to a Hoverfly instance with hoverctl
-----------------------------------------------

Now that a Hoverfly instance has started with authentication enabled, you will need to 
**login** to the instance using hoverctl.

.. literalinclude:: login-hoverctl.sh
   :language: sh 

Running this command will prompt you to enter the username and 
password you set for the Hoverfly instance. Again, this can be bypassed by providing the ``--username``
and ``--password`` flags.

There may be situations in which you need to log into to a Hoverfly instance
that is already running. In this case, it is best practice to create a new **target**
for the instance (please see :ref:`remotehoverfly` for more information on **targets**). You can do this using 
the ``--new-target`` flag.

In this example, a remote Hoverfly instance is already running on the host ```hoverfly.example.com```, with
the ports set to 8880 and 8555 and authentication enabled (the example from :ref:`remotehoverfly`). 
You will need to create a new target (named ``remote``) for the instance and log in with it.

.. literalinclude:: login-hoverctl-new-target.sh
   :language: sh 

You will be prompted to enter the username and password for the instance.

Now run the following commands to capture and simulate a URL using the 
remote Hoverfly:

.. literalinclude:: curl-proxy-basic-auth.sh
   :language: sh