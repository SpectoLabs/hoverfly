.. _proxyauth:

Enabling authentication for the Hoverfly proxy and API
======================================================

Sometimes, you may need authentication on Hoverfly. An example
use case would be when running Hoverfly on a remote host.

Hoverfly can provide authentication for both the admin API
(using a JWT token) and the proxy (using HTTP basic auth).

Setting Hoverfly authentication credentials
-------------------------------------------

To start Hoverfly with authentication, all we need to do is run 
the ``hoverctl start`` command with the flag authentication flag.

.. literalinclude:: hoverfly-start-proxy-auth.sh
   :language: sh 

Running this command will prompt you to enter a username and 
password. This can be bypassed by providing the ``-username``
and ``--password`` flags, though this will leave credentials
in your terminal history.

.. warning::
  
   By default, hoverctl will start Hoverfly with authentication disabled. If you require authentication
   you must make sure the ``--auth`` flag are supplied every time Hoverfly is started. 

Logging in with hoverctl
------------------------

Now that Hoverfly has started, trying to interact with it using
hoverctl will now result in an authentication error. You will now
need to login with hoverctl.

.. literalinclude:: login-hoverctl.sh
   :language: sh 

This will log you in with the default target, assuming that the
Hoverfly instance was started with this target.

There may be situations where a Hoverfly process is started
externally to hoverctl. When this happens, it is often
best practice to create a new target for it. You can do this
with the login command using the ``--new-target`` flag.

In this example, a remote Hoverfly instance has been started
already for us on the ports 8880 and 8550 (the example from 
:ref:`remotehoverfly`). To get started, we need to create a 
new target and log in with it.

.. literalinclude:: login-hoverctl-new-target.sh
   :language: sh 

Run the following commands to capture and simulate a URL using the 
remote Hoverfly:

.. literalinclude:: curl-proxy-basic-auth.sh
   :language: sh
