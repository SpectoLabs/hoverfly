.. _proxyauth:

Enabling authentication for the Hoverfly proxy and API
======================================================

If you are running Hoverfly on a remote host (see :ref:`remotehoverfly`), you may wish to 
enable authentication on the Hoverfly proxy and API.


Setting Hoverfly authentication credentials
-------------------------------------------

In this example, we assume that the steps in the :ref:`remotehoverfly` tutorial have been followed, 
and that the Hoverfly binary is installed **but not running** on a remote host.

On the **remote host**, run the following command to start Hoverfly with authentication credentials, 
and the default admin and proxy ports overridden.

.. literalinclude:: hoverfly-start-proxy-auth.sh
   :language: sh 

.. warning::
  
   By default, Hoverfly starts with authentication disabled. If you require authentication
   you must make sure the ``-auth``, ``-username`` and ``-password`` flags are supplied every
   time Hoverfly is started. 

Configuring hoverctl
--------------------

On your **local machine**, edit your ``~/.hoverfly/config.yml`` so it looks like this:

.. literalinclude:: config.yml
   :language: none

The ``config.yml`` file will now tell hoverctl the location of the remote Hoverfly instance,
and provide the credentials required to authenticate.

Run the following commands **on your local machine** to capture and simulate a URL using the 
remote Hoverfly:

.. literalinclude:: curl-proxy-basic-auth.sh
   :language: sh

.. note::
   The ``hoverfly.username`` and ``hoverfly.password`` values in the ``config.yml`` file 
   allow hoverctl to authenticate against the admin API of a remote Hoverfly instance. 
   
   When using authentication, the Hoverfly proxy port is also authenticated using basic HTTP authentication. 
   This means that any application using Hoverfly (in this example, cURL) must include the authentication credentials as part of the proxy URL. 
