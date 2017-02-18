.. _proxyauth:

Enabling authentication for the Hoverfly proxy and API
======================================================

If you are running Hoverfly on a remote host (see :ref:`remotehoverfly`), you may wish to 
enable authentication on the Hoverfly proxy and API.

In-memory vs BoltDB persistence
-------------------------------

By default, Hoverfly persists data in memory. Simulation data can then be written to disk using 
the ``hoverctl export`` command (see :ref:`simulations_io`).

Optionally, Hoverfly can be configured to persist data on disk in a ``requests.db`` file (see
:ref:`simulations` for more information). This is implemented using the `BoltDB <https://github.com/boltdb/bolt>`_
database system.

Since authentication credentials must be persisted, **Hoverfly must be configured to use BoltDB
in order for authentication to work**.

Setting Hoverfly authentication credentials
-------------------------------------------

In this example, we assume that the steps in the :ref:`remotehoverfly` tutorial have been followed, 
and that the Hoverfly binary is installed **but not running** on a remote host.

On the **remote host**, run the following command to set the Hoverfly authentication credentials.

.. literalinclude:: hoverfly-set-user-pass.sh
   :language: sh 

This command will start Hoverfly, set the credentials, and then terminate Hoverfly. A ``requests.db``
file will be created in the working directory.

Now run the following command to start Hoverfly with authentication enabled, and the default admin and 
proxy ports overridden.

.. literalinclude:: hoverfly-start-proxy-auth.sh
   :language: sh 

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
   allow hoverctl to authenticate against remote Hoverfly instance admin API. 
   
   The proxy authentication credentials must **also be provided to the application making the request**,
   in this case, cURL.

