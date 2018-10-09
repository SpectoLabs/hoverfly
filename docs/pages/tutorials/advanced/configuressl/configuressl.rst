.. _configuressl:

Configuring SSL in Hoverfly
===========================

Hoverfly supports both one-way and two-way SSL authentication.

Hoverfly uses default certificate which you should add to your HTTPS client's trust store for one-way SSL authentication. You have options to provide your own certificate, please see below.

Override default certificate for one-way SSL authentication
-----------------------------------------------------------

In some cases, you may not wish to use Hoverfly's default SSL certificate. Hoverfly allows 
you to generate a new certificate and key.

The following command will start a Hoverfly process and create new ``cert.pem`` and ``key.pem``
files in the current working directory. These newly-created files will be loaded into the
running Hoverfly instance.

.. literalinclude:: hoverfly-generate-ca.sh
   :language: sh 

Optionally, you can provide a custom certificate name and authority:

.. literalinclude:: hoverfly-generate-ca-custom-name.sh
   :language: sh 

Once you have generated ``cert.pem`` and ``key.pem`` files with Hoverfly, you can use hoverctl 
to start an instance of Hoverfly using these files.

.. literalinclude:: hoverctl-start-with-custom-cert.sh
   :language: sh

.. note::
   Both a certificate and a key file must be supplied. The files must be in unencrypted PEM format.


Configure Hoverfly for two-way SSL authentication
-------------------------------------------------

For two-way or mutual SSL authentication, you should provide Hoverfly with a client certificate and a certificate key  that you use to authenticate with the remote server.

Two-way SSL authentication is only enabled for request hosts that match the value you provided to the ``--client-authentication-destination`` flag. You can also pass a regex pattern if you need to match multiple hosts.

.. code:: bash

    hoverctl start --client-authentication-client-cert cert.pem --client-authentication-client-key key.pem --client-authentication-destination <host name of the remote server>


If you need to provide a CA cert, you can do so using the ``--client-authentication-ca-cert`` flag.